package main

import (
	"MP1/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	Req_type string
	Data     string
}

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
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())
	for {
		request, err := bufio.NewReader(conn).ReadString('\r')
		if err != nil {
			fmt.Println(err)
			return
		}
		var req Request
		fmt.Print(request)
		err = json.Unmarshal([]byte(request), &req)
		if err != nil {
			fmt.Println("error in json unmarshall")
			fmt.Print(err)
		}
		fmt.Print("> ", req)
		if req.Req_type == "cmd" {
			// Print the data read from the connection to the terminal
			// Write back the same message to the client.
			data, _ := utils.Grep(req.Data, "sample.txt")
			fmt.Print("Sending data to client: ", strings.Join(data, ""))
			utils.ReturnOutput(conn.RemoteAddr().(*net.TCPAddr).IP.String(), strings.Join(data, ""))
		}
		if req.Req_type == "heartbeat" {
			data := "Alive\n"
			_, err := conn.Write([]byte(data))
			if err != nil {
				fmt.Println("error in json unmarshall")
				fmt.Print(err)
			}
			fmt.Println("------------- send data: " + data + "----------end data")
			//utils.ReturnOutput(conn.RemoteAddr().(*net.TCPAddr).IP.String(), data)
		}
	}
}
