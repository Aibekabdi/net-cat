package tcp

import (
	"fmt"
	"net"
	"sync"
)

func (s *Server) Run(ip, port, transferProtocol string) error {
	srv, err := net.Listen(transferProtocol, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return err
	}
	s.Server = srv
	s.Conn = make(map[net.Conn]string)
	s.AllMsg = []string{}
	s.mutex = sync.Mutex{}
	return nil
}

func (s *Server) Shutdown() error {
	s.mutex.Lock()
	for conn := range s.Conn {
		conn.Write([]byte("Server Was Closed!"))
		conn.Close()
	}
	s.mutex.Unlock()
	return s.Server.Close()
}
