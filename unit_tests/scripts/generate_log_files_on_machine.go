package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/go-faker/faker/v4"
)

type FakeDataStruct struct {
	Latitude    float32 `faker:"lat"`
	Longitude   float32 `faker:"long"`
	PhoneNumber string  `faker:"phone_number"`
	MacAddress  string  `faker:"mac_address"`
	UUID        string  `faker:"uuid_digit"`
}

func writeLinesToFile(lines []string, machineIdx string) {
	file, err := os.Create("test_log" + machineIdx + ".log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write lines to the file
	for _, line := range lines {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("Lines successfully written to file")
}

func generateFakeData(index string) {
	pattern := "PUT"
	filename := "/root/vm" + index + ".log"
	cmd_ := "grep -Ern " + pattern + " " + filename
	cmd := exec.Command("sh", "-c", cmd_)
	//fmt.Println(cmd)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// // Create a scanner
	scanner := bufio.NewScanner(stdout)

	lines := []string{}
	for i := 0; i < 50000; i++ {
		a := FakeDataStruct{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}
		lines = append(lines, fmt.Sprintf("%v", a))
	}

	// Read line by line
	i := 0
	specific_line := true
	output_lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		output_lines = append(output_lines, line)
		if i < 1000 {
			if index == "2" && specific_line {
				output_lines = append(output_lines, lines[i]+"GET 700")
				specific_line = false
			} else {
				output_lines = append(output_lines, lines[i])
			}
			i += 1
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	writeLinesToFile(output_lines, os.Args[1])
}

func main() {
	generateFakeData(os.Args[1])
}
