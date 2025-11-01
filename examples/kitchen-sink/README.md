# Kitchen Sink Demo

A comprehensive demonstration of all RetroForge Engine features.

## Features Showcased

### 1. Module-Based State System
- Uses `rf.import()` to load state modules
- `menu_state.lua` and `play_state.lua` demonstrate module conventions
- Module-level variables persist across state transitions

### 2. Game State Machine
- Built-in engine splash (automatically shown)
- Custom menu state
- Play state with physics
- Built-in credits screen (automatically shown on exit)
- Uses `game.changeState()` and `game.exit()` for transitions

### 3. Automatic Sprite Pooling
- Ball sprite has `isUI=false` and `maxSpawn=100` (meets pooling criteria)
- Pools are created automatically - no developer code needed
- Sprites are reused for better performance

### 4. Physics Engine
- Box2D integration for realistic physics simulation
- Dynamic bodies (bouncing balls)
- Static bodies (boundary walls)
- Circle fixtures and collision detection

### 5. Stats Display
- Real-time FPS counter
- Memory usage (KB)
- Object count tracking
- Maximum objects reached
- Total objects spawned
- Timer showing play duration

## Game Flow

1. **Engine Splash** - Shows RetroForge branding (automatic)
2. **Menu** - Navigate with arrow keys, select with Z/X
3. **Play** - Physics demo with bouncing balls
   - Press any key to spawn more balls
   - Balls automatically spawn every 0.2 seconds
   - Runs for 30 seconds, then auto-transitions to credits
4. **Credits** - Shows feature credits (automatic on exit)
5. **Exit** - App closes when credits are skipped

## Running

```bash
# From retroforge-engine directory
make run-dev FOLDER=examples/kitchen-sink
```

## Code Structure

```
kitchen-sink/
├── manifest.json          # Cart metadata
├── assets/
│   ├── main.lua          # Entry point, imports modules
│   ├── menu_state.lua   # Menu module
│   ├── play_state.lua   # Physics demo module
│   └── sprites.json     # Ball sprite (pooled)
└── README.md            # This file
```

## Technical Details

- **Ball Sprite**: 12x12 pixels, `isUI=false`, `maxSpawn=100` → automatically pooled
- **Physics**: 6 pixel radius circles, dynamic bodies with random velocities
- **Boundaries**: 4 static walls creating a play area
- **Stats**: Updated every frame using `rf.stat()` function (dev mode only)

