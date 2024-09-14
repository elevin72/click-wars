package main

import (
	"context"
	"fmt"
	"log"
	// "time"
)

type Server struct {
	ctx        context.Context
	broadcast  chan *Click
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
	state      *State
}

func newServer() *Server {
	fmt.Println("new server")
	ctx := context.Background()

	return &Server{
		broadcast:  make(chan *Click, 100000),
		register:   make(chan *Client, 10000),
		unregister: make(chan *Client, 10000),
		clients:    make(map[*Client]bool),
		ctx:        ctx,
		state:      &InitState,
	}
}

func (s *Server) run() {
	for {
		select {
		// register new client
		case client := <-s.register:
			log.Println("Registering new client")
			s.clients[client] = true
		// unregister client
		case client := <-s.unregister:
			log.Println("Unregistering client")
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send) // TODO ??
			}
		// broadcast click
		case click := <-s.broadcast:
			var inc int

			if click.color == BLUE {
				inc = 1
			} else {
				inc = -1
			}

			s.state.Lock()
			s.state.linePosition += inc
			s.state.totalHits++
			s.state.Unlock()
			log.Printf("Recieved %v", click)

			s.state.RLock()
			log.Printf("line position: %d\n", s.state.linePosition)
			log.Printf("total hits: %d\n", s.state.totalHits)
			s.state.RUnlock()

			outBoundMessage := click.Serialize(s.state)
			fmt.Printf("num clients %d\n", len(s.clients))
			for client := range s.clients {
				fmt.Printf("client %v\n", client.ip)

			}
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
