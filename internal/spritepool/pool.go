package spritepool

import (
	"fmt"
	"sync"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// SpriteInstance represents a runtime instance of a sprite
type SpriteInstance struct {
	Name       string                 // Sprite name
	X          float64                // X position
	Y          float64                // Y position
	Age        time.Duration          // Age since spawn (for lifetime management)
	IsActive   bool                   // Whether this instance is currently active
	IsPooled   bool                   // Whether this instance came from a pool
	Data       cartio.SpriteData      // Reference to sprite data
	CustomData map[string]interface{} // Custom data storage for Lua
}

// Pool manages a collection of sprite instances for reuse
type Pool struct {
	spriteName string
	spriteData cartio.SpriteData
	maxSize    int                      // Maximum pool size (from MaxSpawn)
	available  []*SpriteInstance        // Available instances ready to be used
	active     map[*SpriteInstance]bool // Active instances currently in use
	mu         sync.RWMutex             // Mutex for thread safety
}

// NewPool creates a new pool for a sprite type
func NewPool(spriteName string, spriteData cartio.SpriteData, initialSize int, maxSize int) *Pool {
	pool := &Pool{
		spriteName: spriteName,
		spriteData: spriteData,
		maxSize:    maxSize,
		available:  make([]*SpriteInstance, 0, initialSize),
		active:     make(map[*SpriteInstance]bool),
	}

	// Pre-allocate initial instances
	for i := 0; i < initialSize; i++ {
		instance := &SpriteInstance{
			Name:       spriteName,
			X:          0,
			Y:          0,
			Age:        0,
			IsActive:   false,
			IsPooled:   true,
			Data:       spriteData,
			CustomData: make(map[string]interface{}),
		}
		pool.available = append(pool.available, instance)
	}

	return pool
}

// Acquire gets an instance from the pool, creating a new one if needed
func (p *Pool) Acquire() (*SpriteInstance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if we have available instances
	if len(p.available) > 0 {
		instance := p.available[len(p.available)-1]
		p.available = p.available[:len(p.available)-1]
		p.active[instance] = true
		instance.IsActive = true
		instance.Age = 0
		// Reset position and custom data
		instance.X = 0
		instance.Y = 0
		for k := range instance.CustomData {
			delete(instance.CustomData, k)
		}
		return instance, nil
	}

	// No available instances - check if we can grow
	activeCount := len(p.active)
	if activeCount >= p.maxSize {
		// Pool exhausted - create overflow instance (non-pooled)
		return &SpriteInstance{
			Name:       p.spriteName,
			X:          0,
			Y:          0,
			Age:        0,
			IsActive:   true,
			IsPooled:   false, // Overflow instance is not pooled
			Data:       p.spriteData,
			CustomData: make(map[string]interface{}),
		}, fmt.Errorf("pool exhausted for %s, creating overflow instance", p.spriteName)
	}

	// Grow pool dynamically
	growthSize := p.calculateGrowthSize(activeCount)
	for i := 0; i < growthSize; i++ {
		instance := &SpriteInstance{
			Name:       p.spriteName,
			X:          0,
			Y:          0,
			Age:        0,
			IsActive:   false,
			IsPooled:   true,
			Data:       p.spriteData,
			CustomData: make(map[string]interface{}),
		}
		p.available = append(p.available, instance)
	}

	// Get the first newly created instance
	instance := p.available[len(p.available)-1]
	p.available = p.available[:len(p.available)-1]
	p.active[instance] = true
	instance.IsActive = true
	return instance, nil
}

// Release returns an instance to the pool
func (p *Pool) Release(instance *SpriteInstance) error {
	if !instance.IsPooled {
		return fmt.Errorf("cannot release non-pooled instance")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.active[instance] {
		return fmt.Errorf("instance not active in pool")
	}

	// Reset instance state
	instance.IsActive = false
	instance.Age = 0
	instance.X = 0
	instance.Y = 0
	for k := range instance.CustomData {
		delete(instance.CustomData, k)
	}

	// Remove from active and add to available
	delete(p.active, instance)
	p.available = append(p.available, instance)
	return nil
}

// Update updates all active instances and handles lifetime expiration
func (p *Pool) Update(deltaTime time.Duration) []*SpriteInstance {
	p.mu.RLock()
	defer p.mu.RUnlock()

	expired := make([]*SpriteInstance, 0)

	// If sprite has lifetime, check for expiration
	if p.spriteData.Lifetime > 0 {
		lifetime := time.Duration(p.spriteData.Lifetime) * time.Millisecond
		for instance := range p.active {
			if instance.IsActive {
				instance.Age += deltaTime
				if instance.Age >= lifetime {
					expired = append(expired, instance)
				}
			}
		}
	}

	return expired
}

// GetStats returns pool statistics
func (p *Pool) GetStats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	activeCount := len(p.active)
	return PoolStats{
		SpriteName:  p.spriteName,
		Available:   len(p.available),
		Active:      activeCount,
		MaxSize:     p.maxSize,
		Utilization: float64(activeCount) / float64(p.maxSize),
	}
}

// calculateGrowthSize calculates how much to grow the pool
func (p *Pool) calculateGrowthSize(activeCount int) int {
	// Grow by 25% of max size, or enough to reach max, whichever is smaller
	growthSize := int(float64(p.maxSize) * 0.25)
	remaining := p.maxSize - activeCount
	if growthSize > remaining {
		growthSize = remaining
	}
	if growthSize < 1 {
		growthSize = 1
	}
	return growthSize
}

// PoolStats contains statistics about a pool
type PoolStats struct {
	SpriteName  string
	Available   int
	Active      int
	MaxSize     int
	Utilization float64
}
