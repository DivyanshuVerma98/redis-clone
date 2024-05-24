package main

import (
	"fmt"
	"log"
	"net"

	"github.com/DivyanshuVerma98/redis-clone/structs"
)

func main() {
	fmt.Println("Starting...🚀")
	listener, err := net.Listen("tcp", "localhost:8069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		log.Println("New connection", conn)
		if err != nil {
			log.Fatal(err)
		}
		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		log.Println("Closing connection", conn)
		conn.Close()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovering from error", err)
		}
	}()
	// Create a buffer to read data into
	// buffer := make([]byte, 1024)

	parser := structs.NewParser(conn)
	for {
		line, err := parser.ReadLine()
		if err != nil {
			log.Println(err)
			conn.Write([]uint8("-ERR " + err.Error() + "\r\n"))
			break
		}
		// r := bufio.NewReader(conn)
		// line, err := r.ReadBytes('\r')
		// if err != nil {
		// 	fmt.Println("Error", err)
		// 	return
		// }

		// fmt.Println(string(line[:len(line)-1]), len(string(line)))
		if string(line) == "ping" {
			fmt.Println("Done !!")
			conn.Write([]byte("Pong\r\n"))
		}
		fmt.Printf("Line --> %s\r\n", line)
		// conn.Write([]byte(line))
		// Read data from the client
		// n, err := conn.Read(buffer)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return
		// }
		// Process and use the data (here, we'll just print it)
		// fmt.Printf("Received: %s\n", buffer[:n])
	}
}
