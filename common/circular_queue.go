package common

import (
	"sync"
)

type CircularQueue struct {
	mtx    *sync.Mutex
	cond   *sync.Cond
	buffer []interface{}
	head   int
	tail   int
}

func (c *CircularQueue) idxInc(i int) int {
	return (i + 1) % len(c.buffer)
}

// Put put item into queue. if queue full, evict head item, put new item, and return evicted item
func (c *CircularQueue) Put(i interface{}) interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var ret interface{}

	// queue full, evict head item
	if c.idxInc(c.tail) == c.head {
		ret = c.buffer[c.head]
		c.head = c.idxInc(c.head)
	}
	c.buffer[c.tail] = i
	c.tail = c.idxInc(c.tail)

	c.cond.Signal()

	return ret
}

// Get get item from queue. if queue empty, caller will sleep until queue has items
func (c *CircularQueue) Get() interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for c.head == c.tail {
		c.cond.Wait()
	}

	ret := c.buffer[c.head]
	c.head = c.idxInc(c.head)

	return ret
}
