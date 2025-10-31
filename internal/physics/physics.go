package physics

import (
	"github.com/ByteArena/box2d"
)

// World represents a Box2D physics world
type World struct {
	world              box2d.B2World
	gravity            box2d.B2Vec2
	timeStep           float64
	velocityIterations int
	positionIterations int
}

// BodyType represents the type of physics body
type BodyType int

const (
	StaticBody BodyType = iota
	RigidBody
	KinematicBody
)

// Body represents a physics body with position and velocity
type Body struct {
	body     *box2d.B2Body
	bodyType BodyType
}

// NewWorld creates a new physics world with gravity
func NewWorld(gravityX, gravityY float64) *World {
	gravity := box2d.MakeB2Vec2(gravityX, gravityY)
	world := box2d.MakeB2World(gravity)
	return &World{
		world:              world,
		gravity:            gravity,
		timeStep:           1.0 / 60.0, // 60 FPS
		velocityIterations: 8,
		positionIterations: 3,
	}
}

// Step advances the physics simulation by one time step
func (w *World) Step() {
	w.world.Step(w.timeStep, w.velocityIterations, w.positionIterations)
}

// GetWorld returns the underlying Box2D world (for advanced usage)
func (w *World) GetWorld() *box2d.B2World {
	return &w.world
}

// CreateBody creates a physics body with the given type and position
func (w *World) CreateBody(bodyType BodyType, x, y float64) *Body {
	bodyDef := box2d.MakeB2BodyDef()

	switch bodyType {
	case StaticBody:
		bodyDef.Type = box2d.B2BodyType.B2_staticBody
	case KinematicBody:
		bodyDef.Type = box2d.B2BodyType.B2_kinematicBody
	case RigidBody:
		bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	}

	bodyDef.Position.Set(x, y)
	body := w.world.CreateBody(&bodyDef)

	return &Body{
		body:     body,
		bodyType: bodyType,
	}
}

// CreateStaticBody creates a static (immovable) body
func (w *World) CreateStaticBody(x, y float64) *Body {
	return w.CreateBody(StaticBody, x, y)
}

// CreateDynamicBody creates a dynamic (physics-simulated) body
func (w *World) CreateDynamicBody(x, y float64) *Body {
	return w.CreateBody(RigidBody, x, y)
}

// CreateKinematicBody creates a kinematic (player-controlled) body
func (w *World) CreateKinematicBody(x, y float64) *Body {
	return w.CreateBody(KinematicBody, x, y)
}

// CreateBoxFixture adds a box-shaped fixture to a body
func (b *Body) CreateBoxFixture(width, height float64, density float64) {
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(width/2.0, height/2.0)
	b.body.CreateFixture(&shape, density)
}

// CreateCircleFixture adds a circle-shaped fixture to a body
func (b *Body) CreateCircleFixture(radius float64, density float64) {
	shape := box2d.MakeB2CircleShape()
	shape.M_radius = radius
	b.body.CreateFixture(&shape, density)
}

// SetPosition sets the body's position
func (b *Body) SetPosition(x, y float64) {
	b.body.SetTransform(box2d.MakeB2Vec2(x, y), b.body.GetAngle())
}

// GetPosition returns the body's position
func (b *Body) GetPosition() (x, y float64) {
	pos := b.body.GetPosition()
	return pos.X, pos.Y
}

// SetVelocity sets the body's linear velocity
func (b *Body) SetVelocity(vx, vy float64) {
	b.body.SetLinearVelocity(box2d.MakeB2Vec2(vx, vy))
}

// GetVelocity returns the body's linear velocity
func (b *Body) GetVelocity() (vx, vy float64) {
	vel := b.body.GetLinearVelocity()
	return vel.X, vel.Y
}

// SetAngle sets the body's rotation angle in radians
func (b *Body) SetAngle(angle float64) {
	pos := b.body.GetPosition()
	b.body.SetTransform(pos, angle)
}

// GetAngle returns the body's rotation angle in radians
func (b *Body) GetAngle() float64 {
	return b.body.GetAngle()
}

// ApplyForce applies a force to the body at the given point
func (b *Body) ApplyForce(fx, fy, px, py float64) {
	b.body.ApplyForce(box2d.MakeB2Vec2(fx, fy), box2d.MakeB2Vec2(px, py), true)
}

// ApplyImpulse applies an impulse to the body at the given point
func (b *Body) ApplyImpulse(ix, iy, px, py float64) {
	b.body.ApplyLinearImpulse(box2d.MakeB2Vec2(ix, iy), box2d.MakeB2Vec2(px, py), true)
}

// SetGravityScale sets the gravity scale for the body
func (b *Body) SetGravityScale(scale float64) {
	b.body.SetGravityScale(scale)
}

// Destroy removes the body from the world
func (b *Body) Destroy() {
	w := b.body.GetWorld()
	if w != nil {
		w.DestroyBody(b.body)
	}
}

// GetWorld returns the world this body belongs to
func (b *Body) GetWorld() *box2d.B2World {
	return b.body.GetWorld()
}
