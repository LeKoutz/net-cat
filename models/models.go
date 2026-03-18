package models

import (
	"net"
	"sync"
)

type Server struct {
	Listener 	net.Listener
	Clients 	map[string]*Client
	MaxClients	int
	Mutex		sync.Mutex
}

type Client struct {
	Conn net.Conn
	Name string
}
