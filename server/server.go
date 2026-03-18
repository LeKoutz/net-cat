package server

import (
	"fmt"
	"os"
	"net"
	"net-cat/client"
	"net-cat/models"
	"sync"
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

// StartServer starts a TCP server that listens on the specified port.
// Prints the message "Listening on the port :$port" to stdout when the server starts successfully.
// It returns an error if there is an issue starting the server.
func StartServer(port string) (*models.Server, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %v", err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	portNum := addr.Port
	fmt.Printf("Listening on the port :%d\n", portNum)

	server := &models.Server{
		Listener:   listener,
		Clients:    make(map[string]*models.Client),
		MaxClients: 10,
		Mutex:      sync.Mutex{},
	}

	return server, nil
}

// AcceptConnections accepts incoming connections in an infinite loop. It returns an error if there is an issue accepting connections.
func AcceptConnections(server *models.Server) error {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		}
		// Check if the number of clients has reached the maximum limit
		server.Mutex.Lock()
		if len(server.Clients) >= server.MaxClients {
			server.Mutex.Unlock()
			conn.Close()
			fmt.Println("Maximum number of clients reached. Cannot accept more connections.")
			continue
		}
		server.Mutex.Unlock()

		// Handle Connection
		go client.HandleConnection(conn)
	}
}

// AddClient adds a new client to the server's clients map. It locks the mutex to ensure thread safety while modifying the clients map.
func AddClient(s *models.Server, cl *models.Client) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if cl.Name == "" {
		return fmt.Errorf("Client name cannot be empty")
	}
	if _, exists := s.Clients[cl.Name]; exists {
		return fmt.Errorf("Client name \"%s\" already exists", cl.Name)
	}
	s.Clients[cl.Name] = cl
	return nil
}