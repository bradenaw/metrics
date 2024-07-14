package metrics

import (
	"maps"
	"sync"
	"sync/atomic"
)

// metricMap is a subset of sync.Map, but uses generics/concrete types for the keys and values, thus
// avoiding wrapping either into any. It also does not allow deletes or overwrites, because the map
// only grows and metrics are only added once, which eliminates one more indirection. As such, it's
// slightly faster.
//
// It works the same way.
type metricMap[V any] struct {
	// read is the immutable part of the map. Lookups can check here with just an atomic load, and
	// on hit, can safely use the value they see.
	read atomic.Pointer[metricMapRO[V]]

	mu sync.Mutex
	// The number of loads since the last `promoteLocked()` with `amended` set that had to take `mu`
	// to examine `dirty`. When this grows too large, promote `dirty` into `read`.
	misses int
	// `dirty` is either nil (when `!read.amended`) or a superset of the map inside of `read`.
	dirty map[metricKey]V
}

type metricMapRO[V any] struct {
	m map[metricKey]V
	// True if `m` is incomplete, that is, that there are keys in `dirty` that are not in `m`.
	amended bool
}

func (m *metricMap[V]) Load(k metricKey) (V, bool) {
	ro := m.loadReadOnly()
	existing, ok := ro.m[k]
	if ok {
		return existing, true
	}
	if !ro.amended {
		// We saw a complete view of the map, so the key doesn't exist.
		var zero V
		return zero, false
	}

	// Slow path, have to check `m.dirty`.
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.dirty == nil {
		// Somebody already promoted.
		existing, ok := m.loadReadOnly().m[k]
		return existing, ok
	}

	existing, ok = m.dirty[k]
	m.misses++
	if m.misses > len(m.dirty) {
		m.promoteLocked()
	}
	return existing, ok
}

func (m *metricMap[V]) LoadOrStore(k metricKey, v V) (V, bool) {
	ro := m.loadReadOnly()
	existing, ok := ro.m[k]
	if ok {
		return existing, true
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.misses++
	if m.dirty == nil {
		// Have to load again in case somebody already showed up and promoted `dirty`.
		ro = m.loadReadOnly()
		// m.dirty==nil means that ro.m was complete up to this point, so we're the first write
		// since the last promotion.
		if ro.m == nil {
			m.dirty = make(map[metricKey]V)
		} else {
			m.dirty = maps.Clone(ro.m)
		}
		// Also mark ro as incomplete so that reads know they might need to check m.dirty.
		m.read.Store(&metricMapRO[V]{
			m:       ro.m,
			amended: true,
		})
	}
	existing, ok = m.dirty[k]
	if ok {
		return existing, true
	}
	m.dirty[k] = v
	return v, false
}

func (m *metricMap[V]) Range(f func(metricKey, V) bool) {
	ro := m.loadReadOnly()
	read := ro.m
	if ro.amended {
		// Same claim as sync.Map, Range is linear anyway so just promote everything immediately and
		// then range over what's left since Range is only guaranteed to see keys that were already
		// present at the point of the call to Range, not any that were added while ranging.
		m.mu.Lock()
		// Re-check to see if somebody else already promoted.
		if m.dirty == nil {
			ro = m.loadReadOnly()
			read = ro.m
		} else {
			read = m.dirty
			m.promoteLocked()
		}
		m.mu.Unlock()
	}
	for k, v := range read {
		if !f(k, v) {
			break
		}
	}
}

// promotes m.dirty into m.read, making m.read complete again.
func (m *metricMap[V]) promoteLocked() {
	m.read.Store(&metricMapRO[V]{
		m:       m.dirty,
		amended: false,
	})
	m.dirty = nil
	m.misses = 0
}

// Loads ro, acccounting for a zero-valued m.
func (m *metricMap[V]) loadReadOnly() metricMapRO[V] {
	ro := m.read.Load()
	if ro == nil {
		// nil map shows as empty, and if it's never been set then amended is false anyway.
		return metricMapRO[V]{}
	}
	return *ro
}
