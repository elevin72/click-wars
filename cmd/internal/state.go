package internal

import (
	"sync/atomic"
)

var LinePosition atomic.Int32

var TotalHits atomic.Int32

var LeftHits atomic.Int32

var RightHits atomic.Int32
