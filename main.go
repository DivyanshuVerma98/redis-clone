package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

var Cache sync.Map

func main() {
	fmt.Println("Launching...🚀")
	listener, err := net.Listen("tcp", "localhost:8069")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening at 8069🔥🔥")
	defer listener.Close()
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		log.Println("New connection", conn)
		if err != nil {
			log.Fatal(err)
		}
		// Handle client connection in a goroutine
		client := newClient(conn)
		go handleClient(client)
	}
}

func handleClient(client *client) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovering from error", err)
			client.conn.Close()
		}
	}()

	for {
		r := bufio.NewReader(client.conn)
		line, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		client.handle(line)
	}
}
