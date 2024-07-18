package arc

import (
	"container/list"
)

// entry represents an individual item stored in cache
type entry[K comparable] struct {
	key   K
	value []byte
	ll    *list.List    // ll - pointer to the list it belongs to (l1, l2, b1, or b2).
	el    *list.Element // el - used to directly manipulate the entry’s position in the list.
	ghost bool
}

// setLRU sets  entry as LRU in provided list
func (e *entry[K]) setLRU(list *list.List) {
	e.detach()
	e.ll = list
	e.el = e.ll.PushBack(e)
}

// setMRU sets entry as MRU in provided list
func (e *entry[K]) setMRU(list *list.List) {
	e.detach()
	e.ll = list
	e.el = e.ll.PushFront(e)
}

// detach remove entry from list
func (e *entry[K]) detach() {
	if e.ll != nil {
		e.ll.Remove(e.el)
	}
}