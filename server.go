package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan string)
	mu        sync.Mutex
	history   []string
)

func sendHistory(conn net.Conn) {
	fmt.Fprintln(conn, "\n--- Chat History ---")
	for _, msg := range history {
		fmt.Fprintln(conn, msg)
	}
	fmt.Fprintln(conn, "---------------------\n")
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read client name
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	mu.Lock()
	clients[conn] = name
	mu.Unlock()

	// Send chat history to the new client
	sendHistory(conn)

	joinMsg := fmt.Sprintf("%s joined the chat", name)

	mu.Lock()
	history = append(history, joinMsg)
	mu.Unlock()

	broadcast <- joinMsg

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = strings.TrimSpace(msg)

		if msg == "exit" {
			exitMsg := fmt.Sprintf("%s left the chat", name)

			mu.Lock()
			history = append(history, exitMsg)
			mu.Unlock()

			broadcast <- exitMsg
			break
		}

		fullMsg := fmt.Sprintf("%s: %s", name, msg)

		mu.Lock()
		history = append(history, fullMsg)
		mu.Unlock()

		broadcast <- fullMsg
	}

	mu.Lock()
	delete(clients, conn)
	active := len(clients)
	mu.Unlock()

	if active == 0 {
		fmt.Println("Last client left → shutting down server.")
		close(broadcast)
	}
}

func sendMessages() {
	for msg := range broadcast {
		fmt.Println("[SERVER HISTORY] →", msg) // show in server console

		mu.Lock()
		for conn := range clients {
			fmt.Fprintln(conn, msg)
		}
		mu.Unlock()
	}
}

func main() {
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server running on port 1234...")

	go sendMessages()

	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		go handleClient(conn)
	}
}
