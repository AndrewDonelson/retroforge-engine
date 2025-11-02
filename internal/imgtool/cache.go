package imgtool

import (
	"sync"
	"sync/atomic"
)

// ColorCache caches RGB to palette index mappings for performance
type ColorCache struct {
	mu     sync.RWMutex
	cache  map[uint32]int // RGB packed → palette index
	hits   uint64
	misses uint64
}

// NewColorCache creates a new color cache
func NewColorCache() *ColorCache {
	return &ColorCache{
		cache: make(map[uint32]int, 10000),
	}
}

// Get retrieves a cached palette index for an RGB color
func (c *ColorCache) Get(r, g, b uint8) (int, bool) {
	key := packRGB(r, g, b)

	c.mu.RLock()
	index, ok := c.cache[key]
	c.mu.RUnlock()

	if ok {
		atomic.AddUint64(&c.hits, 1)
	} else {
		atomic.AddUint64(&c.misses, 1)
	}

	return index, ok
}

// Set stores a palette index for an RGB color
func (c *ColorCache) Set(r, g, b uint8, index int) {
	key := packRGB(r, g, b)

	c.mu.Lock()
	c.cache[key] = index
	c.mu.Unlock()
}

// Stats returns cache hit/miss statistics
func (c *ColorCache) Stats() (hits, misses uint64) {
	return atomic.LoadUint64(&c.hits), atomic.LoadUint64(&c.misses)
}

// Clear clears the cache
func (c *ColorCache) Clear() {
	c.mu.Lock()
	c.cache = make(map[uint32]int, 10000)
	c.mu.Unlock()
}

// packRGB packs RGB values into a single uint32
func packRGB(r, g, b uint8) uint32 {
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

