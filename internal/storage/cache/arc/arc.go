package arc

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"homework/internal/util"
	"sync"
)

type ARC[K comparable, V any] struct {
	p     int        // p Target (dynamic) size of cache
	c     int        // c Defines the maximum number of entries stored.
	l1    *list.List // l1 List of pages that have been seen only 1 recently (capturing recency).
	b1    *list.List // b1 List of evicted from l1 that seen only 1 recently
	l2    *list.List // l2 List of pages that have been seen at least 2 recently (capturing frequency).
	b2    *list.List // b2 List of evicted from l2 that seen at least 2 recently
	mutex sync.RWMutex
	len   int
	cache map[K]*entry[K]
}

// NewArcCache returns a new Adaptive Replacement Cache (ARC).
// c defines the maximum number of entries stored.
func NewArcCache[K comparable, V any](c int) *ARC[K, V] {
	return &ARC[K, V]{
		p:     0,
		c:     c,
		l1:    list.New(),
		b1:    list.New(),
		l2:    list.New(),
		b2:    list.New(),
		len:   0,
		cache: make(map[K]*entry[K], c),
	}
}

// Get retrieves a previous via Set inserted entry.
func (a *ARC[K, V]) Get(key K) (value V, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok {
		a.reconstruct(ent)
		val, err := a.bytesToV(ent.value)
		if err != nil {
			return *new(V), false
		}
		return val, !ent.ghost
	}
	return *new(V), false
}

// Put inserts a new key-value pair into the cache.
func (a *ARC[K, V]) Put(key K, value V) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	byteValue, err := a.VToBytes(value)
	if err != nil {
		return err
	}

	ent, ok := a.cache[key]
	if !ok {
		a.len++
		ent = &entry[K]{
			key:   key,
			value: byteValue,
			ghost: false,
		}
		a.reconstruct(ent)
		a.cache[key] = ent
	} else {
		if ent.ghost {
			a.len++
		}
		ent.value = byteValue
		ent.ghost = false
		a.reconstruct(ent)
	}
	return nil
}

// Delete removes an entry
func (a *ARC[K, V]) Delete(key K) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, ok := a.cache[key]
	if !ok {
		return util.ErrCacheDelete
	}

	delete(a.cache, key)
	a.len--

	return nil
}

// Len determines the number of currently cached entries.
func (a *ARC[K, V]) Len() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.len
}

// reconstruct Dynamically adjusts sizes of l1, l2 based on the workload.
// A hit in b1 increases p, shifting the balance towards recency (increasing the size of l1).
// A hit in b2 decreases p, shifting the balance towards frequency (increasing the size of l2).
func (a *ARC[K, V]) reconstruct(ent *entry[K]) {
	// comparing the pointer
	if ent.ll == a.l1 || ent.ll == a.l2 {
		// Case I	- Hit in l1 or l2
		ent.setMRU(a.l2)
	} else if ent.ll == a.b1 {
		// Case II	- Miss in l1 and l2 (hit in b1)
		var d int
		if a.b1.Len() >= a.b2.Len() {
			d = 1
		} else {
			d = a.b2.Len() / a.b1.Len()
		}
		a.p = min(a.p+d, a.c)

		a.replace(ent)
		ent.setMRU(a.l2)
	} else if ent.ll == a.b2 {
		// Case III	- Miss in l1 and l2 (hit in b2)
		var d int
		if a.b1.Len() <= a.b2.Len() {
			d = 1
		} else {
			d = a.b1.Len() / a.b2.Len()
		}
		a.p = max(a.p-d, 0)

		a.replace(ent)
		ent.setMRU(a.l2)
	} else if ent.ll == nil {
		// Case IV	- new entry
		if a.l1.Len()+a.b1.Len() == a.c {
			// if reached the cache capacity
			if a.l1.Len() < a.c {
				a.delLRU(a.b1)
				a.replace(ent)
			} else {
				a.delLRU(a.l1)
			}
		} else if a.l1.Len()+a.b1.Len() < a.c {
			// if not reached the cache capacity
			if a.l1.Len()+a.l2.Len()+a.b1.Len()+a.b2.Len() >= a.c {
				// if total number of entries of all lists >= cache capacity
				if a.l1.Len()+a.l2.Len()+a.b1.Len()+a.b2.Len() == 2*a.c {
					a.delLRU(a.b2)
				}
				a.replace(ent)
			}
		}
		// entry added to the front of l1, marking it as MRU
		ent.setMRU(a.l1)
	}
}

// delLRU remove the Least Recently Used (LRU) entry
func (a *ARC[K, V]) delLRU(list *list.List) {
	lru := list.Back()
	list.Remove(lru)
	a.len--
	delete(a.cache, lru.Value.(*entry[K]).key)
}

// replace provides balance between recency and frequency when:
// 1. new entry needs to be added.
// 2. cache has reached its capacity
func (a *ARC[K, V]) replace(ent *entry[K]) {
	if a.l1.Len() > 0 && ((a.l1.Len() > a.p) || (ent.ll == a.b2 && a.l1.Len() == a.p)) {
		lru := a.l1.Back().Value.(*entry[K])
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b1)
	} else {
		lru := a.l2.Back().Value.(*entry[K])
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b2)
	}
}

func (a *ARC[K, V]) VToBytes(val V) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (a *ARC[K, V]) bytesToV(data []byte) (value V, err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&value)
	if err != nil {
		return *new(V), err
	}
	return value, nil
}