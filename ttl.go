package cache

import (
	"time"
)

type TTL struct {
	db    map[string]entry
	queue chan func()
	quit  chan bool
}

// NewTTL returns a new TTL cache.
func NewTTL() *TTL {
	self := &TTL{
		db:    make(map[string]entry),
		queue: make(chan func()),
		quit:  make(chan bool),
	}
	go self.run_queue()
	return self
}

// Quit tells the TTL it is no longer needed.
func (self *TTL) Quit() {
	self.quit <- true
}

// Set sets the value associated with the given key, and sets the TTL
// for that value.
func (self *TTL) Set(key string, value interface{}, ttl time.Duration) {
	self.queue <- func() {
		// time.Duration represents number of nanoseconds
		self.db[key] = entry{value, time.Now().UnixNano() + int64(ttl)}
	}
}

// Get returns the value associated with the given key.  When the
// value's expiry has passed, it prunes the value from the map and
// returns nil,false.
func (self *TTL) Get(key string) (interface{}, bool) {
	rq := make(chan result)
	self.queue <- func() {
		entry, ok := self.db[key]
		if ok {
			now := time.Now().UnixNano()
			if entry.expiry > now {
				rq <- result{entry.value, ok}
				return
			}
			delete(self.db, key)
		}
		rq <- result{nil, false}
	}
	res := <-rq
	return res.value, res.ok
}

// Prune removes all values from cache that have expired.
func (self *TTL) Prune() {
	self.queue <- func() {
		now := time.Now().UnixNano()
		for k, v := range self.db {
			if v.expiry > now {
				delete(self.db, k)
			}
		}
	}
}

const (
	nanosPerMilli = time.Millisecond / time.Nanosecond
)

type entry struct {
	value  interface{}
	expiry int64
}

type result struct {
	value interface{}
	ok    bool
}

func (self *TTL) run_queue() {
	for {
		select {
		case fn := <-self.queue:
			fn()
		case <-self.quit:
			break
		}
	}
}