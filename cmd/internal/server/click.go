package server

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type Side = byte

const (
	LEFT  Side = 0
	RIGHT Side = 1
)

const OutgoingMessageLength = 18

type IncomingClick struct {
	x, y       float32
	side       Side
	senderUUID uuid.UUID
}

func (c IncomingClick) String() string {
	return fmt.Sprintf("Click {x: %f, y: %f, side: %d, senderuuid: %v}", c.x, c.y, c.side, c.senderUUID)
}

func parseIncomingClick(message []byte, senderUUID uuid.UUID) (IncomingClick, error) {
	if len(message) != 10 {
		return IncomingClick{}, fmt.Errorf("expected message of 10 bytes, got message [%v] with length of%d", message, len(message))
	}

	x := math.Float32frombits(binary.LittleEndian.Uint32(message[1:5]))
	y := math.Float32frombits(binary.LittleEndian.Uint32(message[5:9]))
	side := message[9]

	if side != LEFT && side != RIGHT {
		return IncomingClick{}, fmt.Errorf("expected 9th byte to be 0 or 1, was %d", side)
	}

	return IncomingClick{
		x:          x,
		y:          y,
		side:       side,
		senderUUID: senderUUID,
	}, nil
}

type ServerClick struct {
	IncomingClick
	linePosition int32
	totalHits    int32
}

func (c ServerClick) String() string {
	return fmt.Sprintf("Click {x: %f, y: %f, side: %d, senderuuid: %v, linePosition: %v, totalHits: %v}", c.x, c.y, c.side, c.senderUUID, c.linePosition, c.totalHits)
}

func (c *ServerClick) outgoingBytes(target *Client) []byte {
	msg := make([]byte, OutgoingMessageLength)
	if c.senderUUID == target.UUID {
		fmt.Println("sending to self")
		msg[0] = uint8(1)
	} else {
		fmt.Println("sending to other")
		msg[0] = uint8(0)
	}
	binary.LittleEndian.PutUint32(msg[1:5], uint32(c.linePosition))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(c.totalHits))
	binary.LittleEndian.PutUint32(msg[9:13], math.Float32bits(c.x))
	binary.LittleEndian.PutUint32(msg[13:17], math.Float32bits(c.y))
	msg[17] = c.side
	return msg
}
