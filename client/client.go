package client

import (
	"bufio"
	"fmt"
	"net"
)

// Client struct stores connection and Username
type Client struct {
	Conn net.Conn // η "γραμμή" σύνδεσης TCP
	Name string   // το όνομά του
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

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	sendWelcome(conn)
	GetName(conn)
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
