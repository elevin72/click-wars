package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn         *websocket.Conn
	LastActivity time.Time
	send         chan []byte
}
