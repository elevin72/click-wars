package server

import (
	"context"
	"fmt"
	"log"

	"github.com/elevin72/click-wars/internal/click"
	"github.com/elevin72/click-wars/internal/db"
)

type Server struct {
	ctx        context.Context
	broadcast  chan click.IncomingClick
	register   chan *Client
	unregister chan *Client
	db         *db.Database
	clients    map[*Client]bool
}

func New() *Server {
	ctx := context.Background()

	db := db.New()
	linePosition, totalHits := db.State()
	log.Printf("loading from db: linePosition: %d, totalHits: %d\n", linePosition, totalHits)
	LinePosition.Store(linePosition)
	TotalHits.Store(totalHits)

	return &Server{
		broadcast:  make(chan click.IncomingClick, 100000),
		register:   make(chan *Client, 10000),
		unregister: make(chan *Client, 10000),
		clients:    make(map[*Client]bool),
		db:         db,
		ctx:        ctx,
	}
}

func (s *Server) Run() {
	defer s.db.Close()
	for {
		select {
		case client := <-s.register:
			log.Printf("Registering new client, %v\n", client)
			s.clients[client] = true
		case client := <-s.unregister:
			log.Printf("Unregistering client, %v\n", client)
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
		case incomingClick := <-s.broadcast:
			var inc int32

			if incomingClick.Side == click.LEFT {
				inc = 1
			} else {
				inc = -1
			}

			linePosition := LinePosition.Add(inc)
			totalHits := TotalHits.Add(1)

			serverClick := click.NewServerClick(incomingClick, linePosition, totalHits)

			// store click
			s.db.InsertClick(&serverClick)

			log.Printf("Recieved %v", serverClick)
			fmt.Printf("num of connected clients %d\n", len(s.clients))
			for client := range s.clients {
				select {
				case client.send <- serverClick.OutgoingBytes(client.UUID):
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}
