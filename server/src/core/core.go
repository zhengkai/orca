package core

import "sync"

// Core ...
type Core struct {
	serial int
	pool   map[[16]byte]*row
	mux    sync.Mutex
}

// NewCore ...
func NewCore() *Core {
	return &Core{
		pool: make(map[[16]byte]*row),
	}
}
