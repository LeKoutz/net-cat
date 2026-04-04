package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	Listener    net.Listener
	Clients     map[string]*Client
	MaxClients  int
	Mutex       sync.Mutex
	broadcastCh chan Message
	History     []Message
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
	if len(args) == 1 && portIsValid {
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
func StartServer(port string) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %v", err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	portNum := addr.Port
	fmt.Printf("Listening on the port :%d\n", portNum)

	server := &Server{
		Listener:    listener,
		Clients:     make(map[string]*Client),
		MaxClients:  10,
		Mutex:       sync.Mutex{},
		broadcastCh: make(chan Message, 15),
		History:     []Message{},
	}

	return server, nil
}

// AcceptConnections accepts incoming connections in an infinite loop. It returns an error if there is an issue accepting connections.
func (server *Server) AcceptConnections() error {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		}
		// Check if the number of clients has reached the maximum limit
		server.Mutex.Lock()
		if len(server.Clients) >= server.MaxClients {
			fmt.Fprintln(conn, "Maximum number of clients reached. Cannot accept more connections.")
			server.Mutex.Unlock()
			conn.Close()
			continue
		}
		server.Mutex.Unlock()

		// Handle Connection
		go server.HandleConnection(conn)
	}
}

// Adds-remove clients, sends the messages.
func (server *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	sendWelcome(conn)
	client := &Client{Name: GetName(conn), Conn: conn}
	mesg := Message{Sender: client, System: true}
	for {
		err := server.AddClient(client)
		if err == nil {
			for _, msg := range server.History {
				fmt.Fprintln(conn, msg.Format())
			}
			defer server.RemoveClient(client)
			mesg.Content = fmt.Sprintf("\n%v has joined our chat...", client.Name)
			fmt.Fprintf(client.Conn, "[%v][%v]:", time.Now().Format("2006-01-02 15:04:05"), client.Name)
			server.broadcastCh <- mesg
			break
		} else {
			fmt.Fprintf(conn, "%v\n", err)
			client = &Client{Name: GetName(conn), Conn: conn}
		}
	}
	// Listen for client input and send message to the server's broadcast channel
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "" {
			continue
		}
		server.broadcastCh <- Message{Sender: client, Content: msg, Time: time.Now()}
		fmt.Fprintf(client.Conn, "[%v][%v]:", time.Now().Format("2006-01-02 15:04:05"), client.Name)

	}
	mesg.Content = fmt.Sprintf("\n%v has left our chat...", client.Name)
	server.broadcastCh <- mesg

}

// AddClient adds a new client to the server's clients map. It locks the mutex to ensure thread safety while modifying the clients map.
func (server *Server) AddClient(cl *Client) error {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	if cl.Name == "" {
		return fmt.Errorf("Client name cannot be empty")
	}
	if _, exists := server.Clients[cl.Name]; exists {
		return fmt.Errorf("Client name \"%s\" already exists", cl.Name)
	}
	server.Clients[cl.Name] = cl
	return nil
}

// StartBroadcastingService listens for messages on the broadcast channel
// and sends them to all connected clients except the sender.
// It also keeps track of the message history.
func (server *Server) StartBroadcastingService() {
	for msg := range server.broadcastCh {
		if !msg.System {
			server.SaveMessageToHistory(msg)
		}
		server.Mutex.Lock()
		clients := make(map[string]*Client)
		for _, client := range server.Clients {
			if client.Name != msg.Sender.Name {
				clients[client.Name] = client
			}
		}
		server.Mutex.Unlock()
		for _, client := range clients {
			fmt.Fprintln(client.Conn, "\n"+msg.Format())
			fmt.Fprintf(client.Conn, "[%v][%v]:", time.Now().Format("2006-01-02 15:04:05"), client.Name)

		}
	}
}

func (server *Server) RemoveClient(cl *Client) error {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	delete(server.Clients, cl.Name)
	return nil
}

func (server *Server) SaveMessageToHistory(msg Message) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	server.History = append(server.History, msg)
}
