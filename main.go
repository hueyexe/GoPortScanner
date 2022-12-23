package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Get the hostname and port range from the user
	fmt.Print("Enter hostname: ")
	var hostname string
	fmt.Scanln(&hostname)

	fmt.Print("Enter start of port range: ")
	var startPortStr string
	fmt.Scanln(&startPortStr)
	startPort, err := strconv.Atoi(startPortStr)
	if err != nil {
		fmt.Println("Invalid port number:", startPortStr)
		os.Exit(1)
	}

	fmt.Print("Enter end of port range: ")
	var endPortStr string
	fmt.Scanln(&endPortStr)
	endPort, err := strconv.Atoi(endPortStr)
	if err != nil {
		fmt.Println("Invalid port number:", endPortStr)
		os.Exit(1)
	}

	// Validate the port range
	if startPort > endPort {
		fmt.Println("Invalid port range:", startPortStr, "-", endPortStr)
		os.Exit(1)
	}

	// Set up a wait group to track the goroutines
	var wg sync.WaitGroup

	// Set up counters for open and closed ports
	var openPorts, closedPorts int

	// Start a timer
	startTime := time.Now()

	// Scan the ports
	fmt.Println("Scanning ports", startPort, "to", endPort, "on", hostname)
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := hostname + ":" + strconv.Itoa(port)
			conn, err := net.Dial("tcp", address)
			if err == nil {
				conn.Close()
				fmt.Println(address, "is open")
				openPorts++
			} else {
				closedPorts++
			}
		}(port)
	}

	// Wait for all the goroutines to complete
	wg.Wait()

	// Stop the timer
	elapsedTime := time.Since(startTime)

	// Print the results
	fmt.Println("Scan complete.")
	fmt.Println("Open ports:", openPorts)
	fmt.Println("Closed ports:", closedPorts)
	fmt.Println("Total time elapsed:", elapsedTime)
}
