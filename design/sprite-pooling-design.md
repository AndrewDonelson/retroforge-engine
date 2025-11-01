# Sprite Pooling System Design Document

## Overview

This document describes the automatic sprite pooling system for the 2D game engine. The pooling system provides transparent object reuse for frequently created/destroyed sprites, improving performance by reducing garbage collection pressure and memory allocation overhead.

### Goals

- **Developer-friendly**: Pooling happens automatically based on sprite configuration
- **Transparent**: Same API for pooled and non-pooled sprites
- **Performance**: Reduce GC pauses and memory allocation overhead
- **Configurable**: Engine-level thresholds with sprite-level overrides

### Key Principles

1. UI elements are never pooled
2. Only sprites with short lifetimes and high spawn counts are pooled
3. Developers use the same API regardless of pooling
4. Pooling decisions are based on static sprite metadata

---

## Sprite Configuration

Each sprite class defines pooling metadata through static configuration properties:

### Required Properties

| Property | Type | Description |
|----------|------|-------------|
| `isUI` | boolean | If true, sprite is never pooled |
| `lifetime` | number | Expected lifespan in seconds (0 = indefinite) |
| `maximumCount` | number | Maximum expected concurrent instances |

### Example

```javascript
class Bullet extends Sprite {
  static config = {
    isUI: false,
    lifetime: 2.0,
    maximumCount: 100
  }
}

class UIButton extends Sprite {
  static config = {
    isUI: true,
    lifetime: 0,
    maximumCount: 10
  }
}
```

---

## Pooling Decision Logic

### Auto-Pool Criteria

A sprite type is automatically pooled if **ALL** conditions are met:

1. `isUI === false`
2. `lifetime > 0` AND `lifetime < POOL_LIFETIME_THRESHOLD`
3. `maximumCount >= POOL_COUNT_THRESHOLD`

### Engine-Level Thresholds

Configurable at engine initialization:

| Threshold | Default | Description |
|-----------|---------|-------------|
| `POOL_LIFETIME_THRESHOLD` | 5.0 seconds | Maximum lifetime for pooling eligibility |
| `POOL_COUNT_THRESHOLD` | 20 | Minimum count for pooling eligibility |

### Rationale

- **Lifetime check**: Sprites with longer lifetimes spawn less frequently, reducing pooling benefits
- **Count check**: Low-count sprites don't justify pooling overhead
- **UI exclusion**: UI elements often have complex state and irregular lifecycles

---

## Architecture

### Core Components

```
┌─────────────────────────────────────────┐
│         Engine / SpriteManager          │
│  - createSprite()                       │
│  - registerSpriteType()                 │
└──────────────┬──────────────────────────┘
               │
               │ delegates pooling
               ▼
┌─────────────────────────────────────────┐
│          PoolManager                    │
│  - pools: Map<string, Pool>             │
│  - shouldPool()                         │
│  - acquire()                            │
│  - release()                            │
└──────────────┬──────────────────────────┘
               │
               │ manages multiple
               ▼
┌─────────────────────────────────────────┐
│            Pool                         │
│  - available: Sprite[]                  │
│  - active: Set<Sprite>                  │
│  - maxSize: number                      │
│  - acquire()                            │
│  - release()                            │
│  - grow()                               │
└─────────────────────────────────────────┘
```

### PoolManager

Singleton responsible for managing all sprite pools.

**Responsibilities:**
- Determine if a sprite type should be pooled
- Create and manage Pool instances per sprite type
- Route acquire/release requests to appropriate pools
- Track pooling statistics for debugging

### Pool

Individual pool for a specific sprite type.

**Responsibilities:**
- Store available and active sprite instances
- Grow dynamically when needed (up to maxSize)
- Reset sprite state when releasing back to pool
- Enforce maximum count limits

---

## Implementation Details

### 1. Sprite Type Registration

When a sprite type is registered with the engine:

```javascript
registerSpriteType(SpriteClass) {
  const config = SpriteClass.config;
  
  if (this.poolManager.shouldPool(config)) {
    // Calculate initial pool size
    const initialSize = Math.min(
      Math.ceil(config.maximumCount * 0.5), // Start at 50% of max
      50 // Cap initial allocation at 50
    );
    
    this.poolManager.createPool(
      SpriteClass.name,
      SpriteClass,
      initialSize,
      config.maximumCount
    );
    
    console.log(`[Pool] Created pool for ${SpriteClass.name}: ${initialSize}/${config.maximumCount}`);
  }
}
```

### 2. Sprite Creation (Transparent API)

Developers use the same API regardless of pooling:

```javascript
createSprite(typeName, x, y, ...args) {
  const SpriteClass = this.spriteTypes[typeName];
  let sprite;
  
  if (this.poolManager.hasPool(typeName)) {
    // Get from pool
    sprite = this.poolManager.acquire(typeName);
    sprite.init(x, y, ...args); // Reset/initialize
  } else {
    // Create new instance
    sprite = new SpriteClass(x, y, ...args);
  }
  
  this.activeSprites.add(sprite);
  return sprite;
}
```

### 3. Automatic Lifetime Management

Engine update loop handles automatic sprite destruction:

```javascript
update(deltaTime) {
  for (const sprite of this.activeSprites) {
    sprite.update(deltaTime);
    
    // Check lifetime
    if (sprite.lifetime > 0) {
      sprite.age += deltaTime;
      
      if (sprite.age >= sprite.lifetime) {
        this.destroySprite(sprite);
      }
    }
  }
}

destroySprite(sprite) {
  this.activeSprites.delete(sprite);
  
  if (sprite.isPooled) {
    this.poolManager.release(sprite.constructor.name, sprite);
  } else {
    sprite.destroy();
  }
}
```

### 4. Sprite State Reset

When returning to pool, sprites must be reset to default state:

```javascript
class Pool {
  release(sprite) {
    // Reset common properties
    sprite.age = 0;
    sprite.active = false;
    sprite.visible = false;
    sprite.x = 0;
    sprite.y = 0;
    sprite.velocityX = 0;
    sprite.velocityY = 0;
    
    // Call sprite-specific reset
    if (sprite.onPoolRelease) {
      sprite.onPoolRelease();
    }
    
    this.active.delete(sprite);
    this.available.push(sprite);
  }
}
```

Developers can define custom reset logic:

```javascript
class Bullet extends Sprite {
  onPoolRelease() {
    this.damage = 0;
    this.piercing = false;
    // Clear any references
    this.owner = null;
  }
}
```

### 5. Dynamic Pool Growth

Pools grow automatically when exhausted:

```javascript
class Pool {
  acquire() {
    if (this.available.length === 0) {
      if (this.active.size < this.maxSize) {
        // Grow pool
        const growthSize = Math.min(
          Math.ceil(this.maxSize * 0.25), // Grow by 25%
          this.maxSize - this.active.size
        );
        
        for (let i = 0; i < growthSize; i++) {
          this.available.push(new this.SpriteClass());
        }
        
        console.warn(`[Pool] Grew ${this.typeName} pool by ${growthSize}`);
      } else {
        // Handle overflow (see Edge Cases)
        return this.handleOverflow();
      }
    }
    
    const sprite = this.available.pop();
    sprite.isPooled = true;
    this.active.add(sprite);
    return sprite;
  }
}
```

---

## Edge Cases

### 1. Pool Exhaustion (exceeding maximumCount)

**Scenario**: More sprites requested than maximumCount.

**Solution Options**:

**Option A: Overflow Creation (Recommended)**
```javascript
handleOverflow() {
  console.warn(`[Pool] ${this.typeName} pool exhausted, creating non-pooled instance`);
  const sprite = new this.SpriteClass();
  sprite.isPooled = false;
  return sprite;
}
```

**Option B: Strict Limit**
```javascript
handleOverflow() {
  console.error(`[Pool] ${this.typeName} pool exhausted, cannot create sprite`);
  return null;
}
```

**Option C: Reuse Oldest**
```javascript
handleOverflow() {
  const oldest = this.findOldestActive();
  this.release(oldest);
  return this.acquire();
}
```

**Recommendation**: Option A with debug warnings to help developers tune maximumCount.

### 2. Manual Destruction

Developers may call `sprite.destroy()` manually:

```javascript
class Sprite {
  destroy() {
    if (this.isPooled) {
      // Return to pool instead of destroying
      engine.poolManager.release(this.constructor.name, this);
    } else {
      // Normal cleanup
      this.cleanup();
      // Remove references, listeners, etc.
    }
  }
}
```

### 3. Long-Running Games

For games that run for hours, pools may accumulate unused memory.

**Solution**: Implement pool shrinking:

```javascript
class Pool {
  shrink() {
    const targetSize = Math.max(
      this.active.size,
      Math.ceil(this.maxSize * 0.25)
    );
    
    while (this.available.length > targetSize) {
      const sprite = this.available.pop();
      sprite.destroy();
    }
  }
}

// Call periodically (e.g., between levels)
poolManager.shrinkAll();
```

### 4. Sprite Type Inheritance

If sprite types inherit from pooled parents:

```javascript
class Enemy extends Sprite {
  static config = { isUI: false, lifetime: 0, maximumCount: 50 }
}

class FastEnemy extends Enemy {
  // Inherits config unless overridden
  static config = { ...Enemy.config, maximumCount: 30 }
}
```

Each subclass gets its own pool to avoid type confusion.

---

## Performance Considerations

### Memory Usage

**Pool overhead**:
- Each pooled sprite type uses memory for maximum count
- Formula: `sizeof(Sprite) × maximumCount` per type
- Example: 100 bullets × 1KB each = 100KB reserved

**Mitigation**:
- Start with smaller initial pool sizes (50% of max)
- Implement pool shrinking for long-running games
- Profile actual concurrent sprite counts and adjust maximumCount

### CPU Impact

**Pool acquisition**: O(1) - pop from array  
**Pool release**: O(1) - push to array  
**Traditional creation**: O(1) but triggers GC  
**Traditional destruction**: O(1) but creates garbage

**Net benefit**: Reduced GC frequency outweighs minimal pooling overhead.

### Benchmarking

Expected improvements for typical scenarios:

| Scenario | Spawn Rate | GC Improvement | Frame Time Improvement |
|----------|-----------|----------------|------------------------|
| Bullet hell | 100/sec | 40-60% | 10-20% |
| Particle systems | 500/sec | 50-70% | 15-30% |
| Typical game | 30/sec | 20-30% | 5-10% |
| Low spawn rate | <10/sec | 0-5% | 0-2% |

---

## API Reference

### Engine Configuration

```javascript
const engine = new Engine({
  pooling: {
    enabled: true,
    lifetimeThreshold: 5.0,    // seconds
    countThreshold: 20,         // instances
    debugLogging: false
  }
});
```

### Sprite Configuration

```javascript
class MySprite extends Sprite {
  static config = {
    isUI: boolean,           // Required
    lifetime: number,        // Required (0 = indefinite)
    maximumCount: number     // Required
  }
  
  // Optional: Custom reset logic
  onPoolRelease() {
    // Reset sprite-specific state
  }
}
```

### Developer API

```javascript
// Create sprite (pooling is transparent)
const bullet = engine.createSprite('Bullet', x, y, angle, speed);

// Manual destruction (automatically returns to pool if pooled)
bullet.destroy();

// Check if sprite is pooled
if (bullet.isPooled) {
  console.log('This sprite came from a pool');
}

// Pool management (advanced)
engine.poolManager.getStats('Bullet');
// Returns: { available: 45, active: 55, maxSize: 100 }

engine.poolManager.shrinkAll(); // Between levels
```

---

## Debug Tools

### Pool Statistics

```javascript
class PoolManager {
  getStats(typeName) {
    const pool = this.pools.get(typeName);
    return {
      available: pool.available.length,
      active: pool.active.size,
      maxSize: pool.maxSize,
      utilization: pool.active.size / pool.maxSize
    };
  }
  
  getAllStats() {
    const stats = {};
    for (const [typeName, pool] of this.pools) {
      stats[typeName] = this.getStats(typeName);
    }
    return stats;
  }
}
```

### Debug UI

Display pool statistics in development mode:

```
Pools:
  Bullet:    45/100 available  (55% utilization)
  Particle:  150/200 available (25% utilization)
  Enemy:     8/50 available    (84% utilization) ⚠️
```

### Warnings

- **High utilization** (>80%): Consider increasing maximumCount
- **Overflow events**: Log when non-pooled instances created
- **Low utilization** (<10%): Consider decreasing maximumCount

---

## Implementation Checklist

- [ ] Create PoolManager class
- [ ] Create Pool class with acquire/release methods
- [ ] Implement shouldPool() decision logic
- [ ] Add pool creation during sprite type registration
- [ ] Modify createSprite() to check for pools
- [ ] Add automatic lifetime management in update loop
- [ ] Implement sprite state reset in release()
- [ ] Add overflow handling strategy
- [ ] Implement dynamic pool growth
- [ ] Add pool shrinking for long-running games
- [ ] Create debug statistics API
- [ ] Add console logging for debugging
- [ ] Write unit tests for pool operations
- [ ] Profile performance improvements
- [ ] Document sprite configuration requirements
- [ ] Create example sprite classes

---

## Example Usage

### Complete Example

```javascript
// Define sprite types
class Bullet extends Sprite {
  static config = {
    isUI: false,
    lifetime: 3.0,
    maximumCount: 200
  }
  
  constructor() {
    super();
    this.damage = 10;
    this.speed = 500;
  }
  
  onPoolRelease() {
    this.damage = 10;
    this.speed = 500;
    this.owner = null;
  }
}

class Explosion extends Sprite {
  static config = {
    isUI: false,
    lifetime: 0.5,
    maximumCount: 100
  }
  
  onPoolRelease() {
    this.animation.reset();
    this.scale = 1.0;
  }
}

class Boss extends Sprite {
  static config = {
    isUI: false,
    lifetime: 0,        // Lives until killed
    maximumCount: 1     // Below threshold, not pooled
  }
}

// Initialize engine with pooling
const engine = new Engine({
  pooling: {
    enabled: true,
    lifetimeThreshold: 5.0,
    countThreshold: 20
  }
});

// Register sprite types (pools created automatically)
engine.registerSpriteType(Bullet);   // Will be pooled
engine.registerSpriteType(Explosion); // Will be pooled
engine.registerSpriteType(Boss);      // Won't be pooled

// Use normally in game code
function shootBullet() {
  const bullet = engine.createSprite('Bullet', x, y);
  bullet.angle = playerAngle;
  // Automatically returns to pool after 3 seconds
}

function createExplosion() {
  const explosion = engine.createSprite('Explosion', x, y);
  // Automatically returns to pool after 0.5 seconds
}

// Check pool status (debug)
console.log(engine.poolManager.getAllStats());
```

---

## Future Enhancements

### Potential Additions

1. **Warm-up phase**: Pre-allocate pools during loading screen
2. **Priority system**: Different pools can have different priorities for memory limits
3. **Per-sprite overrides**: Allow runtime override of pooling behavior
4. **Pool presets**: Common configurations (bullets, particles, enemies)
5. **Analytics**: Track pool misses, overflow events, and utilization over time
6. **Adaptive sizing**: Automatically adjust maximumCount based on usage patterns
7. **Cross-scene pooling**: Share pools across different game scenes/levels

### Performance Monitoring

Add telemetry for production builds:
- Pool hit/miss rates
- Average/max concurrent sprite counts
- GC pause frequency and duration
- Memory usage trends

---

## Conclusion

This automatic sprite pooling system provides significant performance benefits for frequently spawned sprites while maintaining a simple, transparent API for developers. By using declarative sprite configuration, the engine makes intelligent pooling decisions without requiring manual pool management.

**Key Benefits:**
- 10-30% performance improvement for spawn-heavy games
- Reduced GC pauses and smoother frame times
- Zero developer overhead for common use cases
- Configurable and debuggable for fine-tuning

**Next Steps:**
1. Implement core PoolManager and Pool classes
2. Integrate with sprite lifecycle management
3. Add debug tools and logging
4. Profile with real game scenarios
5. Document patterns and best practices
