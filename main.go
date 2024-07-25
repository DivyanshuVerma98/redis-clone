package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var Cache sync.Map

func main() {
	fmt.Println("Launching...🚀")
	redis := CreateRedis()
	listener, err := net.Listen("tcp", "localhost:8069")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening at 8069🔥🔥..")
	defer listener.Close()
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle client connection in a goroutine
		client := NewClient(conn)
		redis.AddClient(client)
		go handleClient(client, redis)
	}
}

func handleClient(client *Client, redis *Redis) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovering from error", err)
			client.conn.Close()
			redis.RemoveClient(client)
		}
	}()
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Split(input, " ")
		redis.HandleCommand(parts[0], parts[1:], client)
	}
}
