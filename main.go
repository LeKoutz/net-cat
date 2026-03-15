package main

import (
	"fmt"
	"net-cat/server"
	"os"
)

func main() {
	// Parse command line arguments to get the port number
	port, err := server.ParseArgs(os.Args[1:])
	if err != nil {
		server.PrintUsageMessage()
		os.Exit(1)
	}

	// Start server on port
	listener, err := server.StartServer(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	// Accept incoming connections
	err = server.AcceptConnections(listener)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accepting connections: %v\n", err)
		os.Exit(1)
	}

}
