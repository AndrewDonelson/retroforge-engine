package spritepool

import (
	"testing"
	"time"
)

func TestPoolEdgeCases(t *testing.T) {
	// Test with maxSpawn exactly 11 (boundary)
	spriteData := createTestSpriteData("boundary", false, 11, 1000)
	if !ShouldPool(spriteData) {
		t.Error("should pool when maxSpawn=11")
	}

	// Test with maxSpawn exactly 10 (boundary, should NOT pool)
	spriteData = createTestSpriteData("boundary2", false, 10, 1000)
	if ShouldPool(spriteData) {
		t.Error("should NOT pool when maxSpawn=10")
	}

	// Test with maxSpawn exactly 1
	spriteData = createTestSpriteData("single", false, 1, 1000)
	if ShouldPool(spriteData) {
		t.Error("should NOT pool when maxSpawn=1")
	}

	// Test very large maxSpawn
	spriteData = createTestSpriteData("large", false, 10000, 1000)
	if !ShouldPool(spriteData) {
		t.Error("should pool when maxSpawn=10000")
	}
}

func TestPoolInitialSize(t *testing.T) {
	// Test initial size calculation (should be 50% of max, capped at 50)
	spriteData := createTestSpriteData("test", false, 100, 1000)
	pool := NewPool("test", spriteData, 50, 100) // Initial size 50
	if len(pool.available) != 50 {
		t.Errorf("expected initial size 50, got %d", len(pool.available))
	}

	// Test with small maxSpawn
	spriteData2 := createTestSpriteData("small", false, 20, 1000)
	pool2 := NewPool("small", spriteData2, 10, 20) // Initial size 10 (50% of 20)
	if len(pool2.available) != 10 {
		t.Errorf("expected initial size 10, got %d", len(pool2.available))
	}

	// Test with very small maxSpawn
	spriteData3 := createTestSpriteData("tiny", false, 15, 1000)
	pool3 := NewPool("tiny", spriteData3, 7, 15) // Initial size 7 (50% of 15, rounded down)
	if len(pool3.available) != 7 {
		t.Errorf("expected initial size 7, got %d", len(pool3.available))
	}
}

func TestPoolGrowth(t *testing.T) {
	spriteData := createTestSpriteData("test", false, 100, 1000)
	pool := NewPool("test", spriteData, 1, 100)

	// Acquire all available
	_, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	// Acquire more to trigger growth
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire after growth failed: %v", err)
	}
	if instance == nil {
		t.Fatal("instance should not be nil")
	}

	stats := pool.GetStats()
	if stats.Active != 2 {
		t.Errorf("expected 2 active after growth, got %d", stats.Active)
	}
}

func TestPoolLifetimeEdgeCases(t *testing.T) {
	// Test with very short lifetime
	spriteData := createTestSpriteData("short", false, 50, 1) // 1ms
	pool := NewPool("short", spriteData, 5, 50)

	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	// Update past lifetime
	expired := pool.Update(2 * time.Millisecond)
	if len(expired) != 1 {
		t.Errorf("expected 1 expired, got %d", len(expired))
	}

	// Test with very long lifetime
	spriteData2 := createTestSpriteData("long", false, 50, 3600000) // 1 hour
	pool2 := NewPool("long", spriteData2, 5, 50)

	instance2, err := pool2.Acquire()
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	// Update for 1 second - should not expire
	expired2 := pool2.Update(1 * time.Second)
	if len(expired2) != 0 {
		t.Errorf("expected 0 expired (long lifetime), got %d", len(expired2))
	}

	pm := NewPoolManager()
	pm.Release(instance)
	pm.Release(instance2)
}

func TestPoolManagerMultiplePools(t *testing.T) {
	pm := NewPoolManager()

	// Register multiple eligible sprites
	spriteData1 := createTestSpriteData("bullet", false, 50, 1000)
	spriteData2 := createTestSpriteData("particle", false, 100, 500)
	spriteData3 := createTestSpriteData("enemy_bullet", false, 25, 2000)

	err := pm.RegisterSprite("bullet", spriteData1)
	if err != nil {
		t.Fatalf("register bullet failed: %v", err)
	}
	err = pm.RegisterSprite("particle", spriteData2)
	if err != nil {
		t.Fatalf("register particle failed: %v", err)
	}
	err = pm.RegisterSprite("enemy_bullet", spriteData3)
	if err != nil {
		t.Fatalf("register enemy_bullet failed: %v", err)
	}

	// All should have pools
	if !pm.HasPool("bullet") {
		t.Error("bullet pool should exist")
	}
	if !pm.HasPool("particle") {
		t.Error("particle pool should exist")
	}
	if !pm.HasPool("enemy_bullet") {
		t.Error("enemy_bullet pool should exist")
	}

	// Acquire from different pools
	bullet, err := pm.Acquire("bullet")
	if err != nil {
		t.Fatalf("acquire bullet failed: %v", err)
	}

	particle, err := pm.Acquire("particle")
	if err != nil {
		t.Fatalf("acquire particle failed: %v", err)
	}

	if bullet.Name != "bullet" {
		t.Errorf("expected bullet name, got %s", bullet.Name)
	}
	if particle.Name != "particle" {
		t.Errorf("expected particle name, got %s", particle.Name)
	}

	// Release
	pm.Release(bullet)
	pm.Release(particle)
}

func TestPoolConcurrentOperations(t *testing.T) {
	pm := NewPoolManager()

	spriteData := createTestSpriteData("test", false, 200, 2000)
	err := pm.RegisterSprite("test", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Concurrent acquire/release operations
	const numGoroutines = 50
	const opsPerGoroutine = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < opsPerGoroutine; j++ {
				instance, err := pm.Acquire("test")
				if err != nil {
					t.Errorf("concurrent acquire failed: %v", err)
					return
				}
				time.Sleep(100 * time.Microsecond)
				err = pm.Release(instance)
				if err != nil {
					t.Errorf("concurrent release failed: %v", err)
					return
				}
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Pool should be in valid state
	stats, err := pm.GetStats("test")
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}

	// All instances should be released
	if stats.Active != 0 {
		t.Errorf("expected 0 active after all releases, got %d", stats.Active)
	}
}

func TestPoolEmptyRelease(t *testing.T) {
	spriteData := createTestSpriteData("test", false, 50, 1000)
	pool := NewPool("test", spriteData, 0, 50) // Start with empty pool

	// Acquire should trigger growth (grows by 25% of max = 12.5, rounded to 13)
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire from empty pool failed: %v", err)
	}

	// Release
	err = pool.Release(instance)
	if err != nil {
		t.Fatalf("release failed: %v", err)
	}

	// Growth size is 25% of 50 = 12.5, converted to int = 12
	// We create 12 instances, acquire 1, leaving 11 available
	// After release, we should have 12 available (11 from growth + 1 released)
	if len(pool.available) != 12 {
		t.Errorf("expected 12 available after release (growth created 12, 1 acquired, 1 released), got %d", len(pool.available))
	}
}

func TestPoolExhaustionWithGrowth(t *testing.T) {
	spriteData := createTestSpriteData("test", false, 10, 1000)
	pool := NewPool("test", spriteData, 1, 10) // Small initial size

	// Acquire all up to max
	instances := make([]*SpriteInstance, 0, 10)
	for i := 0; i < 10; i++ {
		instance, err := pool.Acquire()
		if err != nil && i < 9 {
			t.Fatalf("acquire %d failed: %v", i, err)
		}
		instances = append(instances, instance)
	}

	// Try to acquire one more - should create overflow
	overflow, err := pool.Acquire()
	if err == nil {
		t.Error("expected error when pool exhausted")
	}
	if overflow == nil {
		t.Fatal("overflow instance should not be nil")
	}
	if overflow.IsPooled {
		t.Error("overflow instance should not be pooled")
	}

	// Release some instances
	for i := 0; i < 5; i++ {
		err := pool.Release(instances[i])
		if err != nil {
			t.Fatalf("release failed: %v", err)
		}
	}

	// Should be able to acquire from pool again
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire after release failed: %v", err)
	}
	if !instance.IsPooled {
		t.Error("instance should be pooled")
	}

	pool.Release(instance)
}
