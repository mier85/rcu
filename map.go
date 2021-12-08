package rcu

import (
	"sync"
	"sync/atomic"
)

// Map is a generic thread-safe read-copy-update map
type Map[KEY comparable, VALUE interface{}] struct {
	mutex    sync.Mutex
	internal map[KEY]VALUE
	readMap  *atomic.Value
}

// New creates a new generic RCU-map
func New[KEY comparable, VALUE interface{}]() Map[KEY, VALUE] {
	val := &atomic.Value{}
	val.Store(make(map[KEY]VALUE))
	return Map[KEY, VALUE]{
		internal: make(map[KEY]VALUE),
		readMap:  val,
	}
}

func (m *Map[KEY, VALUE]) updateReadMap() {
	readModel := make(map[KEY]VALUE, len(m.internal))
	for k, v := range m.internal {
		readModel[k] = v
	}
	m.readMap.Store(readModel)
}

// Insert inserts a value into the map
func (m *Map[KEY, VALUE]) Insert(key KEY, value VALUE) {
	m.mutex.Lock()
	m.internal[key] = value
	m.updateReadMap()
	m.mutex.Unlock()
}

// Delete deletes a value from the map
func (m *Map[KEY, VALUE]) Delete(key KEY) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, has := m.internal[key]
	if has {
		delete(m.internal, key)
		m.updateReadMap()
	}
}

// ReadMap returns the current Read-Copy of the map
func (m *Map[KEY, VALUE]) ReadMap() map[KEY]VALUE {
	return m.readMap.Load().(map[KEY]VALUE)
}
