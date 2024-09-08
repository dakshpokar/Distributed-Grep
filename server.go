package main

import (
	"MP1/utils"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Write a server that listens on a specified port and greps a file for a given pattern.
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":1200")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	fmt.Printf("Listening on %s\n", tcpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		// Handle new connections in a Goroutine for concurrency
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	pattern, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data read from the connection to the terminal
	fmt.Print("> ", string(pattern))

	// Write back the same message to the client.
	data, _ := utils.Grep(pattern, "sample.txt")
	fmt.Print("Sending data to client: ", strings.Join(data, ""))
	utils.ReturnOutput(conn.RemoteAddr().(*net.TCPAddr).IP.String(), strings.Join(data, ""))
	conn.Close()
}
