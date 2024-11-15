package main

import (
	"MP1/types"
	"MP1/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

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

/*
 * handleConnection: Handle incoming connections
 *
 * args:
 * 	{conn} net.Conn: Connection object
 *
 * returns:
 * 	{nil}
 */
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())
	for {
		request, err := bufio.NewReader(conn).ReadString('\r')
		if err != nil {
			fmt.Println(err)
			return
		}
		var req types.Request
		err = json.Unmarshal([]byte(request), &req)
		if err != nil {
			fmt.Println("error in json unmarshall")
			fmt.Print(err)
		}
		if req.Req_type == "cmd" {
			data, _ := utils.Grep(req.Data, req.File)
			joinedString := strings.Join(data, "\\n")
			utils.ReturnOutput(conn, joinedString)
		}
		if req.Req_type == "heartbeat" {
			data := "Alive\n"
			_, err := conn.Write([]byte(data))
			if err != nil {
				fmt.Println("error in json unmarshall")
				fmt.Print(err)
			}
		}
	}
}
