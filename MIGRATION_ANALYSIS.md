# Migration Analysis: SDL2 → Ebiten & MIT → Apache 2.0

## Question 1: Switching from SDL2 to Ebiten

### Complexity Assessment: **MEDIUM** (2-3 days of focused work)

### Current SDL2 Usage

SDL2 is used in **2 isolated modules** with build tags:

1. **`internal/sdlrun/sdlrun.go`** (`//go:build !js && !wasm`)
   - Window creation and management
   - Event loop (keyboard input)
   - Texture rendering (pixel buffer → SDL texture → screen)
   - Screenshot functionality
   - ~147 lines of code

2. **`internal/audio/audio.go`** (`//go:build !js`)
   - Audio device initialization
   - Audio mixing loop
   - SDL audio playback
   - ~206 lines of code

**Key Points:**
- WASM already doesn't use SDL (has separate implementation)
- SDL usage is cleanly isolated with build tags
- Your engine abstraction layer (`internal/engine`) doesn't depend on SDL directly

### Migration Steps

#### Phase 1: Window & Rendering (Ebiten) - ~6 hours

**Replace `internal/sdlrun/sdlrun.go` with `internal/ebitenrun/ebitenrun.go`:**

```go
//go:build !js && !wasm

package ebitenrun

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image"
)

type Game struct {
    engine *engine.Engine
    scale  int
}

func (g *Game) Update() error {
    // Input handling (Ebiten's input API)
    // Engine frame update
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Draw engine's pixel buffer to screen
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 480 * g.scale, 270 * g.scale
}

func RunWindow(e *engine.Engine, scale int) error {
    game := &Game{engine: e, scale: scale}
    ebiten.SetWindowSize(480*scale, 270*scale)
    ebiten.SetWindowTitle("RetroForge")
    return ebiten.RunGame(game)
}
```

**Key Changes:**
- `sdl.CreateWindow()` → `ebiten.SetWindowSize()` / `ebiten.SetWindowTitle()`
- `sdl.PollEvent()` → `ebitenutil.IsKeyPressed()` / `ebitenutil.KeyPressDuration()`
- `sdl.CreateRenderer()` / `tex.Update()` / `ren.Copy()` → `ebiten.Image` operations
- Event loop → Ebiten's `Update()` / `Draw()` callbacks

**Effort:** Medium - API is different but conceptually similar

#### Phase 2: Audio (Pure Go) - ~4 hours

**Replace `internal/audio/audio.go` with pure Go implementation:**

**Option A: Use `oto` (Pure Go audio)**
```go
import "github.com/hajimehoshi/oto/v3"

// Simple audio driver - no CGO needed
```

**Option B: Use Ebiten's built-in audio (recommended)**
```go
import (
    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/wav"
)
```

**Effort:** Medium - Need to rewrite mixer but logic stays the same

#### Phase 3: Input Mapping - ~2 hours

**Update keyboard mappings:**
- SDL key constants → Ebiten key constants
- `sdl.K_UP` → `ebiten.KeyArrowUp`
- `sdl.K_ESCAPE` → `ebiten.KeyEscape`
- etc.

**Effort:** Low - Straightforward mapping

#### Phase 4: Testing & Cleanup - ~4 hours

- Remove `github.com/veandco/go-sdl2` dependency
- Update build scripts (remove SDL2 deps)
- Test all platforms
- Update documentation

**Effort:** Low-Medium

### Estimated Total Time: **16-20 hours** (2-3 days)

### Benefits

✅ **Zero CGO dependencies** - Pure Go builds
✅ **Simple cross-compilation** - `GOOS=windows go build` just works
✅ **Better WASM support** - Ebiten has excellent WASM integration
✅ **Smaller binaries** - No SDL2 library dependencies
✅ **Better performance** - Ebiten is highly optimized
✅ **Active maintenance** - Ebiten is actively developed

### Challenges

⚠️ **API differences** - Need to learn Ebiten's API patterns
⚠️ **Audio rewrite** - Need new audio implementation (but simpler overall)
⚠️ **Testing required** - Need to verify all input/graphics work correctly

---

## Question 2: Changing License from MIT to Apache 2.0

### Complexity Assessment: **LOW** (2-4 hours)

### Files to Update

#### Engine (`retroforge-engine/`):

1. **`LICENSE`** - Replace MIT with Apache 2.0 text
2. **`README.md`** - Update license mention
3. **`go.mod`** - No change needed (license not in go.mod)

#### Webapp (`retroforge-webapp/`):

1. **`LICENSE`** - Replace MIT with Apache 2.0 text
2. **`public/json/license.json`** - Update license text and metadata
3. **`src/components/common/License.tsx`** - Update if hardcoded text
4. **`package.json`** - Update `license` field (if present)
5. **Any README files** - Update license mentions

#### Other Considerations:

- **Git history** - Old commits will still have MIT references (this is OK)
- **Dependencies** - No impact (you're not changing dependency licenses)
- **Contributors** - No action needed (Apache 2.0 is permissive like MIT)
- **Legal notices** - May want to add Apache 2.0 NOTICE file (optional)

### Apache 2.0 Template

You'll need to replace the LICENSE files with Apache 2.0 text. Key differences from MIT:

1. **Patent grant** - Apache 2.0 includes explicit patent protection
2. **Attribution** - Must include Apache 2.0 license text in distributions
3. **Notice file** - Optional but recommended for proper attribution
4. **More verbose** - Apache 2.0 is longer but more legally precise

### Estimated Total Time: **2-4 hours**

### Benefits

✅ **Better patent protection** - Explicit patent grant clause
✅ **More legally precise** - Better protection for you and users
✅ **Industry standard** - Many major projects use Apache 2.0
✅ **Still permissive** - Compatible with all use cases (commercial, GPL, etc.)

### No Downsides

- Same permissiveness as MIT
- Compatible with all existing MIT dependencies
- No breaking changes for users

---

## Recommended Approach

### Option 1: Migrate Both (Recommended)
**Timeline:** 3-4 days total
- Days 1-2: SDL2 → Ebiten migration
- Day 3: License update + testing
- Day 4: Buffer for issues

**Benefits:**
- Eliminate all CGO issues permanently
- Improve cross-compilation experience
- Better license protection
- Future-proof the project

### Option 2: License First, Ebiten Later
**Timeline:** License in 1 day, Ebiten when ready

**If you're not ready for Ebiten migration yet:**
- Update license now (low risk, high value)
- Plan Ebiten migration for next sprint
- Keep SDL2 working in the meantime

---

## Next Steps

1. **Decision:** Choose Option 1 or 2
2. **I can help with:**
   - Writing the Ebiten replacement code
   - Updating all license files
   - Testing the migration
   - Updating build scripts

Would you like me to start with either migration?

