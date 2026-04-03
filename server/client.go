package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Conn net.Conn
	Name string
}

type Message struct {
	Sender  *Client
	Content string
	Time    time.Time
}

const Logo = `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | ` + "`" + `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'

 `

// Welcome message with logo
func sendWelcome(conn net.Conn) {
	fmt.Fprint(conn, "Welcome to TCP-Chat!\n", Logo, "\n[ENTER YOUR NAME]:")

}

// Get users name
func GetName(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		name := scanner.Text()
		if name == "" {
			fmt.Fprint(conn, "[ENTER YOUR NAME]:")
			continue
		}
		return name
	}
	return ""
}

// Timestamp formation
func (m *Message) Format() string {
	formattedt := m.Time.Format("2006-01-02 15:04:05")
	return "[" + formattedt + "][" + m.Sender.Name + "]:" + m.Content
}
