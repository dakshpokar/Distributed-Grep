package main

import (
	"fmt"
	"bufio"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Getenv("REMOTE_IP")+":1200")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Connect to the address with tcp
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	fmt.Printf("Sending to %s\n", tcpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send a message to the server
	_, err = conn.Write([]byte(`{
		"req_type" : "cmd",
		"data" : "quick.*fox"
	}` + "\n\r"))
	fmt.Println("send...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
func handleHeartbeat(conn net.Conn){
	isAlive := true
	for isAlive{
		_, err = conn.Write([]byte(`{
			"req_type" : "heartbeat",
			"data" : "test"
		}` + "\n\r"))
		fmt.Println("send hb")
		if err != nil {
			fmt.Print(err)
		}
		err := conn.SetReadDeadline(time.Now().Add(2*time.Second))
		if err != nil{
			fmt.Println("Read deadline set error")
			isAlive = false
			break
		}
		request, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil{
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timeout waiting for heartbeat response")
			} else {
				fmt.Println("Error reading heartbeat response:", err)						}
			isAlive = false
			break
		}
		fmt.Println(string(request))
		time.Sleep(2 * time.Second)

	}
}
