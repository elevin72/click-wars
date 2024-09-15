package main

import (
	"sync/atomic"
)

var linePosition atomic.Int32

var totalHits atomic.Int32

var leftHits atomic.Int32

var rightHits atomic.Int32
