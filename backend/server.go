package main

import (
	// "fmt"
	"context"
)

type Server struct {
	ctx        context.Context
	broadcast  chan *Click
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
	state      State
}

func newServer() *Server {
	ctx := context.Background()

	return &Server{
		broadcast:  make(chan *Click, 100000),
		register:   make(chan *Client, 10000),
		unregister: make(chan *Client, 10000),
		ctx:        ctx,
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send) // TODO
			}
		case message := <-s.broadcast:
			outBoundMessage := message.Serialize(&s.state)
			for client := range s.clients {
				select {
				case client.send <- outBoundMessage:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}
