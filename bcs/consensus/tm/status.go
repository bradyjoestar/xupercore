package tm

import "sync"

type TMStatus struct {
	startHeight int64
	mutex       sync.RWMutex
	newHeight   int64
	index       int
	config      *TMConfig
}
