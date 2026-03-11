package client

import (
	"fmt"
	"net"
)

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

func sendWelcome(conection net.Conn) {
	fmt.Fprint(conection, "Welcome to TCP-Chat!\n", Logo, "\n[ENTER YOUR NAME]:")
}
