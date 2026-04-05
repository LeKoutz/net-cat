package server

import (
	"net"
	"strings"
	"testing"
	"time"
)

func newTestServer() *Server {
	return &Server{
		Clients: make(map[string]*Client),
	}
}

func TestAddClient(t *testing.T) {
	tests := []struct {
		name         string
		clientName   string
		setup        func(s *Server)
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid client",
			clientName:  "User",
			setup:       func(s *Server) {},
			expectError: false,
		},
		{
			name:         "Empty name",
			clientName:   "",
			setup:        func(s *Server) {},
			expectError:  true,
			errorMessage: "client name cannot be empty",
		},
		{
			name:       "Duplicate client",
			clientName: "User",
			setup: func(s *Server) {
				s.Clients["User"] = &Client{Name: "User"}
			},
			expectError:  true,
			errorMessage: "client name %s already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := newTestServer()

			// setup initial state
			tt.setup(testServer)

			// Print clients map before adding client for debugging
			t.Logf("Clients map before: %v", testServer.Clients)

			cl := &Client{Name: tt.clientName}
			err := testServer.AddClient(cl)

			// Print clients map for debugging
			t.Logf("Clients map after: %+v", testServer.Clients)
			if err != nil {
				t.Logf("Error message: %v", err)
			}

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBroadcasting(t *testing.T) {
	s := &Server{
		Clients:     make(map[string]*Client),
		broadcastCh: make(chan Message, 15),
	}
	connA1, _ := net.Pipe()
	connB1, connB2 := net.Pipe()
	sender := &Client{Name: "Alice", Conn: connA1}
	reciever := &Client{Name: "Bob", Conn: connB1}
	s.Clients["Alice"] = sender
	s.Clients["Bob"] = reciever
	go s.StartBroadcastingService()
	s.broadcastCh <- Message{Sender: sender, Content: "hello", Time: time.Now()}
	buf := make([]byte, 1024)
	n, _ := connB2.Read(buf)
	got := string(buf[:n])

	if !strings.Contains(got, "hello") {
		t.Fatalf("expected 'hello', got : %v", got)
	}
}
