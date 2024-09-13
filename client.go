package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
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
		"data" : "qui|do"
	}` + "\n\r"))
	fmt.Println("send...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
