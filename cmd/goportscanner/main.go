package main

import (
	"fmt"
	"os"

	"github.com/rancmd/goportscanner/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	hostname  string
	startPort int
	endPort   int
	timeout   string
	workers   int
	output    string
	format    string
	verbose   bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goportscanner",
		Short: "A fast and efficient port scanner written in Go",
		Long: `GoPortScanner is a high-performance port scanner designed for security testing and network reconnaissance.

Features:
- Fast concurrent scanning with configurable worker threads
- Multiple output formats (text, JSON, CSV)
- Banner grabbing for service identification
- Rate limiting to avoid overwhelming targets
- Comprehensive error handling and logging

Use responsibly and only scan systems you own or have permission to test.`,
		RunE: runScanner,
	}

	// Add flags
	rootCmd.Flags().StringVarP(&hostname, "hostname", "h", "", "Target hostname or IP address (required)")
	rootCmd.Flags().IntVarP(&startPort, "start-port", "s", 1, "Start of port range")
	rootCmd.Flags().IntVarP(&endPort, "end-port", "e", 1024, "End of port range")
	rootCmd.Flags().StringVarP(&timeout, "timeout", "t", "1s", "Connection timeout (e.g., 1s, 500ms)")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", 100, "Number of concurrent workers")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json, csv)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Mark required flags
	rootCmd.MarkFlagRequired("hostname")

	// Add examples
	rootCmd.Example = `  # Scan common ports on localhost
  goportscanner -h localhost -s 1 -e 1024

  # Scan specific ports with JSON output
  goportscanner -h scanme.nmap.org -s 20 -e 25 -f json

  # Fast scan with more workers
  goportscanner -h example.com -w 500 -t 500ms

  # Save results to file
  goportscanner -h target.com -o results.txt`

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runScanner(cmd *cobra.Command, args []string) error {
	// Validate inputs
	if hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	if startPort < 1 || startPort > 65535 {
		return fmt.Errorf("start port must be between 1 and 65535")
	}

	if endPort < 1 || endPort > 65535 {
		return fmt.Errorf("end port must be between 1 and 65535")
	}

	if startPort > endPort {
		return fmt.Errorf("start port cannot be greater than end port")
	}

	if workers < 1 {
		return fmt.Errorf("workers must be at least 1")
	}

	// Create scanner configuration
	config := scanner.Config{
		Hostname:  hostname,
		StartPort: startPort,
		EndPort:   endPort,
		Timeout:   timeout,
		Workers:   workers,
		Output:    output,
		Format:    format,
		Verbose:   verbose,
	}

	// Run the scanner
	return scanner.Run(config)
}