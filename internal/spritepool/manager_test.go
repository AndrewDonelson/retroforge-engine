package spritepool

import (
	"testing"
	"time"
)

func TestShouldPool(t *testing.T) {
	// Should pool: isUI=false, maxSpawn > 10
	spriteData := createTestSpriteData("bullet", false, 20, 1000)
	if !ShouldPool(spriteData) {
		t.Error("should pool: isUI=false, maxSpawn=20")
	}

	// Should NOT pool: isUI=true
	spriteData = createTestSpriteData("ui_button", true, 20, 1000)
	if ShouldPool(spriteData) {
		t.Error("should NOT pool: isUI=true")
	}

	// Should NOT pool: maxSpawn <= 10
	spriteData = createTestSpriteData("bullet", false, 10, 1000)
	if ShouldPool(spriteData) {
		t.Error("should NOT pool: maxSpawn=10")
	}

	// Should NOT pool: maxSpawn = 0
	spriteData = createTestSpriteData("bullet", false, 0, 1000)
	if ShouldPool(spriteData) {
		t.Error("should NOT pool: maxSpawn=0")
	}

	// Should pool: exactly maxSpawn > 10
	spriteData = createTestSpriteData("bullet", false, 11, 1000)
	if !ShouldPool(spriteData) {
		t.Error("should pool: isUI=false, maxSpawn=11")
	}
}

func TestNewPoolManager(t *testing.T) {
	pm := NewPoolManager()
	if pm == nil {
		t.Fatal("pool manager should not be nil")
	}
	if pm.pools == nil {
		t.Fatal("pools map should not be nil")
	}
	if len(pm.pools) != 0 {
		t.Errorf("expected 0 pools, got %d", len(pm.pools))
	}
}

func TestPoolManagerRegisterSprite(t *testing.T) {
	pm := NewPoolManager()

	// Register eligible sprite
	spriteData := createTestSpriteData("bullet", false, 50, 2000)
	err := pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if !pm.HasPool("bullet") {
		t.Error("pool should exist for bullet")
	}

	// Register ineligible sprite (isUI=true)
	spriteData2 := createTestSpriteData("ui_button", true, 50, 2000)
	err = pm.RegisterSprite("ui_button", spriteData2)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if pm.HasPool("ui_button") {
		t.Error("pool should NOT exist for ui_button (isUI=true)")
	}

	// Register ineligible sprite (maxSpawn <= 10)
	spriteData3 := createTestSpriteData("enemy", false, 5, 2000)
	err = pm.RegisterSprite("enemy", spriteData3)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if pm.HasPool("enemy") {
		t.Error("pool should NOT exist for enemy (maxSpawn <= 10)")
	}

	// Register same sprite twice (should be idempotent)
	err = pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("second register failed: %v", err)
	}

	// Should still have pool
	if !pm.HasPool("bullet") {
		t.Error("pool should still exist after second register")
	}
}

func TestPoolManagerAcquireRelease(t *testing.T) {
	pm := NewPoolManager()

	// Register sprite
	spriteData := createTestSpriteData("bullet", false, 50, 2000)
	err := pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Acquire instance
	instance, err := pm.Acquire("bullet")
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	if instance == nil {
		t.Fatal("instance should not be nil")
	}
	if instance.Name != "bullet" {
		t.Errorf("expected name 'bullet', got %s", instance.Name)
	}
	if !instance.IsPooled {
		t.Error("instance should be pooled")
	}
	if !instance.IsActive {
		t.Error("instance should be active")
	}

	// Release instance
	err = pm.Release(instance)
	if err != nil {
		t.Fatalf("release failed: %v", err)
	}

	// Acquire from non-existent pool
	_, err = pm.Acquire("nonexistent")
	if err == nil {
		t.Error("expected error when acquiring from non-existent pool")
	}
}

func TestPoolManagerUpdate(t *testing.T) {
	pm := NewPoolManager()

	// Register sprite with lifetime
	spriteData := createTestSpriteData("bullet", false, 50, 1000) // 1 second
	err := pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Acquire instances
	instances := make([]*SpriteInstance, 0, 3)
	for i := 0; i < 3; i++ {
		instance, err := pm.Acquire("bullet")
		if err != nil {
			t.Fatalf("acquire failed: %v", err)
		}
		instances = append(instances, instance)
	}

	// Update - nothing should expire yet
	expired := pm.Update(500 * time.Millisecond)
	if len(expired) != 0 {
		t.Errorf("expected 0 expired, got %d", len(expired))
	}

	// Update past lifetime
	expired = pm.Update(600 * time.Millisecond)
	if len(expired) != 3 {
		t.Errorf("expected 3 expired, got %d", len(expired))
	}

	// Clean up expired bullet instances first
	for _, instance := range expired {
		pm.Release(instance)
	}

	// Register sprite without lifetime
	spriteData2 := createTestSpriteData("enemy", false, 50, 0)
	err = pm.RegisterSprite("enemy", spriteData2)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	instance, err := pm.Acquire("enemy")
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	// Update - enemy should not expire (no lifetime)
	// Note: This might still include expired bullets from previous update if not released
	// So we create a fresh pool manager for this test
	pm2 := NewPoolManager()
	err = pm2.RegisterSprite("enemy", spriteData2)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	instance2, err := pm2.Acquire("enemy")
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	// Update - enemy should not expire (no lifetime)
	expired2 := pm2.Update(10 * time.Second)
	if len(expired2) != 0 {
		t.Errorf("expected 0 expired (no lifetime), got %d", len(expired2))
	}

	// Cleanup
	pm.Release(instance)
	pm2.Release(instance2)
}

func TestPoolManagerGetStats(t *testing.T) {
	pm := NewPoolManager()

	// Register sprite
	spriteData := createTestSpriteData("bullet", false, 50, 2000)
	err := pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Get stats for non-existent pool
	_, err = pm.GetStats("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent pool")
	}

	// Get stats
	stats, err := pm.GetStats("bullet")
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}
	if stats.SpriteName != "bullet" {
		t.Errorf("expected sprite name 'bullet', got %s", stats.SpriteName)
	}

	// Acquire some instances
	for i := 0; i < 5; i++ {
		_, err := pm.Acquire("bullet")
		if err != nil {
			t.Fatalf("acquire failed: %v", err)
		}
	}

	stats, err = pm.GetStats("bullet")
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}
	if stats.Active != 5 {
		t.Errorf("expected 5 active, got %d", stats.Active)
	}
}

func TestPoolManagerGetAllStats(t *testing.T) {
	pm := NewPoolManager()

	// No pools yet
	stats := pm.GetAllStats()
	if len(stats) != 0 {
		t.Errorf("expected 0 pools, got %d", len(stats))
	}

	// Register multiple sprites
	spriteData1 := createTestSpriteData("bullet", false, 50, 2000)
	spriteData2 := createTestSpriteData("particle", false, 100, 500)

	err := pm.RegisterSprite("bullet", spriteData1)
	if err != nil {
		t.Fatalf("register bullet failed: %v", err)
	}
	err = pm.RegisterSprite("particle", spriteData2)
	if err != nil {
		t.Fatalf("register particle failed: %v", err)
	}

	stats = pm.GetAllStats()
	if len(stats) != 2 {
		t.Errorf("expected 2 pools, got %d", len(stats))
	}

	if _, exists := stats["bullet"]; !exists {
		t.Error("bullet pool stats should exist")
	}
	if _, exists := stats["particle"]; !exists {
		t.Error("particle pool stats should exist")
	}
}

func TestPoolManagerRemovePool(t *testing.T) {
	pm := NewPoolManager()

	// Register sprite
	spriteData := createTestSpriteData("bullet", false, 50, 2000)
	err := pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if !pm.HasPool("bullet") {
		t.Error("pool should exist")
	}

	// Remove pool
	pm.RemovePool("bullet")

	if pm.HasPool("bullet") {
		t.Error("pool should not exist after removal")
	}
}

func TestPoolManagerClear(t *testing.T) {
	pm := NewPoolManager()

	// Register multiple sprites
	spriteData1 := createTestSpriteData("bullet", false, 50, 2000)
	spriteData2 := createTestSpriteData("particle", false, 100, 500)

	err := pm.RegisterSprite("bullet", spriteData1)
	if err != nil {
		t.Fatalf("register bullet failed: %v", err)
	}
	err = pm.RegisterSprite("particle", spriteData2)
	if err != nil {
		t.Fatalf("register particle failed: %v", err)
	}

	if len(pm.GetAllStats()) != 2 {
		t.Error("expected 2 pools before clear")
	}

	// Clear all pools
	pm.Clear()

	if len(pm.GetAllStats()) != 0 {
		t.Error("expected 0 pools after clear")
	}
	if pm.HasPool("bullet") {
		t.Error("bullet pool should not exist after clear")
	}
	if pm.HasPool("particle") {
		t.Error("particle pool should not exist after clear")
	}
}

func TestPoolManagerConcurrentAccess(t *testing.T) {
	pm := NewPoolManager()

	spriteData := createTestSpriteData("bullet", false, 100, 2000)
	err := pm.RegisterSprite("bullet", spriteData)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Test concurrent acquire/release
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			instance, err := pm.Acquire("bullet")
			if err != nil {
				t.Errorf("concurrent acquire failed: %v", err)
				return
			}
			time.Sleep(1 * time.Millisecond)
			err = pm.Release(instance)
			if err != nil {
				t.Errorf("concurrent release failed: %v", err)
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Pool should be in valid state
	stats := pm.GetAllStats()
	if len(stats) != 1 {
		t.Errorf("expected 1 pool, got %d", len(stats))
	}
}

func TestPoolManagerReleaseNonPooled(t *testing.T) {
	pm := NewPoolManager()

	// Try to release non-pooled instance
	nonPooled := &SpriteInstance{
		Name:       "bullet",
		IsPooled:   false,
		IsActive:   true,
		CustomData: make(map[string]interface{}),
	}

	err := pm.Release(nonPooled)
	if err == nil {
		t.Error("expected error when releasing non-pooled instance")
	}
}
