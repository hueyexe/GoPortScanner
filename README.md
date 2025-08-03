# GoPortScanner

[![Go Report Card](https://goreportcard.com/badge/github.com/rancmd/goportscanner)](https://goreportcard.com/report/github.com/rancmd/goportscanner)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A fast, efficient, and feature-rich port scanner written in Go. Designed for security professionals, penetration testers, and network administrators who need reliable port scanning capabilities.

## Features

- ⚡ **High Performance**: Concurrent scanning with configurable worker threads
- 🔍 **Service Detection**: Automatic service identification and banner grabbing
- 📊 **Multiple Output Formats**: Text, JSON, and CSV output options
- 🎯 **Flexible Targeting**: Custom port ranges and timeout configuration
- 🛡️ **Responsible Design**: Built-in rate limiting and timeout controls
- 📝 **Comprehensive Logging**: Verbose mode for detailed scan information
- 🔧 **Easy Integration**: Simple API for programmatic use

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Network access to target systems

### Installation

#### From Source
```bash
git clone https://github.com/rancmd/goportscanner.git
cd goportscanner
go build -o goportscanner cmd/main.go
```

#### Using Go Install
```bash
go install github.com/rancmd/goportscanner/cmd/goportscanner@latest
```

### Basic Usage

```bash
# Scan common ports on localhost
goportscanner -H localhost -s 1 -e 1024

# Scan specific ports with JSON output
goportscanner -H scanme.nmap.org -s 20 -e 25 -f json

# Fast scan with more workers
goportscanner -H example.com -w 500 -t 500ms

# Save results to file
goportscanner -H target.com -o results.txt
```

## Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--hostname` | `-H` | Target hostname or IP address (required) | - |
| `--start-port` | `-s` | Start of port range | 1 |
| `--end-port` | `-e` | End of port range | 1024 |
| `--timeout` | `-t` | Connection timeout (e.g., 1s, 500ms) | 1s |
| `--workers` | `-w` | Number of concurrent workers | 100 |
| `--output` | `-o` | Output file (default: stdout) | - |
| `--format` | `-f` | Output format (text, json, csv) | text |
| `--verbose` | `-v` | Verbose output | false |

## Examples

### Basic Port Scan
```bash
goportscanner -H scanme.nmap.org -s 1 -e 100
```

**Output:**
```
Scan Summary
============
Target: scanme.nmap.org
Port Range: 1-100
Open Ports: 2
Closed Ports: 98
Total Time: 2.3s

Open Ports:
===========
scanme.nmap.org:22 (SSH) - SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
scanme.nmap.org:80 (HTTP) - HTTP/1.1 200 OK
```

### JSON Output
```bash
goportscanner -H localhost -s 80 -e 90 -f json
```

**Output:**
```json
{
  "hostname": "localhost",
  "start_port": 80,
  "end_port": 90,
  "open_ports": 1,
  "closed_ports": 10,
  "total_time": "1.2s",
  "results": [
    {
      "hostname": "localhost",
      "port": 80,
      "status": "open",
      "service": "HTTP",
      "banner": "HTTP/1.1 200 OK"
    }
  ]
}
```

### CSV Output
```bash
goportscanner -H example.com -s 20 -e 25 -f csv -o scan_results.csv
```

### High-Speed Scanning
```bash
goportscanner -H target.com -w 1000 -t 200ms -s 1 -e 65535
```

## Service Detection

GoPortScanner automatically identifies common services based on port numbers and banner information:

- **Web Services**: HTTP (80, 8080), HTTPS (443, 8443)
- **Remote Access**: SSH (22), Telnet (23)
- **Email Services**: SMTP (25), POP3 (110), IMAP (143)
- **Database Services**: MySQL (3306), PostgreSQL (5432), Redis (6379)
- **DNS**: DNS (53)

## Performance Tuning

### Worker Threads
Adjust the number of workers based on your network capacity and target system:

```bash
# Conservative scanning (100 workers)
goportscanner -H target.com -w 100

# Aggressive scanning (1000 workers)
goportscanner -H target.com -w 1000
```

### Timeout Settings
Optimize timeout values for your network conditions:

```bash
# Fast local network
goportscanner -H localhost -t 100ms

# Slower internet connection
goportscanner -H remote.com -t 2s
```

## Security and Legal Considerations

⚠️ **IMPORTANT**: Port scanning can be considered a hostile activity by some systems and networks.

### Best Practices

1. **Always obtain permission** before scanning any system you don't own
2. **Use responsibly** - avoid overwhelming target systems
3. **Respect rate limits** - adjust worker count and timeouts appropriately
4. **Follow local laws** - ensure compliance with applicable regulations
5. **Document your activities** - keep records of authorized scans

### Legal Disclaimer

This tool is provided for educational and authorized security testing purposes only. Users are responsible for ensuring they have proper authorization before scanning any systems. **Unauthorized port scanning may violate local laws and regulations.** The authors disclaim all liability for any misuse of this software.

**By using this software, you agree to:**
- Only scan systems you own or have explicit written permission to test
- Comply with all applicable local, state, and federal laws
- Use reasonable rate limiting to avoid disrupting target systems
- Accept full responsibility for your actions

## Development

### Building from Source

```bash
git clone https://github.com/rancmd/goportscanner.git
cd goportscanner
go mod download
go build -o goportscanner cmd/main.go
```

### Running Tests

```bash
go test ./...
```

### Code Quality

```bash
# Run linter
golangci-lint run

# Run security checks
gosec ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by tools like `nmap` and `masscan`
- Built with Go standard library
- Uses [Cobra](https://github.com/spf13/cobra) for CLI functionality
