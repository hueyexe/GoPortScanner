package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Define the command-line flags
	var hostname string
	var startPort, endPort int
	var timeout time.Duration
	flag.StringVar(&hostname, "hostname", "", "The hostname or IP address to scan")
	flag.IntVar(&startPort, "start-port", 1, "The start of the port range to scan")
	flag.IntVar(&endPort, "end-port", 65535, "The end of the port range to scan")
	flag.DurationVar(&timeout, "timeout", time.Second, "The timeout for connection attempts")
	flag.Parse()

	// Validate the port range
	if startPort > endPort {
		fmt.Println("Invalid port range:", startPort, "-", endPort)
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
			conn, err := net.DialTimeout("tcp", address, timeout)
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
