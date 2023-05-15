package main

import (
	"errors"
	"fmt"
	"log"
	"net-cat/internal/tcp"
	"os"
	"os/signal"
)

var (
	ip               = "localhost"
	port             = "8989"
	transferProtocol = "tcp"
	err              error
)

func main() {
	port, err = GetPort()
	if err != nil {
		fmt.Println(err)
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	srv := new(tcp.Server)
	go func() {
		if err := srv.Run(ip, port, transferProtocol); err != nil {
			log.Fatalf("error occured while trying to create server: %s", err)
			return
		}
		log.Printf("Listening on the port : %s\n", port)

		for {
			conn, err := srv.Server.Accept()
			if err != nil {
				break
			}
			go srv.ConnectMessenger(conn)
		}

	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Println("Server is Closing")

	if err = srv.Shutdown(); err != nil {
		log.Fatal(err)
	}
	log.Println("Server was closed")

}

func GetPort() (string, error) {
	args := os.Args
	if len(args) < 2 {
		return port, nil
	} else if len(args) > 2 {
		return "", errors.New("input arguments more than 1")
	}
	return ":" + os.Args[1], nil
}
