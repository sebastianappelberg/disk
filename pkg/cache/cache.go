package cache

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Cache[T any] struct {
	dir    string
	path   string
	buffer *sync.Map
}

func NewCache[T any](dir, name string) *Cache[T] {
	var storedCache cacheStruct[T]
	buffer := &sync.Map{}

	cache := &Cache[T]{dir: dir, path: filepath.Join(dir, name+"_cache"), buffer: buffer}

	file, err := os.Open(cache.path)
	if err != nil {
		return cache
	}
	defer file.Close()

	err = gob.NewDecoder(file).Decode(&storedCache)
	if err != nil {
		return cache
	}
	for k, v := range storedCache.Value {
		buffer.Store(k, v)
	}
	return cache
}

// Internal structure for cached values
type cacheStruct[T any] struct {
	Created time.Time
	Value   map[string]T
}

// Put stores a value in the cache. Nothing is written to disk though.
func (c *Cache[T]) Put(key string, value T) {
	c.buffer.Store(key, value)
}

// Flush writes the cache to disk.
func (c *Cache[T]) Flush() {
	err := os.MkdirAll(c.dir, os.ModePerm)
	if err != nil {
		return
	}

	flushBuffer := make(map[string]T)
	c.buffer.Range(func(k, v interface{}) bool {
		flushBuffer[k.(string)] = v.(T)
		return true
	})

	cacheVal := cacheStruct[T]{Created: time.Now(), Value: flushBuffer}
	file, err := os.Create(c.path)
	if err != nil {
		return
	}
	defer file.Close()

	gob.NewEncoder(file).Encode(cacheVal)
}

// Get retrieves a value from the cache. If it didn't find the value in the buffer the second value is false.
func (c *Cache[T]) Get(key string) (T, bool) {
	val, ok := c.buffer.Load(key)
	if !ok {
		var zero T
		return zero, false
	}
	return val.(T), ok
}
