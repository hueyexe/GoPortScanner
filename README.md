# GoPortScanner
This is a simple tool that allows you to scan a range of ports on a specified host to see which ones are open. It is written in Go, a fast and efficient programming language that is easy to learn (and fun to use).

## Prerequisites
Before you can use the port scanner, you will need to have Go installed on your computer. If you don't already have Go installed, you can download it from the Go website. Follow the installation instructions for your operating system.

## Usage
To use the port scanner, follow these steps:
1. Download or clone the source code for the port scanner from GitHub.
2. Open a terminal or command prompt, and navigate to the directory where you downloaded the source code.
3. Run the following command to build the port scanner: `go build`
4. Run the port scanner by typing ./port-scanner (on Linux or macOS) or port-scanner.exe (on Windows).
5. When prompted, enter the hostname or IP address of the host you want to scan.
6. When prompted, enter the start and end of the port range you want to scan.
7. The port scanner will scan the specified ports and print a message for each open port it finds. When the scan is complete, it will display the total number of open and closed ports, and the time elapsed.

## Example
Here is an example of the scanner in action:
```
$ ./port-scanner
Enter hostname: example.com
Enter start of port range: 1
Enter end of port range: 100
Scanning ports 1 to 100 on example.com
example.com:22 is open
example.com:80 is open
Scan complete.
Open ports: 2
Closed ports: 98
Total time elapsed: 2.475722ms
```
In this example, the program scanned ports 1 to 100 on the host example.com, and found that ports 22 (SSH) and 80 (HTTP) are open.

## Permission first
The program is designed to be fast and efficient, but be aware that port scanning can generate a lot of network traffic and may be perceived as a malicious activity. Use caution when running the port scanner, and obtain permission before scanning any hosts that you do not own, control or otherwise have permission to scan.
