package server

import (
	"fmt"
	"os"
)

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
