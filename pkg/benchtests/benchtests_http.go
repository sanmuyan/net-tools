package benchtests

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sanmuyan/xpkg/xcrypto"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"io"
	"net-tools/pkg/benchtest"
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
			receiveMsg, err := benchtest.Unmarshal(data)
			if err != nil {
				logrus.Warnf("failed to unmarshal: %s %s", err, conn.RemoteAddr())
				return
			}
			logrus.Infof("%s message: %s from %s", s.Protocol, receiveMsg.GetRequestID(), conn.RemoteAddr())
			sendMsg := benchtest.GenerateMessage(receiveMsg.GetRequestID())
			logrus.Infof("%s message: %s to %s", s.Protocol, sendMsg.GetRequestID(), conn.RemoteAddr())
			err = conn.WriteMessage(messageType, xutil.RemoveError(benchtest.Marshal(sendMsg)))
			if err != nil {
				logrus.Warnf("failed to write: %s %s", err, conn.RemoteAddr())
				return
			}
		case <-ctx.Done():
			logrus.Debugf("%s test finished in %s", s.Protocol, conn.RemoteAddr())
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
	receiveMsg, err := benchtest.Unmarshal(data)
	if err != nil {
		logrus.Warnf("failed to unmarshal: %s %s", err, r.RemoteAddr)
		return
	}
	logrus.Infof("%s message: %s from %s", s.Protocol, receiveMsg.GetRequestID(), r.RemoteAddr)
	sendMsg := benchtest.GenerateMessage(receiveMsg.GetRequestID())
	logrus.Debugf("%s message: %s to %s", s.Protocol, sendMsg.GetRequestID(), r.RemoteAddr)
	_, err = w.Write(xutil.RemoveError(benchtest.Marshal(sendMsg)))
	if err != nil {
		logrus.Warnf("failed to write: %s %s", err, r.RemoteAddr)
	}
}

var upgrader = websocket.Upgrader{}

func (s *HTTPServer) run() {
	tlsConfig := xutil.RemoveError(xcrypto.CreateCertToTLS(nil))
	r := http.NewServeMux()
	srv := &http.Server{
		Addr:        s.ServerBind,
		Handler:     r,
		ReadTimeout: s.Timeout,
		TLSConfig:   tlsConfig,
	}
	r.HandleFunc("/ping", s.httpHandler)
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
		var err error
		if s.Protocol == "https" {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("listen error: %v", err)
		}
	}()
	logrus.Infof("%s server listening on %s", s.Protocol, s.ServerBind)
	<-s.ctx.Done()
	_ = srv.Shutdown(context.Background())
}
