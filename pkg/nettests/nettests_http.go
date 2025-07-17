package nettests

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"io"
	"net-tools/pkg/nettest"
	"net/http"
	"time"
)

type HTTPServer struct {
	*Server
}

func NewHTTPServer(server *Server) *HTTPServer {
	return &HTTPServer{Server: server}
}

func (s *HTTPServer) wsHandler(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		default:
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				return
			}
			receiveMsg, err := nettest.Unmarshal(data)
			if err != nil {
				logrus.Warnf("failed to unmarshal: %s %s", err, conn.RemoteAddr())
				return
			}
			logrus.Infof("ws message: %s from %s", receiveMsg.GetRequestID(), conn.RemoteAddr())
			sendMsg := nettest.GenerateMessage(receiveMsg.GetRequestID())
			logrus.Infof("ws message: %s to %s", sendMsg.GetRequestID(), conn.RemoteAddr())
			err = conn.WriteMessage(messageType, xutil.RemoveError(nettest.Marshal(sendMsg)))
			if err != nil {
				logrus.Warnf("failed to write: %s %s", err, conn.RemoteAddr())
				return
			}
		case <-ctx.Done():
			logrus.Debugf("ws test finished in %s", conn.RemoteAddr())
			return
		}
	}
}

func (s *HTTPServer) httpHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Debugf("failed to read: %s %s", err, r.RemoteAddr)
		return
	}
	receiveMsg, err := nettest.Unmarshal(data)
	if err != nil {
		logrus.Warnf("failed to unmarshal: %s %s", err, r.RemoteAddr)
		return
	}
	logrus.Infof("http message: %s from %s", receiveMsg.GetRequestID(), r.RemoteAddr)
	sendMsg := nettest.GenerateMessage(receiveMsg.GetRequestID())
	logrus.Debugf("http message: %s to %s", sendMsg.GetRequestID(), r.RemoteAddr)
	_, err = w.Write(xutil.RemoveError(nettest.Marshal(sendMsg)))
	if err != nil {
		logrus.Warnf("failed to write: %s %s", err, r.RemoteAddr)
	}
}

var upgrader = websocket.Upgrader{}

func (s *HTTPServer) run() {
	r := http.NewServeMux()
	srv := &http.Server{
		Addr:        s.ServerBind,
		Handler:     r,
		ReadTimeout: s.Timeout,
	}
	r.HandleFunc("/http", s.httpHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Warnf("upgrade error: %v", err)
			return
		}
		defer func() {
			_ = conn.Close()
		}()
		s.wsHandler(s.ctx, conn)
	})
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logrus.Fatalf("listen error: %v", err)
		}
	}()
	logrus.Infof("http and ws server listening on %s", s.ServerBind)
	<-s.ctx.Done()
	_ = srv.Shutdown(context.Background())
}
