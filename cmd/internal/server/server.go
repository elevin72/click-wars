package server

import (
	"context"
	"fmt"
	"log"
	// "time"
)

type Server struct {
	ctx        context.Context
	broadcast  chan IncomingClick
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
}

func NewServer() *Server {
	fmt.Println("new server")
	ctx := context.Background()

	return &Server{
		broadcast:  make(chan IncomingClick, 100000),
		register:   make(chan *Client, 10000),
		unregister: make(chan *Client, 10000),
		clients:    make(map[*Client]bool),
		ctx:        ctx,
	}
}

func (s *Server) Run() {
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
		case click := <-s.broadcast:
			var inc int32

			if click.side == LEFT {
				inc = 1
			} else {
				inc = -1
			}

			LinePosition.Add(inc)
			TotalHits.Add(1)

			serverClick := ServerClick{
				click,
				LinePosition.Load(),
				TotalHits.Load(),
			}

			log.Printf("Recieved %v", serverClick)
			fmt.Printf("num of connected clients %d\n", len(s.clients))
			for client := range s.clients {
				select {
				case client.send <- serverClick.outgoingBytes(client):
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}
