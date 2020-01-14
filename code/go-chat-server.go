package main

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const Welcome = `
Welcome to go-chat!
Whatever you send will be broadcast to everyone, unless it's a command.
Commands:
	/nick <your name>
	/whoami
`

type Connection struct {
	Created time.Time
	Name    string
}

var connections map[net.Conn]*Connection
var connectionsMutex sync.Mutex

func main() {
	srv, err := net.Listen("tcp", ":5678")
	if err != nil {
		log.Fatalln(err)
	}
	defer srv.Close()

	connections = make(map[net.Conn]*Connection)

	log.Println("Listening for connections")
	for {
		c, err := srv.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		} else {
			log.Println(c.RemoteAddr())
			go handleConnection(c)
		}
	}
}

func handleConnection(c net.Conn) {
	log.Println(c.RemoteAddr(), " CONNECTED")

	connection := &Connection{
		Created: time.Now(),
		Name:    c.RemoteAddr().String(),
	}

	// Add connection to the pool
	connectionsMutex.Lock()
	connections[c] = connection
	connectionsMutex.Unlock()

	// Remove self when connection closes
	defer func() {
		log.Println(c.RemoteAddr(), " DISCONNECTED")
		connectionsMutex.Lock()
		delete(connections, c)
		connectionsMutex.Unlock()
	}()

	c.Write([]byte(Welcome))

	b := make([]byte, 4096)
	for {
		n, err := c.Read(b)
		if err != nil {
			return
		}
		in := strings.TrimSpace(string(b[:n]))
		if len(in) == 0 {
			continue
		}
		log.Println(c.RemoteAddr(), " '"+in+"'")
		if in[0] != '/' {
			broacast(connection.Name, in)
		} else { // It's a command
			command := strings.ToLower(in[1:])
			params := ""
			if spaceIdx := strings.Index(in, " "); spaceIdx > -1 {
				command = strings.ToLower(in[1:spaceIdx])
				if spaceIdx < len(in) {
					params = in[spaceIdx+1:]
				}
			}

			switch command {
			case "nick":
				connection.Name = params
				broacast(connection.Name, "has set their nick name")

			case "whoami":
				c.Write([]byte("You are " + connection.Name + "\n"))

			default:
				c.Write([]byte("unknown command: " + command + "\n"))
			}
		}
	}
}

func broacast(sender string, message string) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	b := []byte("[" + sender + "] " + message + "\n")

	for c := range connections {
		c.Write(b)
	}
}
