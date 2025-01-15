package main

import (
	"sync"
	"time"
)

type chacheEntry struct {
	createtAt time.Time
	val       []byte
}

type cache struct {
	memory   map[string]chacheEntry
	mux      *sync.Mutex
	interval time.Duration
	handler  time.Ticker
}

func (ch *cache) Add(key string, val []byte) {
	defer ch.mux.Unlock()
	ch.mux.Lock()
	if len(key) <= 0 || val == nil {
		return
	}
	ch.memory[key] = chacheEntry{createtAt: time.Now(), val: val}
}

func (ch *cache) Get(key string) ([]byte, bool) {
	defer ch.mux.Unlock()

	ch.mux.Lock()

	entry, state := ch.memory[key]
	if state {
		return entry.val, state
	}
	return nil, state
}

func (ch *cache) removeEntry(key string) {
	defer ch.mux.Unlock()

	ch.mux.Lock()
	delete(ch.memory, key)
}

func (ch *cache) readLoop() {
	go func() {
		for {
			tick := <-ch.handler.C

			for key, value := range ch.memory {
				passedItme := value.createtAt.Add(ch.interval)
				if passedItme.Before(tick) {
					ch.removeEntry(key)
				}
			}
		}
	}()
}

func (ch *cache) Stop() {
	ch.handler.Stop()
}

func NewCache(interval time.Duration) cache {
	c := cache{memory: make(map[string]chacheEntry), mux: &sync.Mutex{}, handler: *time.NewTicker(interval), interval: interval}
	c.readLoop()
	return c
}
