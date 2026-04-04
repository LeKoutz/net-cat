# net-cat
A TCP-based group chat server written in Go, modelled after the NetCat (`nc`) utility.

## Overview

Net-Cat runs as a TCP server that accepts multiple client connections. Clients connect using the standard `nc` command and can send messages to the group in real time. The server handles join/leave notifications, message history, and concurrent connections safely.

## Usage

```console
$ ./TCPChat                        # Listens on default port 8989
$ ./TCPChat 2525                   # Listens on selected port 2525. Accepts ports 0-65535
```

## Connecting as a Client

```console
$ nc <host-ip> <port>
```

## Features

- TCP server accepting up to 10 simultaneous clients
- Messages are timestamped and identified by sender: `[2020-01-20 15:48:41][name]:message`
- New clients receive full message history on join
- All clients are notified when someone joins or leaves
- Empty messages are not broadcast
- Client disconnection does not affect other clients
- Concurrent connections handled safely with goroutines, channels and mutexes

## Connecting from Another Machine

All machines must be on the same local network. Find the server machine's local IP with:

```console
$ ip a
```

Then connect from another machine:

```console
$ nc 192.168.x.x <port>
```

## Testing the program locally

Start the server on localhost
```console
./TCPChat localhost <port>
```

Simulate different clients by opening different terminals
```console
nc localhost <port>
```
Send messages from one client terminal to the other.