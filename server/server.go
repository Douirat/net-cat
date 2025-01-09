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
	Time       string
	CilentName string
	Content    string
}

// Declare a structure to represent a server:
type Server struct {
	Clients  []*Client
	Messages []*Message
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
	if len(server.Clients) == 10 {
		server.HandleResponse(client.Connection, "The chat box is full!\n")
	}
	server.Clients = append(server.Clients, client)
	// Read data from the client
	reader := bufio.NewReader(conn)
	for {
		if client.Name == "" {
			server.HandleResponse(client.Connection, "[ENTER YOUR NAME]:")
		}
		res, err := reader.ReadString('\n')
		if err != nil {
			for _, v := range server.Clients {
				if v.Connection != client.Connection {
					server.HandleResponse(v.Connection, client.Name[:len(client.Name)-1]+" has left our chat...\n")
				}
			}
			return
		}
		if client.Name == "" && res != "" {
			client.Name = res
			for _, v := range server.Clients {
				if v.Connection != client.Connection {
					server.HandleResponse(v.Connection, client.Name[:len(client.Name)-1]+" has joined our chat...\n")
				}
			}
		} else {
			if res != "" {
				message := &Message{}
				message.Time = time.Now().Format("2006-01-02 15:04:05")
				message.CilentName = client.Name
				message.Content = res
				server.Messages = append(server.Messages, message)
				server.Broadcast(client.Connection)
			} else {
			}
			// Respond to the client
			server.HandleResponse(client.Connection, "Message received\n")
		}
	}
}

// Helper function to broadcast to the network:
func (server *Server) Broadcast(con net.Conn) {
	for _, client := range server.Clients {
		if client.Connection != con {
			for _, message := range server.Messages {
				response := "[" + message.Time + "]" + "[" + message.CilentName + "]: " + message.Content
				_, err := con.Write([]byte(response))
				if err != nil {
					fmt.Println("Error sending response:", err)
					return
				}
			}
		}
	}
}

// Helper function to broadcast error messages:
func (server *Server) HandleResponse(con net.Conn, str string) {
	_, err := con.Write([]byte(str))
	if err != nil {
		fmt.Println("Error sending response:", err)
		return
	}
}
