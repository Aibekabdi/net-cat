package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	ip   = "localhost"
	port = "8989"
)

func main() {
	if len(os.Args[1:]) == 1 {
		port = os.Args[1]
	} else if len(os.Args[1:]) > 1 {
		fmt.Println("[USAGE]: go run server/server.go $port")
		fmt.Println("[USAGE]: go run server/server.go (which will start the default server on port 8989)")
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	time := time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(time)
}
