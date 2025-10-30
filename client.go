package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Fprintln(conn, name)

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	for {
		fmt.Print("Enter message (or 'exit' to quit): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Fprintln(conn, text)
			fmt.Println("Goodbye!")
			return
		}

		fmt.Fprintln(conn, text)
	}
}
