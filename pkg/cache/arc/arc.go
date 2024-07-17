package arc

import (
	"container/list"
	"sync"
)

type ARC struct {
	p     int        // p Target (dynamic) size of cache
	c     int        // c Defines the maximum number of entries stored.
	l1    *list.List // l1 List of pages that have been seen only 1 recently (capturing recency).
	b1    *list.List // b1 List of evicted from l1 that seen only 1 recently
	l2    *list.List // l2 List of pages that have been seen at least 2 recently (capturing frequency).
	b2    *list.List // b2 List of evicted from l2 that seen at least 2 recently
	mutex sync.RWMutex
	len   int
	cache map[any]*entry
}

// New returns a new Adaptive Replacement Cache (ARC).
func New(c int) *ARC {
	return &ARC{
		p:     0,
		c:     c,
		l1:    list.New(),
		b1:    list.New(),
		l2:    list.New(),
		b2:    list.New(),
		len:   0,
		cache: make(map[interface{}]*entry, c),
	}
}

// Put inserts a new key-value pair into the cache.
func (a *ARC) Put(key, value interface{}) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok != true {
		a.len++

		ent = &entry{
			key:   key,
			value: value,
			ghost: false,
		}

		a.req(ent)
		a.cache[key] = ent
	} else {
		if ent.ghost {
			a.len++
		}
		ent.value = value
		ent.ghost = false
		a.req(ent)
	}
	return ok
}

// Get retrieves a previous via Set inserted entry.
func (a *ARC) Get(key interface{}) (value interface{}, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok {
		a.req(ent)
		return ent.value, !ent.ghost
	}
	return nil, false
}

// Len determines the number of currently cached entries.
func (a *ARC) Len() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.len
}

// req Dynamically adjusts sizes of l1, l2 based on the workload.
// A hit in b1 increases p, shifting the balance towards recency (increasing the size of l1).
// A hit in b2 decreases p, shifting the balance towards frequency (increasing the size of l2).
func (a *ARC) req(ent *entry) {
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
func (a *ARC) delLRU(list *list.List) {
	lru := list.Back()
	list.Remove(lru)
	a.len--
	delete(a.cache, lru.Value.(*entry).key)
}

// replace provides balance between recency and frequency when:
// 1. new entry needs to be added.
// 2. cache has reached its capacity
func (a *ARC) replace(ent *entry) {
	if a.l1.Len() > 0 && ((a.l1.Len() > a.p) || (ent.ll == a.b2 && a.l1.Len() == a.p)) {
		lru := a.l1.Back().Value.(*entry)
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b1)
	} else {
		lru := a.l2.Back().Value.(*entry)
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b2)
	}
}