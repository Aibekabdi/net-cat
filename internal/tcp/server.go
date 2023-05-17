package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
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
	if !s.IsConnectable(conn) {
		fmt.Fprintln(conn, "The room is full, please try again later...")
		conn.Close()
		return
	}
	fmt.Fprint(conn, AuthorizationMessage)
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Authorization: %v", err)
	}
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		fmt.Fprintln(conn, "Try again, name too large, max lenght 20 symbols")
		conn.Close()
		return
	}
	if err := s.addConnection(conn, name); err != nil {
		fmt.Fprint(conn, err.Error())
		conn.Close()
		return
	}
	s.startMessaging(conn)
	s.removeConnection(conn)
}

func (s *Server) addConnection(conn net.Conn, name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if MaxConnections != 0 && len(s.Conn) > MaxConnections {
		return errors.New("The room is full, please try again later...")
	}
	for _, names := range s.Conn {
		if names == name {
			return fmt.Errorf("Name '%s' is Exist [%v]", name, conn.RemoteAddr())
		}
	}
	s.Conn[conn] = name
	log.Printf("Client %s connected by %v", name, conn.RemoteAddr())
	return nil
}

func (s *Server) removeConnection(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Conn, conn)
	log.Printf("Connection %v has left", conn.RemoteAddr())
}

func (s *Server) startMessaging(conn net.Conn) {
	s.loadAllMsg(conn)
	message := fmt.Sprintf(PatternJoinChat, s.Conn[conn])
	s.sendMessage(conn, message)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if message != "\n" {
			time := time.Now().Format(TimeLayout)
			message = fmt.Sprintf(PatternMessage, time, s.Conn[conn], message)
		} else {
			message = ""
		}
		s.sendMessage(conn, message)
		s.saveMsg(message)
	}
	fmt.Fprintf(conn, PatternLeftChat, s.Conn[conn])
}

func (s *Server) loadAllMsg(conn net.Conn) {
	for _, message := range s.AllMsg {
		fmt.Fprint(conn, message)
	}
}

func (s *Server) saveMsg(message string) {
	s.mutex.Lock()
	s.AllMsg = append(s.AllMsg, message)
	s.mutex.Unlock()
}

func (s *Server) sendMessage(conn net.Conn, message string) {
	time := time.Now().Format(TimeLayout)
	if message == "" {
		fmt.Fprintf(conn, PatternTyping, time, s.Conn[conn])
		return
	}

	s.mutex.Lock()
	for con := range s.Conn {
		if con != conn {
			fmt.Fprint(con, "!\n"+message)
		}
		fmt.Fprintf(con, PatternTyping, time, s.Conn[con])
	}
	s.mutex.Unlock()
}
