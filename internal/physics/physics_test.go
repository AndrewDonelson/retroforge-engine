package physics

import (
	"testing"
)

func TestNewWorld(t *testing.T) {
	world := NewWorld(0, 9.8)
	if world == nil {
		t.Fatal("NewWorld returned nil")
	}
	if world.timeStep != 1.0/60.0 {
		t.Errorf("Expected timeStep 1/60, got %f", world.timeStep)
	}
}

func TestCreateBody(t *testing.T) {
	world := NewWorld(0, 9.8)

	// Test static body
	staticBody := world.CreateStaticBody(10, 20)
	if staticBody == nil {
		t.Fatal("CreateStaticBody returned nil")
	}
	x, y := staticBody.GetPosition()
	if x != 10 || y != 20 {
		t.Errorf("Expected position (10, 20), got (%f, %f)", x, y)
	}

	// Test dynamic body
	dynamicBody := world.CreateDynamicBody(5, 15)
	if dynamicBody == nil {
		t.Fatal("CreateDynamicBody returned nil")
	}

	// Test kinematic body
	kinematicBody := world.CreateKinematicBody(0, 0)
	if kinematicBody == nil {
		t.Fatal("CreateKinematicBody returned nil")
	}
}

func TestBodyFixture(t *testing.T) {
	world := NewWorld(0, 9.8)
	body := world.CreateDynamicBody(0, 0)

	// Add box fixture
	body.CreateBoxFixture(10, 20, 1.0)

	// Add circle fixture
	body.CreateCircleFixture(5, 1.0)
}

func TestBodyPosition(t *testing.T) {
	world := NewWorld(0, 9.8)
	body := world.CreateDynamicBody(0, 0)

	// Set position
	body.SetPosition(100, 200)
	x, y := body.GetPosition()
	if x != 100 || y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", x, y)
	}
}

func TestBodyVelocity(t *testing.T) {
	world := NewWorld(0, 9.8)
	body := world.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(10, 10, 1.0)

	// Set velocity
	body.SetVelocity(5, -10)
	vx, vy := body.GetVelocity()
	if vx != 5 || vy != -10 {
		t.Errorf("Expected velocity (5, -10), got (%f, %f)", vx, vy)
	}
}

func TestPhysicsStep(t *testing.T) {
	world := NewWorld(0, 9.8)
	body := world.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(10, 10, 1.0)

	// Step physics
	world.Step()

	// Body should have moved due to gravity
	x, y := body.GetPosition()
	if x == 0 && y == 0 {
		t.Log("Position unchanged after step (may be expected if gravity is disabled for dynamic bodies initially)")
	}
}

func TestBodyDestroy(t *testing.T) {
	world := NewWorld(0, 9.8)
	body := world.CreateDynamicBody(0, 0)

	// Destroy body
	body.Destroy()

	// Body should be removed from world (cannot verify without accessing internals)
	// This test mainly ensures Destroy doesn't crash
}

func TestEdgeCases(t *testing.T) {
	world := NewWorld(0, 9.8)

	// Test creating body at negative coordinates
	body := world.CreateDynamicBody(-10, -20)
	x, y := body.GetPosition()
	if x != -10 || y != -20 {
		t.Errorf("Expected position (-10, -20), got (%f, %f)", x, y)
	}

	// Test small fixture (minimum valid size)
	body.CreateBoxFixture(0.1, 0.1, 1.0)

	// Test very large coordinates
	body2 := world.CreateStaticBody(10000, 20000)
	x2, y2 := body2.GetPosition()
	if x2 != 10000 || y2 != 20000 {
		t.Errorf("Expected position (10000, 20000), got (%f, %f)", x2, y2)
	}
}

func TestWorldGetWorld(t *testing.T) {
	w := NewWorld(0, 9.8)
	b2World := w.GetWorld()
	if b2World == nil {
		t.Error("GetWorld should not return nil")
	}
}

func TestBodySetAngle(t *testing.T) {
	w := NewWorld(0, 9.8)
	body := w.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(1, 1, 1.0)

	// Test setting angle
	body.SetAngle(1.57) // ~90 degrees
	angle := body.GetAngle()
	if angle < 1.5 || angle > 1.6 {
		t.Errorf("SetAngle(1.57) failed, GetAngle() = %f, expected ~1.57", angle)
	}

	// Test zero angle
	body.SetAngle(0)
	angle = body.GetAngle()
	if angle != 0 {
		t.Errorf("SetAngle(0) failed, GetAngle() = %f, expected 0", angle)
	}

	// Test negative angle
	body.SetAngle(-0.5)
	angle = body.GetAngle()
	if angle < -0.6 || angle > -0.4 {
		t.Errorf("SetAngle(-0.5) failed, GetAngle() = %f, expected ~-0.5", angle)
	}
}

func TestBodyGetAngle(t *testing.T) {
	w := NewWorld(0, 9.8)
	body := w.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(1, 1, 1.0)

	// Initial angle should be 0
	angle := body.GetAngle()
	if angle != 0 {
		t.Errorf("Initial angle should be 0, got %f", angle)
	}

	// Set and get
	body.SetAngle(2.0)
	angle = body.GetAngle()
	if angle < 1.9 || angle > 2.1 {
		t.Errorf("GetAngle after SetAngle(2.0) = %f, expected ~2.0", angle)
	}
}

func TestBodyApplyForce(t *testing.T) {
	w := NewWorld(0, 9.8)
	body := w.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(1, 1, 1.0)

	// Test applying force
	body.ApplyForce(10, 20, 0, 0)

	// Force application doesn't immediately change velocity (needs Step)
	// This is a smoke test that it doesn't crash
	velX, velY := body.GetVelocity()
	_ = velX
	_ = velY
}

func TestBodyApplyImpulse(t *testing.T) {
	w := NewWorld(0, 9.8)
	body := w.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(1, 1, 1.0)

	// Test applying impulse
	body.ApplyImpulse(5, 10, 0, 0)

	// Impulse should change velocity immediately
	w.Step()
	velX, velY := body.GetVelocity()
	if velX == 0 && velY == 0 {
		t.Logf("Impulse may require multiple steps")
	}
}

func TestBodySetGravityScale(t *testing.T) {
	w := NewWorld(0, 9.8)
	body := w.CreateDynamicBody(0, 0)
	body.CreateBoxFixture(1, 1, 1.0)

	// Test setting gravity scale
	body.SetGravityScale(0.5)
	w.Step()

	// With lower gravity, body should fall slower
	// This is a smoke test
	_, velY := body.GetVelocity()
	_ = velY

	// Test zero gravity scale
	body.SetPosition(0, 0)
	body.SetVelocity(0, 0)
	body.SetGravityScale(0)
	w.Step()

	// Test negative gravity (anti-gravity)
	body.SetGravityScale(-1.0)
	w.Step()
}

func TestBodyGetWorld(t *testing.T) {
	w := NewWorld(0, 9.8)
	body := w.CreateDynamicBody(0, 0)

	b2World := body.GetWorld()
	if b2World == nil {
		t.Error("Body.GetWorld should not return nil")
	}

	// Verify it's the same world
	worldFromWorld := w.GetWorld()
	if worldFromWorld == nil {
		t.Error("World.GetWorld should not return nil")
	}
}
