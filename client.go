package main

import (
	"MP1/types"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var serverStates = make(map[string]bool)
var runningHosts = make(map[string]bool)
var output = make(map[string]string)

//var text string

func main() {
	if len(os.Args) < 3 {
		fmt.Println("pattern and log file name needed")
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	startTime := time.Now()
	for i := 1; i <= 10; i++ {
		host := ""
		if i < 10 {
			host = "0" + strconv.Itoa(i)
		} else {
			host = strconv.Itoa(i)
		}
		go createRequestAndHeartbeat(host, strconv.Itoa(i))
	}
	time.Sleep(100 * time.Millisecond)
	for {
		if len(runningHosts) == 0 && len(serverStates) != 0 {
			fmt.Println("dgrep executed successfully")
			break
		}
	}
	fmt.Println("###########################")
	var req types.Request
	keys := make([]string, 0, len(output))
	for k := range output {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	totalLines := 0
	for _, key := range keys {
		value := output[key]
		err = json.Unmarshal([]byte(value), &req)
		if err != nil {
			fmt.Println("error in json unmarshal")
			fmt.Println(err)
			continue
		}
		count, err := strconv.Atoi(req.Data)
		if err == nil {
			fmt.Printf("Matching lines on machine : %s ::: %v\n", key, req.Data)
			totalLines += count
		} else {
			fmt.Printf("Failed to match any lines on : %s ::: %v\n", key, req.Data)
		}
	}
	endTime := time.Now()
	fmt.Println("###########################")
	fmt.Println("Total Lines: ", totalLines)
	fmt.Println("Execution time: ", endTime.Sub(startTime))
}

func createRequestAndHeartbeat(host string, index string) {
	hostName := os.Getenv("REMOTE_IP") + host + ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", hostName)

	if err != nil {
		fmt.Println(err)
		serverStates[hostName] = false
		delete(runningHosts, hostName)
		output[hostName] = `{
			"req_type": "op",
			"data":	"Server ` + hostName + ` is dead"
		}`
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println(err)
		serverStates[hostName] = false
		output[hostName] = `{
			"req_type": "op",
			"data":	"Server ` + hostName + ` is dead"
		}`
		delete(runningHosts, hostName)
		return
	}
	// go handleHeartbeat(conn, hostName)
	go executeGrep(conn, hostName, index)
}
func executeGrep(conn net.Conn, hostName string, index string) {
	runningHosts[hostName] = true
	_, err := conn.Write([]byte(`{
		"req_type" : "cmd",
		"data" : "` + os.Args[1] + `",
		"file" : "` + os.Args[2] + `.log"
	}` + "\n\r"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	request, err := bufio.NewReader(conn).ReadString('\r')
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output[hostName] = string(request)
	_, ok := runningHosts[hostName]
	if ok {
		delete(runningHosts, hostName)
	}
}

func handleHeartbeat(conn net.Conn, hostName string) {
	isAlive := true
	serverStates[hostName] = true
	for isAlive {
		_, err := conn.Write([]byte(`{
			"req_type" : "heartbeat",
			"data" : "ping"
		}` + "\n\r"))
		if err != nil {
			fmt.Print(err)
		}
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("Read deadline set error")
			isAlive = false
			break
		}
		_, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timeout waiting for heartbeat response")
			} else {
				fmt.Println("Error reading heartbeat response:", err)
			}
			isAlive = false
			break
		}
	}
	if !isAlive {
		fmt.Println("Server " + hostName + " is dead")
		serverStates[hostName] = false
		_, ok := runningHosts[hostName]
		if ok {
			delete(runningHosts, hostName)
		}
		if _, ok := output[hostName]; !ok {
			output[hostName] = `{
				"req_type": "op",
				"data":	"Server ` + hostName + ` is dead"
			}`
		}
	}
}
