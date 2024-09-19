package server

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Color = byte

const (
	BLUE Color = 0
	RED  Color = 1
)

type Click struct {
	x     float32
	y     float32
	color Color
}

func (c *Click) String() string {
	return fmt.Sprintf("Click {x: %f, y: %f, color: %d}", c.x, c.y, c.color)
}

func deserializeClick(message []byte) (*Click, error) {
	if len(message) != 10 {
		return nil, fmt.Errorf("expected message of 10 bytes, got message [%v] with length of%d", message, len(message))
	}

	x := math.Float32frombits(binary.LittleEndian.Uint32(message[1:5]))
	y := math.Float32frombits(binary.LittleEndian.Uint32(message[5:9]))
	color := message[9]

	if color != BLUE && color != RED {
		return nil, fmt.Errorf("expected 9th byte to be 0 or 1, was %d", color)
	}

	return &Click{
		x:     x,
		y:     y,
		color: color,
	}, nil
}

func (click *Click) Serialize() []byte {

	// this is dumb
	msg := make([]byte, 18)
	msg[0] = 0
	binary.LittleEndian.PutUint32(msg[1:5], uint32(LinePosition.Load()))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(TotalHits.Load()))
	binary.LittleEndian.PutUint32(msg[9:13], math.Float32bits(click.x))
	binary.LittleEndian.PutUint32(msg[13:17], math.Float32bits(click.y))
	msg[17] = click.color

	return msg
}
