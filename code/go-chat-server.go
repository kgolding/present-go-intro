package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const Welcome = `
Welcome to go-chat!
Whatever you send will be broadcast to everyone, unless it's a command.
Commands:
	/nick <your name>
	/whoami
`

type Server struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

type client struct {
	conn    net.Conn
	Name    string
	Created time.Time
}

func (s *Server) Run() {
	go func() {
		for {
			select {
			case message := <-s.forward:
				for c := range s.clients {
					c.conn.Write(message)
				}

			case c := <-s.leave:
				delete(s.clients, c)

			case c := <-s.join:
				s.clients[c] = true
			}
		}
	}()
}

func (s *Server) Join(client *client) {
	s.join <- client
}

func (s *Server) Leave(client *client) {
	s.leave <- client
}

func (s *Server) Broadcast(sender string, message string) {
	b := []byte("[" + sender + "] " + message + "\n")
	s.forward <- b
}

var server Server

func main() {
	conn, err := net.Listen("tcp", ":5678")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	server = Server{
		forward: make(chan []byte, 10),
		join:    make(chan *client, 10),
		clients: make(map[*client]bool),
	}

	log.Println("Listening for connections")
	server.Run()
	for {
		c, err := conn.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		} else {
			go handleConnection(c)
		}
	}
}

func handleConnection(conn net.Conn) {
	log.Println(conn.RemoteAddr(), " CONNECTED")

	// Create the client
	client := &client{
		conn:    conn,
		Created: time.Now(),
		Name:    conn.RemoteAddr().String(),
	}

	// Add client to the pool
	server.join <- client

	// Remove self when connection closes
	defer func() {
		log.Println(conn.RemoteAddr(), " DISCONNECTED")
		server.leave <- client
	}()

	conn.Write([]byte(Welcome))

	b := make([]byte, 4096)
	for {
		n, err := conn.Read(b)
		if err != nil {
			return
		}
		in := strings.TrimSpace(string(b[:n]))
		if len(in) == 0 {
			continue
		}
		log.Println(conn.RemoteAddr(), " '"+in+"'")
		if in[0] != '/' {
			server.Broadcast(client.Name, in)
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
				client.Name = params
				conn.Write([]byte("Hence forth you shall be known as " + client.Name + "\n"))
				server.Broadcast(client.Name, "has set their nick name")

			case "whoami":
				conn.Write([]byte("You are " + client.Name + "\n"))

			default:
				conn.Write([]byte("Unknown command: " + command + "\n"))
			}
		}
	}
}
