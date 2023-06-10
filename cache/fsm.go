package cache

import (
	"io"

	"github.com/hashicorp/raft"
)

type fsm struct {
	cache *Cache
}

func (f *fsm) Apply(l *raft.Log) interface{} {
	data := l.Data
	_ = data

	// Implement the Apply method to process the Raft log
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	// Implement the Snapshot method to create a snapshot of the FSM state
	return nil, nil
}

func (f *fsm) Restore(serialized io.ReadCloser) error {
	// Implement the Restore method to restore the FSM state from a snapshot
	return nil
}
