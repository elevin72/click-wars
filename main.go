package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn         *websocket.Conn
	LastActivity time.Time
}

var (
	clients      = make(map[*websocket.Conn]*Client)
	clientsMutex sync.Mutex
	linePosition int32 = 0
	lineMutex    sync.Mutex
	totalHits    int32 = 0
	hitsMutex    sync.Mutex
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Conn:         ws,
		LastActivity: time.Now(),
	}

	clientsMutex.Lock()
	log.Println("Adding client")
	clients[ws] = client
	clientsMutex.Unlock()

	defer func() {
		delete(clients, ws)
		ws.Close()
	}()

	initMessage := make([]byte, 5)
	initMessage[0] = 1
	binary.LittleEndian.PutUint32(initMessage[1:5], uint32(linePosition))
	err = client.Conn.WriteMessage(websocket.BinaryMessage, initMessage)
	if err != nil {
		log.Println(err)
	}

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Client unexpectedly closed")
				log.Print(err)
				clientsMutex.Lock()
				delete(clients, client.Conn)
				clientsMutex.Unlock()
			}
			return
		}

		if len(message) != 10 && message[0] != 0xFF {
			log.Println("Invalid message or ping")
			log.Println(string(message))
			continue
		}

		if message[0] == 0xFF {
			client.LastActivity = time.Now()
			continue
		}

		if len(message) != 10 {
			log.Println("Invalid message length")
			continue
		}

		x := int32(binary.LittleEndian.Uint32(message[1:5]))
		y := int32(binary.LittleEndian.Uint32(message[5:9]))
		color := message[9] // 0 for blue, 1 for red

		var lineIncrement int32
		if color == 0 {
			lineIncrement = 1
		} else { // Red side
			lineIncrement = -1
		}
		lineMutex.Lock()
		linePosition += lineIncrement
		lineMutex.Unlock()
		hitsMutex.Lock()
		totalHits++
		hitsMutex.Unlock()

		log.Printf("Received x: %d, y: %d, color: %d\n", x, y, color)
		log.Printf("New line position: %d\n", linePosition)

		clientsMutex.Lock()
		responseMessage := createUpdatedMessage(x, y, color, linePosition, totalHits)
		for _, c := range clients {
			// err := c.Conn.WriteMessage(websocket.BinaryMessage, responseMessage)
			// if err != nil {
			// 	log.Println(err)
			// }
			if c != client {
				err := c.Conn.WriteMessage(websocket.BinaryMessage, responseMessage)
				if err != nil {
					log.Println(err)
				}
			}
		}
		clientsMutex.Unlock()

	}
}

func createUpdatedMessage(x, y int32, color byte, linePosition, totalHits int32) []byte {
	msg := make([]byte, 18)
	msg[0] = 0
	binary.LittleEndian.PutUint32(msg[1:5], uint32(x))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(y))
	msg[9] = color
	binary.LittleEndian.PutUint32(msg[10:14], uint32(linePosition))
	binary.LittleEndian.PutUint32(msg[14:], uint32(totalHits))
	return msg
}

func closeStaleConnections() {
	clientsMutex.Lock()
	now := time.Now()
	for conn, client := range clients {
		if now.Sub(client.LastActivity) > 10*time.Second {
			log.Println("Closing client")
			client.Conn.Close()
			delete(clients, conn)
		}
	}
	clientsMutex.Unlock()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

var s http.Handler = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

func staticHandler(w http.ResponseWriter, r *http.Request) {
	s.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/static/", staticHandler)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", handleConnections)

	// go func() {
	// 	for {
	// 		closeStaleConnections()
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
