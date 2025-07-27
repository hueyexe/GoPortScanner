package scanner

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config holds the scanner configuration
type Config struct {
	Hostname  string
	StartPort int
	EndPort   int
	Timeout   string
	Workers   int
	Output    string
	Format    string
	Verbose   bool
}

// Result represents a scan result for a single port
type Result struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Status   string `json:"status"`
	Service  string `json:"service,omitempty"`
	Banner   string `json:"banner,omitempty"`
	Error    string `json:"error,omitempty"`
}

// ScanSummary contains the overall scan results
type ScanSummary struct {
	Hostname    string    `json:"hostname"`
	StartPort   int       `json:"start_port"`
	EndPort     int       `json:"end_port"`
	OpenPorts   int       `json:"open_ports"`
	ClosedPorts int       `json:"closed_ports"`
	TotalTime   string    `json:"total_time"`
	Results     []Result  `json:"results"`
}

// Scanner represents the port scanner
type Scanner struct {
	config  Config
	results []Result
	mutex   sync.RWMutex
	summary ScanSummary
}

// NewScanner creates a new scanner instance
func NewScanner(config Config) *Scanner {
	return &Scanner{
		config: config,
		results: make([]Result, 0),
	}
}

// Run executes the port scan with the given configuration
func Run(config Config) error {
	scanner := NewScanner(config)
	return scanner.Scan()
}

// Scan performs the actual port scanning
func (s *Scanner) Scan() error {
	// Parse timeout
	timeout, err := time.ParseDuration(s.config.Timeout)
	if err != nil {
		return fmt.Errorf("invalid timeout format: %v", err)
	}

	if s.config.Verbose {
		fmt.Printf("Starting scan of %s (ports %d-%d) with %d workers\n", 
			s.config.Hostname, s.config.StartPort, s.config.EndPort, s.config.Workers)
	}

	// Initialize summary
	s.summary = ScanSummary{
		Hostname:  s.config.Hostname,
		StartPort: s.config.StartPort,
		EndPort:   s.config.EndPort,
	}

	startTime := time.Now()

	// Create worker pool
	portChan := make(chan int, s.config.EndPort-s.config.StartPort+1)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go s.worker(portChan, &wg, timeout)
	}

	// Send ports to workers
	go func() {
		defer close(portChan)
		for port := s.config.StartPort; port <= s.config.EndPort; port++ {
			portChan <- port
		}
	}()

	// Wait for all workers to complete
	wg.Wait()

	// Calculate summary
	elapsedTime := time.Since(startTime)
	s.summary.TotalTime = elapsedTime.String()

	// Count results
	for _, result := range s.results {
		if result.Status == "open" {
			s.summary.OpenPorts++
		} else {
			s.summary.ClosedPorts++
		}
	}

	// Output results
	return s.outputResults()
}

// worker processes ports from the channel
func (s *Scanner) worker(portChan <-chan int, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()

	for port := range portChan {
		result := s.scanPort(port, timeout)
		s.addResult(result)
	}
}

// scanPort scans a single port
func (s *Scanner) scanPort(port int, timeout time.Duration) Result {
	result := Result{
		Hostname: s.config.Hostname,
		Port:     port,
	}

	address := fmt.Sprintf("%s:%d", s.config.Hostname, port)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create dialer with timeout
	dialer := net.Dialer{
		Timeout: timeout,
	}

	// Attempt connection
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		result.Status = "closed"
		result.Error = err.Error()
		return result
	}
	defer conn.Close()

	// Port is open
	result.Status = "open"

	// Try to get banner
	banner, err := s.getBanner(conn, timeout)
	if err == nil && banner != "" {
		result.Banner = strings.TrimSpace(banner)
		result.Service = s.identifyService(port, banner)
	}

	return result
}

// getBanner attempts to read a banner from the connection
func (s *Scanner) getBanner(conn net.Conn, timeout time.Duration) (string, error) {
	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(timeout))

	// Try to read banner
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return banner, nil
}

// identifyService attempts to identify the service based on port and banner
func (s *Scanner) identifyService(port int, banner string) string {
	// Common port mappings
	portServices := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		993:  "IMAPS",
		995:  "POP3S",
		3306: "MySQL",
		5432: "PostgreSQL",
		6379: "Redis",
		8080: "HTTP-Alt",
		8443: "HTTPS-Alt",
	}

	// Check if we have a known port
	if service, exists := portServices[port]; exists {
		return service
	}

	// Try to identify from banner
	banner = strings.ToLower(banner)
	switch {
	case strings.Contains(banner, "ssh"):
		return "SSH"
	case strings.Contains(banner, "http"):
		return "HTTP"
	case strings.Contains(banner, "ftp"):
		return "FTP"
	case strings.Contains(banner, "smtp"):
		return "SMTP"
	case strings.Contains(banner, "pop3"):
		return "POP3"
	case strings.Contains(banner, "imap"):
		return "IMAP"
	case strings.Contains(banner, "mysql"):
		return "MySQL"
	case strings.Contains(banner, "postgresql"):
		return "PostgreSQL"
	case strings.Contains(banner, "redis"):
		return "Redis"
	}

	return "Unknown"
}

// addResult adds a result to the scanner in a thread-safe manner
func (s *Scanner) addResult(result Result) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.results = append(s.results, result)
}

// outputResults outputs the scan results in the specified format
func (s *Scanner) outputResults() error {
	var output *os.File
	var err error

	// Determine output destination
	if s.config.Output == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(s.config.Output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer output.Close()
	}

	// Output based on format
	switch strings.ToLower(s.config.Format) {
	case "json":
		return s.outputJSON(output)
	case "csv":
		return s.outputCSV(output)
	case "text":
		fallthrough
	default:
		return s.outputText(output)
	}
}

// outputJSON outputs results in JSON format
func (s *Scanner) outputJSON(output *os.File) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.summary)
}

// outputCSV outputs results in CSV format
func (s *Scanner) outputCSV(output *os.File) error {
	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write header
	header := []string{"Hostname", "Port", "Status", "Service", "Banner", "Error"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, result := range s.results {
		row := []string{
			result.Hostname,
			strconv.Itoa(result.Port),
			result.Status,
			result.Service,
			result.Banner,
			result.Error,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// outputText outputs results in human-readable text format
func (s *Scanner) outputText(output *os.File) error {
	// Print summary first
	fmt.Fprintf(output, "Scan Summary\n")
	fmt.Fprintf(output, "============\n")
	fmt.Fprintf(output, "Target: %s\n", s.summary.Hostname)
	fmt.Fprintf(output, "Port Range: %d-%d\n", s.summary.StartPort, s.summary.EndPort)
	fmt.Fprintf(output, "Open Ports: %d\n", s.summary.OpenPorts)
	fmt.Fprintf(output, "Closed Ports: %d\n", s.summary.ClosedPorts)
	fmt.Fprintf(output, "Total Time: %s\n\n", s.summary.TotalTime)

	// Print open ports
	if s.summary.OpenPorts > 0 {
		fmt.Fprintf(output, "Open Ports:\n")
		fmt.Fprintf(output, "===========\n")
		for _, result := range s.results {
			if result.Status == "open" {
				if result.Service != "" {
					fmt.Fprintf(output, "%s:%d (%s)", result.Hostname, result.Port, result.Service)
				} else {
					fmt.Fprintf(output, "%s:%d", result.Hostname, result.Port)
				}
				if result.Banner != "" {
					fmt.Fprintf(output, " - %s", result.Banner)
				}
				fmt.Fprintf(output, "\n")
			}
		}
	} else {
		fmt.Fprintf(output, "No open ports found.\n")
	}

	return nil
}