package server

import (
	"fmt"
	"os"
	"net"
	"net-cat/client"
	"sync"
)

type Server struct {
	clients map[string]*client.Client
	mutex	sync.Mutex
}

// ParseArgs parses command line arguments and returns the port as a string.
// If no arguments are provided, it returns the default port "8989".
// It returns an error if the number of arguments are more than one or if the provided port number is invalid.
func ParseArgs(args []string) (string, error) {
	defaultPort := "8989"
	if len(args) == 0 {
		return defaultPort, nil
	}
	if len(args) > 1 {
		return "", fmt.Errorf("invalid number of arguments")
	}
	port := args[0]
	portIsValid := ValidatePort(port)
	if len(args) == 1 && portIsValid{
		return port, nil
	} else {
		return "", fmt.Errorf("invalid port number")
	}
}

// PrintUsageMessage prints the usage message "<USAGE>: ./TCPChat $port" to stderr.
func PrintUsageMessage() {
	fmt.Fprintln(os.Stderr, "[USAGE]: ./TCPChat $port")
}

// MyAtoi is similar to strconv.Atoi, but it is implemented manually because strconv package is not allowed in this project.
// Converts a string to an integer. It returns an error if the string is not a valid integer.
func MyAtoi(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("invalid input")
	}
	num := 0
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("invalid input")
		}
		num = num*10 + int(char-'0')
	}
	return num, nil
}

// ValidatePort checks if the provided port number is within the valid range (0-65535).
func ValidatePort(port string) bool {
	portNum, err := MyAtoi(port)
	if err != nil {
		return false
	}
	if 0 <= portNum && portNum <= 65535 {
		return true
	} else {
		return false
	}
}

// StartServer starts a TCP server that listens on the specified port.
// Prints the message "Listening on the port :$port" to stdout when the server starts successfully.
// It returns an error if there is an issue starting the server.
func StartServer(port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %v", err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	portNum := addr.Port
	fmt.Printf("Listening on the port :%d\n", portNum)

	return listener, nil
}

// AcceptConnections accepts incoming connections in an infinite loop. It returns an error if there is an issue accepting connections.
func AcceptConnections(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		}
		// Handle Connection
		go client.HandleConnection(conn)
	}
}
