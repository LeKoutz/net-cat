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
	fmt.Printf("Listening on the port :%s\n", port)

}
