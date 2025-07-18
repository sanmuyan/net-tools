package benchtestc

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"io"
	"net-tools/pkg/benchtest"
	"net/http"
	"sync"
	"time"
)

type HTTPClient struct {
	*Client
}

func NewHTTPClient(client *Client) *HTTPClient {
	return &HTTPClient{
		Client: client,
	}
}

func (c *HTTPClient) sendMessage() (*int64, error) {
	startTime := time.Now().UnixMilli()
	sendMsg := benchtest.GenerateMessage(benchtest.GenerateRequestID())
	logrus.Debugf("http message: %s to %s", sendMsg.GetRequestID(), c.Server)
	helloReader := bytes.NewBuffer(xutil.RemoveError(benchtest.Marshal(sendMsg)))
	client := &http.Client{
		Timeout: c.Timeout,
	}
	var err error
	var req *http.Request
	if c.Protocol == "https" {
		req, err = http.NewRequest("GET", fmt.Sprintf("https://%s/ping", c.Server), helloReader)
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		req, err = http.NewRequest("GET", fmt.Sprintf("http://%s/ping", c.Server), helloReader)
	}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("%s request error: %v", c.Protocol, err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Warnf("failed to read: %s %s", err, c.Server)
		return nil, err
	}
	receiveMsg, err := benchtest.Unmarshal(data)
	if err != nil {
		logrus.Warnf("failed to unmarshal: %s %s", err, c.Server)
		return nil, err
	}
	timing := time.Now().UnixMilli() - startTime
	logrus.Infof("%s message: %s from %s %s", c.Protocol, receiveMsg.GetRequestID(), c.Server, timeToStrUnit(timing))
	return &timing, nil
}

func (c *HTTPClient) sendHandler() {
	for i := 0; i < c.MaxMessages || c.MaxMessages <= 0; i++ {
		select {
		case <-c.ctx.Done():
			return
		default:
			timing, err := c.sendMessage()
			if err != nil {
				c.addErrorCount()
				return
			}
			c.addSuccessCount(*timing)
		}
		time.Sleep(c.Interval)
	}
}

func (c *HTTPClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	c.sendHandler()
}
