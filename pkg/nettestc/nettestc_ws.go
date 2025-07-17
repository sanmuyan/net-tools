package nettestc

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net-tools/pkg/nettest"
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

func (c *WSClient) sendHandler(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			_ = conn.Close()
			return
		default:
			startTime := time.Now().UnixMilli()
			sendMsg := nettest.GenerateMessage(nettest.GenerateRequestID())
			logrus.Debugf("ws message: %s to %s", sendMsg.RequestID, conn.RemoteAddr())
			err := conn.WriteMessage(websocket.TextMessage, xutil.RemoveError(nettest.Marshal(sendMsg)))
			if err != nil {
				logrus.Warnf("failed to write to tcp server: %v", err)
				return
			}
			_ = conn.SetReadDeadline(time.Now().Add(c.Timeout))
			_, data, err := conn.ReadMessage()
			if err != nil {
				logrus.Warnf("failed to read: %v %s", err, conn.RemoteAddr())
				return
			}
			receiveMsg, err := nettest.Unmarshal(data)
			if err != nil {
				logrus.Warnf("failed to unmarshal: %s %s", err, conn.RemoteAddr())
				return
			}
			endTime := time.Now().UnixMilli()
			logrus.Infof("ws message: %s from %s %dms", receiveMsg.GetRequestID(), conn.RemoteAddr(), endTime-startTime)
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
		logrus.Debugf("ws test finished in %s", conn.RemoteAddr())
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
