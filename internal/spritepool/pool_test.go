package spritepool

import (
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	pool := NewPool("bullet", spriteData, 10, 100)

	if pool == nil {
		t.Fatal("pool should not be nil")
	}

	if pool.spriteName != "bullet" {
		t.Errorf("expected sprite name 'bullet', got %s", pool.spriteName)
	}

	if pool.maxSize != 100 {
		t.Errorf("expected max size 100, got %d", pool.maxSize)
	}

	if len(pool.available) != 10 {
		t.Errorf("expected 10 available instances, got %d", len(pool.available))
	}

	if len(pool.active) != 0 {
		t.Errorf("expected 0 active instances, got %d", len(pool.active))
	}
}

func TestPoolAcquire(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	pool := NewPool("bullet", spriteData, 5, 100)

	// Acquire all available instances
	instances := make([]*SpriteInstance, 0, 5)
	for i := 0; i < 5; i++ {
		instance, err := pool.Acquire()
		if err != nil {
			t.Fatalf("acquire failed: %v", err)
		}
		if instance == nil {
			t.Fatal("instance should not be nil")
		}
		if !instance.IsPooled {
			t.Error("instance should be pooled")
		}
		if !instance.IsActive {
			t.Error("instance should be active")
		}
		if instance.Name != "bullet" {
			t.Errorf("expected name 'bullet', got %s", instance.Name)
		}
		instances = append(instances, instance)
	}

	stats := pool.GetStats()
	if stats.Active != 5 {
		t.Errorf("expected 5 active instances, got %d", stats.Active)
	}
	if stats.Available != 0 {
		t.Errorf("expected 0 available instances, got %d", stats.Available)
	}

	// Acquire more to trigger growth
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire after growth failed: %v", err)
	}
	if instance == nil {
		t.Fatal("instance should not be nil")
	}
	if !instance.IsPooled {
		t.Error("instance should be pooled")
	}
}

func TestPoolRelease(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	pool := NewPool("bullet", spriteData, 5, 100)

	// Acquire instance
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	// Set some state
	instance.X = 100
	instance.Y = 200
	instance.Age = time.Second
	instance.CustomData["test"] = "value"

	// Release instance
	err = pool.Release(instance)
	if err != nil {
		t.Fatalf("release failed: %v", err)
	}

	// Check state was reset
	if instance.IsActive {
		t.Error("instance should not be active after release")
	}
	if instance.X != 0 {
		t.Errorf("expected X=0 after release, got %f", instance.X)
	}
	if instance.Y != 0 {
		t.Errorf("expected Y=0 after release, got %f", instance.Y)
	}
	if instance.Age != 0 {
		t.Errorf("expected Age=0 after release, got %v", instance.Age)
	}
	if len(instance.CustomData) != 0 {
		t.Errorf("expected empty CustomData after release, got %d items", len(instance.CustomData))
	}

	stats := pool.GetStats()
	if stats.Active != 0 {
		t.Errorf("expected 0 active instances, got %d", stats.Active)
	}
	if stats.Available != 5 {
		t.Errorf("expected 5 available instances, got %d", stats.Available)
	}
}

func TestPoolExhaustion(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 10, 2000)
	pool := NewPool("bullet", spriteData, 5, 10)

	// Acquire all instances (including growing to max)
	instances := make([]*SpriteInstance, 0, 10)
	for i := 0; i < 10; i++ {
		instance, err := pool.Acquire()
		if err != nil && i < 9 { // Last one might create overflow
			t.Fatalf("acquire %d failed: %v", i, err)
		}
		instances = append(instances, instance)
	}

	// Try to acquire one more - should create overflow instance
	overflowInstance, err := pool.Acquire()
	if err == nil {
		t.Error("expected error when pool is exhausted")
	}
	if overflowInstance == nil {
		t.Fatal("overflow instance should not be nil")
	}
	if overflowInstance.IsPooled {
		t.Error("overflow instance should not be pooled")
	}
	if !overflowInstance.IsActive {
		t.Error("overflow instance should be active")
	}
}

func TestPoolUpdate(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 1000) // 1 second lifetime
	pool := NewPool("bullet", spriteData, 5, 100)

	// Acquire some instances
	instances := make([]*SpriteInstance, 0, 3)
	for i := 0; i < 3; i++ {
		instance, err := pool.Acquire()
		if err != nil {
			t.Fatalf("acquire failed: %v", err)
		}
		instances = append(instances, instance)
	}

	// Update with time less than lifetime - nothing should expire
	expired := pool.Update(500 * time.Millisecond)
	if len(expired) != 0 {
		t.Errorf("expected 0 expired instances, got %d", len(expired))
	}

	// Update past lifetime - all should expire
	expired = pool.Update(600 * time.Millisecond)
	if len(expired) != 3 {
		t.Errorf("expected 3 expired instances, got %d", len(expired))
	}

	for _, instance := range expired {
		if instance.Name != "bullet" {
			t.Errorf("expected expired instance name 'bullet', got %s", instance.Name)
		}
	}
}

func TestPoolUpdateNoLifetime(t *testing.T) {
	spriteData := createTestSpriteData("enemy", false, 100, 0) // No lifetime
	pool := NewPool("enemy", spriteData, 5, 100)

	// Acquire instances
	for i := 0; i < 3; i++ {
		_, err := pool.Acquire()
		if err != nil {
			t.Fatalf("acquire failed: %v", err)
		}
	}

	// Update should return nothing since no lifetime
	expired := pool.Update(10 * time.Second)
	if len(expired) != 0 {
		t.Errorf("expected 0 expired instances (no lifetime), got %d", len(expired))
	}
}

func TestPoolGetStats(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	pool := NewPool("bullet", spriteData, 10, 100)

	stats := pool.GetStats()
	if stats.SpriteName != "bullet" {
		t.Errorf("expected sprite name 'bullet', got %s", stats.SpriteName)
	}
	if stats.Available != 10 {
		t.Errorf("expected 10 available, got %d", stats.Available)
	}
	if stats.Active != 0 {
		t.Errorf("expected 0 active, got %d", stats.Active)
	}
	if stats.MaxSize != 100 {
		t.Errorf("expected max size 100, got %d", stats.MaxSize)
	}
	if stats.Utilization != 0.0 {
		t.Errorf("expected utilization 0.0, got %f", stats.Utilization)
	}

	// Acquire some instances
	for i := 0; i < 5; i++ {
		_, err := pool.Acquire()
		if err != nil {
			t.Fatalf("acquire failed: %v", err)
		}
	}

	stats = pool.GetStats()
	if stats.Active != 5 {
		t.Errorf("expected 5 active, got %d", stats.Active)
	}
	if stats.Available != 5 {
		t.Errorf("expected 5 available, got %d", stats.Available)
	}
	if stats.Utilization != 0.05 {
		t.Errorf("expected utilization 0.05, got %f", stats.Utilization)
	}
}

func TestPoolReleaseNonPooled(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	pool := NewPool("bullet", spriteData, 5, 100)

	// Create non-pooled instance
	nonPooled := &SpriteInstance{
		Name:       "bullet",
		IsPooled:   false,
		IsActive:   true,
		Data:       spriteData,
		CustomData: make(map[string]interface{}),
	}

	// Try to release - should fail
	err := pool.Release(nonPooled)
	if err == nil {
		t.Error("expected error when releasing non-pooled instance")
	}
}

func TestPoolReleaseInactive(t *testing.T) {
	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	pool := NewPool("bullet", spriteData, 5, 100)

	// Acquire and release instance
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	err = pool.Release(instance)
	if err != nil {
		t.Fatalf("first release failed: %v", err)
	}

	// Try to release again - should fail
	err = pool.Release(instance)
	if err == nil {
		t.Error("expected error when releasing inactive instance")
	}
}
