package nettestc

import (
	"bytes"
	"fmt"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"io"
	"net-tools/pkg/nettest"
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

func (c *HTTPClient) sendHandler() {
	startTime := time.Now().UnixMilli()
	sendMsg := nettest.GenerateMessage(nettest.GenerateRequestID())
	logrus.Debugf("http message: %s to %s", sendMsg.GetRequestID(), c.Server)
	helloReader := bytes.NewBuffer(xutil.RemoveError(nettest.Marshal(sendMsg)))
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/http", c.Server), helloReader)
	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("http request error: %v", err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Warnf("failed to read: %s %s", err, c.Server)
		return
	}
	receiveMsg, err := nettest.Unmarshal(data)
	if err != nil {
		logrus.Warnf("failed to unmarshal: %s %s", err, c.Server)
		return
	}
	endTime := time.Now().UnixMilli()
	logrus.Infof("http message: %s from %s %dms", receiveMsg.GetRequestID(), c.Server, endTime-startTime)
}

func (c *HTTPClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.sendHandler()
		}
		time.Sleep(c.Interval)
	}
}
