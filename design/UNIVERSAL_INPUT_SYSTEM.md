# Universal Input System Design

## Overview
A cross-platform 11-button input system that works consistently across desktop, mobile, and tablet devices. This system replaces the previous 6-button system to provide a cleaner, more consistent input experience across all platforms.

## 11 Universal Buttons

| Index | Name    | Purpose                          |
|-------|---------|----------------------------------|
| 0     | SELECT  | Menu navigation, secondary action |
| 1     | START   | Pause/menu, primary action        |
| 2     | UP      | Directional input                 |
| 3     | DOWN    | Directional input                 |
| 4     | LEFT    | Directional input                 |
| 5     | RIGHT   | Directional input                 |
| 6     | A       | Primary action button             |
| 7     | B       | Secondary action button           |
| 8     | X       | Tertiary action button            |
| 9     | Y       | Quaternary action button          |
| 10    | TURBO   | Modifier button (e.g., boost, run)|

## Key Mapping System

### Default Keyboard Mappings

| Button | Key Codes                                |
|--------|------------------------------------------|
| SELECT | `Enter`                                  |
| START  | `Space`                                  |
| UP     | `ArrowUp`                                |
| DOWN   | `ArrowDown`                              |
| LEFT   | `ArrowLeft`                              |
| RIGHT  | `ArrowRight`                             |
| A      | `KeyA`                                   |
| B      | `KeyS`                                   |
| X      | `KeyZ`                                   |
| Y      | `KeyX`                                   |
| TURBO  | `ShiftLeft`, `ShiftRight`               |

### Mobile/Tablet Support
- On-screen virtual controller (Controller component)
- Touch events mapped directly to buttons
- Portrait mode: Canvas at top, Controller at bottom
- Automatic detection via `isMobilePortrait()` function
- Controller shown/hidden automatically based on device orientation

## Architecture

### Engine Layer (`internal/input/input.go`)
- 11-button state array (`cur[11]`, `prev[11]`)
- `Set(i int, down bool)` - Set button by index (0-10)
- `SetByName(name string, down bool)` - Set button by name (e.g., "UP", "A", "TURBO")
- `Btn(i int) bool` - Check if button is pressed
- `Btnp(i int) bool` - Check if button was just pressed (edge-triggered)
- `Shift() bool` - Alias for `Btn(10)` (TURBO button, backward compatibility)

### WASM Interface Layer (`wasmInterface.ts`)
- `rf_set_button(name, down)` - Set button by name (new preferred method)
- `rf_set_btn(idx, down)` - Legacy index-based (backward compat)
- Key mapping handled in JavaScript layer with `defaultKeyMap`
- Supports both new button-by-name and legacy index-based systems

### Controller Component (`Controller/Controller.tsx`)
- React component with visual controller UI
- Sends button press/release events via `onButtonPress`/`onButtonRelease` callbacks
- Works on touch devices with `onPointerDown`/`onPointerUp` events
- Integrates with WASM via `useController()` hook and `rf_set_button()`
- Responsive design for mobile/tablet

### Platform Detection (`useController.ts`)
- `isMobilePortrait()` - Detects mobile/tablet in portrait mode
- Automatically shows/hides controller based on device type and orientation
- Updates on window resize events

### SDL Native Layer (`internal/sdlrun/sdlrun.go`)
- Maps SDL keyboard events to new 11-button system
- Supports same key mappings as WASM layer
- Direct mapping from SDL key codes to button indices

## Implementation Details

### Lua API (`internal/luabind/luabind.go`)
- `rf.btn(button)` - Supports all 11 buttons (0-10)
- `rf.btnp(button)` - Edge-triggered detection for all 11 buttons
- `rf.btnr(button)` - Button release detection (future)
- `rf.shift()` - Returns `rf.btn(10)` for backward compatibility

### Backward Compatibility

The old 6-button system is mapped to the new 11-button system:

| Old System    | New System          |
|---------------|---------------------|
| `BtnLeft (0)` | `BtnLEFT (4)`       |
| `BtnRight (1)`| `BtnRIGHT (5)`      |
| `BtnUp (2)`   | `BtnUP (2)`         |
| `BtnDown (3)` | `BtnDOWN (3)`       |
| `BtnO (4)`    | `BtnA (6)`          |
| `BtnX (5)`    | `BtnX (8)` or `BtnB (7)` |

Legacy constants are maintained for backward compatibility:
```go
BtnLeft  = BtnLEFT  // Deprecated: use BtnLEFT
BtnRight = BtnRIGHT // Deprecated: use BtnRIGHT
BtnUp    = BtnUP    // Deprecated: use BtnUP
BtnDown  = BtnDOWN  // Deprecated: use BtnDOWN
BtnO     = BtnA     // Deprecated: use BtnA
```

Lua API maintains `rf.btn()` and `rf.btnp()` for compatibility, but now supports all 11 buttons.

## Benefits

1. **Consistent Experience**: Same button layout across all platforms
2. **Mobile Support**: Native touch controller integration
3. **No Conflicts**: Eliminates ESC/DOWN button conflicts
4. **Clear Semantics**: Button names (START, SELECT, TURBO) are self-explanatory
5. **Future-Proof**: Easy to extend or remap keys if needed

## Example Usage

```lua
function _HANDLE_INPUT()
  -- Pause with START button
  if rf.btnp(1) then
    game.pushState("pause")
  end
  
  -- Movement
  if rf.btn(4) then  -- LEFT
    player.vx = -3
  elseif rf.btn(5) then  -- RIGHT
    player.vx = 3
  end
  
  -- Jump
  if rf.btnp(6) then  -- A button
    player.vy = -10
  end
  
  -- Boost with TURBO
  if rf.btn(10) then  -- TURBO
    speed = speed * 1.5
  end
end
```

