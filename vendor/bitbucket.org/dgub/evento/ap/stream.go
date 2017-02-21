package ap

import (
	"sync"
)

type StreamSub struct {
	lock      sync.Mutex
	StreamMap map[string]uint64
}
