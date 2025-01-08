package server

import (
	// "io"
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// Declare a structure to represent a client:
type Client struct {
	Connection net.Conn
	Name       string
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
		go handleConnection(conn, ServerInst) // handle one connection at a time
	}
}

func handleConnection(conn net.Conn, server *Server) {
	defer conn.Close()
	client := &Client{}
	client.Connection = conn
	server.Clients = append(server.Clients, client)
	// Read data from the client
	reader := bufio.NewReader(conn)
	for {
		if client.Name == "" {
			_, err := conn.Write([]byte("[ENTER YOUR NAME]:"))
			if err != nil {
				fmt.Println("Error sending response:", err)
				return
			}
		}
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by client")
			return
		}
		if client.Name == "" && res != "" {
			client.Name = res
		} else {
			message := &Message{}
			message.Time = time.Now()
			// Format("2006-01-02 15:04:05")
			message.CilentName = client.Name
			message.Content = res
			_, err = conn.Write([]byte(""\n"))
			if err != nil {
				fmt.Println("Error sending response:", err)
				return
			}
			fmt.Printf("Message received: %s", res)
		}
		// Respond to the client
		_, err = conn.Write([]byte("Message received\n"))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}
