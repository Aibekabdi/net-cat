package tcp

import (
	"net"
	"sync"
)

type Server struct {
	Server net.Listener
	Conn   map[net.Conn]string
	AllMsg []string
	mutex  sync.Mutex
}
