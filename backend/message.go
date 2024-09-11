package main

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/gorilla/websocket"
)

type Color = uint8

const (
	BLUE Color = 0
	RED  Color = 1
)

type Client struct {
	Conn         *websocket.Conn
	LastActivity time.Time
	send         chan []byte
}

func Deserialize(messageBytes []byte) *Click {
	x := math.Float32frombits(binary.LittleEndian.Uint32(messageBytes[1:5]))
	y := math.Float32frombits(binary.LittleEndian.Uint32(messageBytes[5:9]))
	color := messageBytes[9] // 0 for blue, 1 for red
	return &Click{
		x:     x,
		y:     y,
		color: color,
	}
}

type Click struct {
	x     float32
	y     float32
	color Color
}

func (click *Click) Serialize(state *State) []byte {

	state.RLock()
	linePosition := state.linePosition
	totalHits := state.totalHits
	state.RUnlock()

	msg := make([]byte, 18)
	msg[0] = LONG
	binary.LittleEndian.PutUint32(msg[1:5], uint32(linePosition))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(totalHits))
	binary.LittleEndian.PutUint32(msg[9:13], math.Float32bits(click.x))
	binary.LittleEndian.PutUint32(msg[13:17], math.Float32bits(click.y))
	msg[17] = click.color

	return msg
}
