package utils

// Write a function that greps a file for a given pattern.
// The function should return the lines that match the pattern.

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func Grep(pattern, filename string) ([]string, error) {
	pattern = `"` + pattern + `"`
	pattern = strings.ReplaceAll(pattern, "|", `\|`)
	cmd := exec.Command("grep", pattern, filename, "-rn")
	// Print cmd structure	
	fmt.Println(cmd.Path + " " + strings.Join(cmd.Args[1:], " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StdoutPipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Error starting command: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	err = cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// grep returns exit status 1 if no lines were matched
			if exitErr.ExitCode() != 1 {
				return nil, fmt.Errorf("Command finished with error: %w", err)
			}
		} else {
			return nil, fmt.Errorf("Error waiting for command: %w", err)
		}
	}
	fmt.Println("lines:", lines)
	return lines, nil
}

func ReturnOutput(ip string, data string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":1200")

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
	//fmt.Print([]byte(data))
	// Send a message to the server
	_, err = conn.Write([]byte(`{
		"req_type": "op",
		"data" : "` + data + `"
	}` + "\n\r"))
	fmt.Println("send...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn.Write([]byte("\r"))
}
