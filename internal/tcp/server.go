package tcp

import (
	"fmt"
	"log"
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

func (s *Server) IsConnectable(conn net.Conn) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.Conn) > MaxConnections || MaxConnections == 0 {
		return false
	} else {
		return true
	}
}

func (s *Server) ConnectMessenger(conn net.Conn) {
	log.Println(conn)
	s.Conn[conn] = "lol"
	log.Println(s.Conn)
	if !s.IsConnectable(conn) {
		log.Println("The room is full, please try again later...")
		conn.Write([]byte("The room is full, please try again later..."))
		conn.Close()
		return
	}
}
