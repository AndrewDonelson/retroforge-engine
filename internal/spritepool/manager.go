package spritepool

import (
	"fmt"
	"sync"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// PoolManager manages all sprite pools
type PoolManager struct {
	pools map[string]*Pool // Map of sprite name to pool
	mu    sync.RWMutex     // Mutex for thread safety
}

// NewPoolManager creates a new pool manager
func NewPoolManager() *PoolManager {
	return &PoolManager{
		pools: make(map[string]*Pool),
	}
}

// ShouldPool determines if a sprite should be pooled based on its properties
// Criteria: isUI == false AND maxSpawn > 10
func ShouldPool(spriteData cartio.SpriteData) bool {
	return !spriteData.IsUI && spriteData.MaxSpawn > 10
}

// RegisterSprite checks if a sprite should be pooled and creates a pool if needed
func (pm *PoolManager) RegisterSprite(spriteName string, spriteData cartio.SpriteData) error {
	if !ShouldPool(spriteData) {
		return nil // Not eligible for pooling
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if pool already exists
	if _, exists := pm.pools[spriteName]; exists {
		return nil // Pool already exists
	}

	// Calculate initial pool size (50% of max, capped at 50)
	initialSize := int(float64(spriteData.MaxSpawn) * 0.5)
	if initialSize > 50 {
		initialSize = 50
	}
	if initialSize < 1 {
		initialSize = 1
	}

	// Create pool
	pool := NewPool(spriteName, spriteData, initialSize, spriteData.MaxSpawn)
	pm.pools[spriteName] = pool

	return nil
}

// HasPool checks if a pool exists for a sprite
func (pm *PoolManager) HasPool(spriteName string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, exists := pm.pools[spriteName]
	return exists
}

// Acquire gets an instance from the pool for a sprite, or returns nil if no pool exists
func (pm *PoolManager) Acquire(spriteName string) (*SpriteInstance, error) {
	pm.mu.RLock()
	pool, exists := pm.pools[spriteName]
	pm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no pool exists for sprite %s", spriteName)
	}

	return pool.Acquire()
}

// Release returns an instance to its pool
func (pm *PoolManager) Release(instance *SpriteInstance) error {
	if !instance.IsPooled {
		return fmt.Errorf("cannot release non-pooled instance")
	}

	pm.mu.RLock()
	pool, exists := pm.pools[instance.Name]
	pm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no pool exists for sprite %s", instance.Name)
	}

	return pool.Release(instance)
}

// Update updates all pools and returns expired instances
func (pm *PoolManager) Update(deltaTime time.Duration) []*SpriteInstance {
	pm.mu.RLock()
	pools := make([]*Pool, 0, len(pm.pools))
	for _, pool := range pm.pools {
		pools = append(pools, pool)
	}
	pm.mu.RUnlock()

	expired := make([]*SpriteInstance, 0)
	for _, pool := range pools {
		poolExpired := pool.Update(deltaTime)
		expired = append(expired, poolExpired...)
	}

	return expired
}

// GetStats returns statistics for a specific pool
func (pm *PoolManager) GetStats(spriteName string) (PoolStats, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pool, exists := pm.pools[spriteName]
	if !exists {
		return PoolStats{}, fmt.Errorf("no pool exists for sprite %s", spriteName)
	}

	return pool.GetStats(), nil
}

// GetAllStats returns statistics for all pools
func (pm *PoolManager) GetAllStats() map[string]PoolStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	stats := make(map[string]PoolStats)
	for spriteName, pool := range pm.pools {
		stats[spriteName] = pool.GetStats()
	}

	return stats
}

// RemovePool removes a pool (for cleanup)
func (pm *PoolManager) RemovePool(spriteName string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.pools, spriteName)
}

// Clear removes all pools
func (pm *PoolManager) Clear() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.pools = make(map[string]*Pool)
}
