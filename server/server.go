package server

import (
	// "io"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
	Max      int
}

// Instantiate  a new server:
func NewServer() *Server {
	server := new(Server)
	server.Clients = []*Client{}
	server.Messages = []*Message{}
	server.Max = 10
	return server
}

func ListenAndServe() {
	args := os.Args[1:]
	Port := "8989"
	if len(args) == 1 && IsValidPort(args[0]) {
		Port = args[0]
	}
	fmt.Println("localhost:"+Port)
	ServerInst := NewServer()
	listener, err := net.Listen("tcp", "localhost:"+Port)
	fmt.Printf("Listening on the port :%v\n", Port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
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
	wellcome := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]:"

	for {
		if client.Name == "" {
			server.HandleResponse(client.Connection, wellcome)
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
				_, err := client.Connection.Write([]byte(response))
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

// IsValidPort checks if the given port string is a valid port number.
func IsValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr) // Convert the string to an integer.
	if err != nil {
		return false // Not a valid number.
	}
	return port > 0 && port <= 65535 // Valid port range.
}
