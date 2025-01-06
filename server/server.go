package server

import (
	"io"
	"log"
	"net"
	"time"
)

// Declare a structure to represent a client:
type Client struct {
	Name string
}

// Declare a structure to represent a message:
type Message struct {
	Time       time.Time
	CilentName string
	Content    string
}

// Declare a structure to represent a server:
type Server struct {
	Clients []*Client
	Message []*Message
}

// Instantiate  a new server:
func NewServer() *Server {
	server := new(Server)
	server.Clients = []*Client{}

	return server
}

func ListenAndServe() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		ServerInst := new(Server)
		handleConn(conn,  ServerInst) // handle one connection at a time
	}
}

func handleConn(c net.Conn, server *Server) {
	defer c.Close()
}
