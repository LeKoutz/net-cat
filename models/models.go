package models

import (
	"net"
	"sync"
)

type Server struct {
	clients map[string]*Client
	mutex	sync.Mutex
}

type Client struct {
	conn net.Conn
	name string
}
