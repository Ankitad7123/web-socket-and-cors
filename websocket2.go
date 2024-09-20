package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


first reading from the client then writing to other clients 


// Upgrader handles WebSocket connection upgrades
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Clients to manage connected WebSocket clients
var clients = make(map[*websocket.Conn]bool)

// Message struct for storing chat messages
type Message struct {
	gorm.Model
	Content string
}




// WebSocket handler
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		// Read a message from the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}

		// Broadcast message to all clients
		broadcastMessage(msg)

		// Save message to the database
		db.Create(&Message{Content: string(msg)})
	}
}

// Function to broadcast message to all clients
func broadcastMessage(msg []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

// Main function to set up Gin routes and WebSocket handling
func main() {
	// Initialize database
	InitializeDatabase()

	// Gin router setup
	r := gin.Default()

	// WebSocket route
	r.GET("/ws", handleWebSocket)

	// Run the Gin server on port 8080
	r.Run(":8080")
}
