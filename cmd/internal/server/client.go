package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	server   *Server
	conn     *websocket.Conn
	lastPing time.Time //fairly certain that lastping is pointless now
	send     chan []byte
	UUID     uuid.UUID // TODO: use counter atomic.Int64 instead of uuid
}

func (c *Client) String() string {
	return fmt.Sprintf("Client{UUID: %v}", c.UUID)
}

const (
	writeWait  = 5 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// origin := r.Header.Get("Origin")
		// return origin == "http://127.0.0.1:8080" || origin == "http://localhost:8080"
	},
}

func ServeWs(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		server:   server,
		conn:     conn,
		lastPing: time.Now(),
		send:     make(chan []byte),
		UUID:     uuid.New(),
	}

	server.register <- client

	go client.readFromSocket()
	go client.writeToSocket()

}

func (c *Client) readFromSocket() {
	defer func() {
		fmt.Println("In defer of readFromSocket")
		c.server.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, incomingBytes, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Client closed")
			}
			break
		}

		if incomingBytes[0] == websocket.PingMessage {
			log.Println("Pong")
			c.lastPing = time.Now()
			continue
		}

		click, err := parseIncomingClick(incomingBytes, c.UUID)
		if err != nil {
			log.Println(err)
		}

		c.server.broadcast <- click
	}
}

func (c *Client) writeMessage(message []byte) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) writeToSocket() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.writeMessage(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				c.writeMessage(message)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Print(err)
				return
			}
		}
	}
}
