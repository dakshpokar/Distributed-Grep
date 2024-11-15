package utils

// Write a function that greps a file for a given pattern.
// The function should return the lines that match the pattern.

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
)

func Grep(pattern, filename string) ([]string, error) {
	pattern = `"` + pattern + `"`
	cmd_ := "grep -Ern " + pattern + " " + filename
	cmd := exec.Command("sh", "-c", cmd_)
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
	count := strconv.Itoa(len(lines))
	lines = append(lines, count)
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
	WriteOutputToFile(lines, pattern)
	return []string{count}, nil
}

func WriteOutputToFile(lines []string, pattern string) {
	filenm := "output_" + pattern + ".txt"
	file, err := os.Create(filenm)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, line := range lines {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	return

}

func ReturnOutput(conn net.Conn, data string) {
	_, err := conn.Write([]byte(`{
		"req_type": "op",
		"data" : "` + data + `"
	}` + "\n\r"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn.Write([]byte("\r"))
}
