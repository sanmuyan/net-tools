package speedtests

type ServerConn interface {
	run()
}

type Server struct {
	ServerBind string
	ServerPort int
	Protocol   string
}

func NewServer(serverBind string, serverPort int, protocol string) *Server {
	return &Server{
		ServerBind: serverBind,
		ServerPort: serverPort,
		Protocol:   protocol,
	}
}

func Start(server *Server) {
	var s ServerConn
	if server.Protocol == "tcp" {
		s = NewTCPServer(server)
	} else {
		s = NewUDPServer(server)
	}
	s.run()
}
