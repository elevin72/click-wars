package click

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
)

type Side = byte

const (
	LEFT  Side = 0
	RIGHT Side = 1
)

const OutgoingMessageLength = 18

type IncomingClick struct {
	X          float32   `db:"x"`
	Y          float32   `db:"y"`
	Side       Side      `db:"side"`
	SenderUUID uuid.UUID `db:"uuid"`
}

func (c IncomingClick) String() string {
	return fmt.Sprintf("Click {x: %f, y: %f, side: %d, senderuuid: %v}", c.X, c.Y, c.Side, c.SenderUUID)
}

func ParseIncomingClick(message []byte, senderUUID uuid.UUID) (IncomingClick, error) {
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
		X:          x,
		Y:          y,
		Side:       side,
		SenderUUID: senderUUID,
	}, nil
}

type ServerClick struct {
	X            float32   `db:"x"`
	Y            float32   `db:"y"`
	Side         Side      `db:"side"`
	SenderUUID   uuid.UUID `db:"sender_uuid"`
	LinePosition int32     `db:"line_position"`
	TotalHits    int32     `db:"total_hits"`
	Time         string    `db:"time"`
	bytes        []byte
}

func NewServerClick(ic IncomingClick, linePosition int32, totalHits int32) ServerClick {
	sc := ServerClick{
		X:            ic.X,
		Y:            ic.Y,
		Side:         ic.Side,
		SenderUUID:   ic.SenderUUID,
		LinePosition: linePosition,
		TotalHits:    totalHits,
		Time:         time.Now().UTC().Format(time.RFC3339Nano),
	}
	sc.setAsBytes()
	return sc
}

func (c ServerClick) String() string {
	return fmt.Sprintf("Click {x: %f, y: %f, side: %d, senderuuid: %v, linePosition: %v, totalHits: %v, time: %v}", c.X, c.Y, c.Side, c.SenderUUID, c.LinePosition, c.TotalHits, c.Time)
}

func (c *ServerClick) OutgoingBytes(targetUUID uuid.UUID) []byte {
	if c.SenderUUID == targetUUID {
		log.Println("creating bytes to self")
		bytesToSelf := make([]byte, OutgoingMessageLength)
		log.Printf("len(c.bytes) %d, len(bytesToSelf) %d\n", len(c.bytes), len(bytesToSelf))
		copy(bytesToSelf, c.bytes)
		log.Printf("bytes to self %v\n", bytesToSelf)
		bytesToSelf[0] = uint8(1)
		return bytesToSelf
	}
	log.Printf("bytes to other %v\n", c.bytes)
	return c.bytes
}

func (c *ServerClick) setAsBytes() {
	log.Println("creating bytes to other")
	msg := make([]byte, OutgoingMessageLength)
	msg[0] = uint8(0)
	binary.LittleEndian.PutUint32(msg[1:5], uint32(c.LinePosition))
	binary.LittleEndian.PutUint32(msg[5:9], uint32(c.TotalHits))
	binary.LittleEndian.PutUint32(msg[9:13], math.Float32bits(c.X))
	binary.LittleEndian.PutUint32(msg[13:17], math.Float32bits(c.Y))
	msg[17] = c.Side
	c.bytes = msg
}
