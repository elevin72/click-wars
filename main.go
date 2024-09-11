package main

import (
	"encoding/binary"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Color = uint8

const (
	BLUE Color = 0
	RED  Color = 1
)

const (
	LONG  byte = 0
	SHORT byte = 1
)

type Click struct {
	x     float32
	y     float32
	color Color
}

func makeLongMessage(click Click) []byte {
	msg := make([]byte, 18)
	msg[0] = LONG
	binary.LittleEndian.PutUint32(msg[1:5], uint32(LINE_POSITION))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(TOTAL_HITS))
	binary.LittleEndian.PutUint32(msg[9:13], math.Float32bits(click.x))
	binary.LittleEndian.PutUint32(msg[13:17], math.Float32bits(click.y))
	msg[17] = click.color
	return msg
}

func makeShortMessage() []byte {
	msg := make([]byte, 9)
	msg[0] = SHORT
	binary.LittleEndian.PutUint32(msg[1:5], uint32(LINE_POSITION))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(TOTAL_HITS))
	return msg
}

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
	clients       = make(map[*websocket.Conn]*Client)
	clientsMutex  sync.Mutex
	LINE_POSITION int32 = 0
	lineMutex     sync.Mutex
	TOTAL_HITS    int32 = 0
	hitsMutex     sync.Mutex
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		ws.Close()
	}()

	client := &Client{
		Conn:         ws,
		LastActivity: time.Now(),
	}

	clientsMutex.Lock()
	log.Println("Adding client")
	clients[ws] = client
	clientsMutex.Unlock()

	initMessage := make([]byte, 9)
	initMessage[0] = 1
	binary.LittleEndian.PutUint32(initMessage[1:5], uint32(LINE_POSITION))
	binary.LittleEndian.PutUint32(initMessage[5:9], uint32(TOTAL_HITS))

	shortMessage := makeShortMessage()
	preparedShortMessage, err := websocket.NewPreparedMessage(websocket.BinaryMessage, shortMessage)
	if err != nil {
		log.Println(err)
	}
	err = client.Conn.WritePreparedMessage(preparedShortMessage)
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

		if len(message) != 10 {
			if message[0] == 0xFF {
				// ping pong
				log.Println("ping")
				client.LastActivity = time.Now()
			} else {
				log.Println("Invalid message length")
				log.Println(string(message))
			}
			continue
		}

		if len(message) != 10 {
			log.Println("Invalid message length")
			continue
		}

		client.LastActivity = time.Now()
		x := math.Float32frombits(binary.LittleEndian.Uint32(message[1:5]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(message[5:9]))
		color := message[9] // 0 for blue, 1 for red

		var lineIncrement int32
		if color == 0 {
			lineIncrement = 1
		} else {
			lineIncrement = -1
		}
		if LINE_POSITION > -80 && LINE_POSITION < 80 {
			lineMutex.Lock()
			LINE_POSITION += lineIncrement
			lineMutex.Unlock()
		}
		hitsMutex.Lock()
		TOTAL_HITS++
		hitsMutex.Unlock()

		log.Printf("Received x: %f, y: %f, color: %d\n", x, y, color)
		log.Printf("New line position: %d\n", LINE_POSITION)
		log.Printf("New total hits : %d\n", TOTAL_HITS)

		click := Click{
			x:     x,
			y:     y,
			color: color,
		}

		ShortMessage := makeShortMessage()
		longMessage := makeLongMessage(click)

		clientsMutex.Lock()
		for _, c := range clients {
			var err error
			if c != client {
				err = c.Conn.WriteMessage(websocket.BinaryMessage, longMessage)
			} else {
				err = c.Conn.WriteMessage(websocket.BinaryMessage, ShortMessage)
			}
			if err != nil {
				log.Println(err)
			}
		}
		clientsMutex.Unlock()

	}
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
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Print(err)
	}
	t.Execute(w, struct {
		LinePosition int32
		TotalHits    int32
	}{
		LinePosition: LINE_POSITION,
		TotalHits:    TOTAL_HITS,
	})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else {
		w.Header().Set("Content-Type", "text/css")
	}
	http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
}

func percentageOnLoadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("hit percentage.css\n")
	t, err := template.ParseFiles("static/css/percentage.css.tmpl")
	if err != nil {
		log.Print(err)
	}
	percentage := float32(LINE_POSITION/-2) + 50
	w.Header().Set("Content-Type", "text/css")
	t.Execute(w, struct {
		LeftSide  float32
		RightSide float32
	}{
		LeftSide:  100 - percentage,
		RightSide: percentage,
	})
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/static/", staticHandler)
	// http.HandleFunc("/static/css/percentage.css", percentageOnLoadHandler)
	http.HandleFunc("/ws", handleConnections)

	go func() {
		for {
			closeStaleConnections()
			time.Sleep(5 * time.Second)
		}
	}()

	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
