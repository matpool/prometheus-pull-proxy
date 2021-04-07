package common

import (
	"sync"
	"testing"
)

func TestCircularQueue(t *testing.T) {
	mtx := &sync.Mutex{}
	q := &CircularQueue{
		mtx:    mtx,
		cond:   sync.NewCond(mtx),
		buffer: make([]interface{}, 6),
		tail:   0,
		head:   0,
	}

	for i := 0; i < 5; i++ {
		q.Put(i)
	}
	evicted := q.Put(5)
	if evicted.(int) != 0 {
		t.Fatal("unexpected evicted item")
	}
	for i := 0; i < 5; i++ {
		if (q.Get()).(int) != i+1 {
			t.Fatal("unexpected get item")
		}
	}
}
