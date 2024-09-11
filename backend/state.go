package main

import (
	"sync"
)

type State struct {
	linePosition int
	totalHits    int
	leftHits     int
	rightHits    int
	sync.RWMutex
}

var InitState State = State{
	linePosition: 0,
	totalHits:    0,
	leftHits:     0,
	rightHits:    0,
}
