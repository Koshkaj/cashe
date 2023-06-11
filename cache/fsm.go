package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Koshkaj/cashe/core"
	"github.com/hashicorp/raft"
)

type CacheFSM struct {
	*Cache
}

func (f *CacheFSM) applySet(key []byte, value []byte, ttl time.Duration) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data[string(key)] = value
	// if ttl > 0 {
	// 	go func() {
	// 		// Move to the one goroutine with locking (channels)
	// 		<-time.After(ttl)
	// 		delete(f.data, string(key))
	// 	}()
	// }
	return nil
}

func (f *CacheFSM) applyDelete(key []byte) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.data, string(key))
	return nil
}

func NewCacheFSM(c Cacher) *CacheFSM {
	return &CacheFSM{
		Cache: c.(*Cache),
	}
}

type fsmSnapshot struct {
	data map[string][]byte
}

func (f *CacheFSM) Apply(l *raft.Log) interface{} {
	fmt.Printf("Address of data in fsm: %p %v\n", f.data, f.data)
	reader := bytes.NewReader(l.Data)
	cmd, err := core.ParseCommand(reader)
	fmt.Println(reader)
	if err != nil {
		return err
	}
	switch v := cmd.(type) {
	case *core.CommandSet:
		f.applySet(v.Key, v.Value, time.Duration(v.TTL))
	case *core.CommandDel:
		f.applyDelete(v.Key)
	default:
		return nil
	}
	fmt.Println("after set", f.data)
	return nil
}

func (f *CacheFSM) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	m := make(map[string][]byte)
	for k, v := range f.data {
		m[k] = v
	}
	return &fsmSnapshot{data: m}, nil
}

func (f *CacheFSM) Restore(serialized io.ReadCloser) error {
	m := make(map[string][]byte)
	if err := json.NewDecoder(serialized).Decode(&m); err != nil {
		return err
	}
	f.data = m
	return nil
}

func (fs *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		b, err := json.Marshal(fs.data)
		if err != nil {
			return err
		}
		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}
	return err
}

func (fs *fsmSnapshot) Release() {}
