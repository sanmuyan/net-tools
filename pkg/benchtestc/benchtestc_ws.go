package benchtestc

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net-tools/pkg/benchtest"
	"sync"
	"time"
)

type WSClient struct {
	*Client
}

func NewWSClient(client *Client) *WSClient {
	return &WSClient{
		Client: client,
	}
}
func (c *WSClient) sendMessage(conn *websocket.Conn) (*int64, error) {
	startTime := time.Now().UnixMilli()
	sendMsg := benchtest.GenerateMessage(benchtest.GenerateRequestID())
	logrus.Debugf("%s message: %s to %s", c.Protocol, sendMsg.RequestID, conn.RemoteAddr())
	err := conn.WriteMessage(websocket.TextMessage, xutil.RemoveError(benchtest.Marshal(sendMsg)))
	if err != nil {
		logrus.Warnf("failed to write:: %v", err)
		return nil, err
	}
	_ = conn.SetReadDeadline(time.Now().Add(c.Timeout))
	_, data, err := conn.ReadMessage()
	if err != nil {
		logrus.Warnf("failed to read: %v %s", err, conn.RemoteAddr())
		return nil, err
	}
	receiveMsg, err := benchtest.Unmarshal(data)
	if err != nil {
		logrus.Warnf("failed to unmarshal: %s %s", err, conn.RemoteAddr())
		return nil, err
	}
	timing := time.Now().UnixMilli() - startTime
	logrus.Infof("%s message: %s from %s %s", c.Protocol, receiveMsg.GetRequestID(), conn.RemoteAddr(), timeToStrUnit(timing))
	return &timing, err
}

func (c *WSClient) sendHandler(ctx context.Context, conn *websocket.Conn) {
	for i := 0; i < c.MaxMessages || c.MaxMessages <= 0; i++ {
		select {
		case <-ctx.Done():
			_ = conn.Close()
			return
		default:
			timing, err := c.sendMessage(conn)
			if err != nil {
				c.addErrorCount()
				return
			}
			c.addSuccessCount(*timing)
		}
		time.Sleep(c.Interval)
	}
}

func (c *WSClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := c.createConn()
	if err != nil {
		logrus.Errorf("failed to dial server: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	c.sendHandler(c.ctx, conn)
}

func (c *WSClient) createConn() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", c.Server), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
