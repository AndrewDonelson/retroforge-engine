# Input Flow: Desktop vs WASM

## DESKTOP (SDL) - WORKING ✓

### Flow:
```
1. SDL Event Loop (sdlrun.go:52-103)
   └─> input.Step()              // STEP 1: Copy cur → prev (clear edges for this frame)
   └─> sdl.PollEvent()
       └─> KeyboardEvent
           └─> input.Set(buttonIndex, down)  // STEP 2: Set cur[buttonIndex] = down

2. Engine Frame (engine.go:66-128)
   └─> e.RunFrames(1)
       └─> "tick" event fires
           └─> HandleInput()     // STEP 3: Check input (prev from last frame, cur from this frame)
               └─> GSM.HandleInput()
                   └─> Lua _HANDLE_INPUT()
                       └─> rf.btnp(index)
                           └─> input.Btnp(index)  // Returns cur[i] && !prev[i]
                               └─> prev[i] = false (from last frame)
                               └─> cur[i] = true (just set)
                               └─> Returns TRUE ✓
           └─> Update()
           └─> Draw()
           └─> input.Step()     // STEP 4: Copy cur → prev for NEXT frame
```

### Key Points:
- `input.Step()` is called **BEFORE** processing events
- `prev` contains state from **previous frame**
- `cur` is updated **during** event processing
- When `HandleInput()` runs, `prev` is still old, `cur` is new → `btnp()` works ✓

---

## WASM (Browser) - FIXED ✓

### Flow (After Fix):
```
1. Browser Event (Controller.tsx or Keyboard)
   └─> handleButtonPress/Release
       └─> useController.ts
           └─> window.rf_set_btn(idx, true/false)  // Called asynchronously

2. WASM Export (cmd/wasm/main.go:87-105)
   └─> rfSetBtn()
       └─> input.Set(idx, down)  // Updates cur[idx] = down
           └─> cur[idx] = down

3. Engine Frame (engine.go:66-128) - FIXED ✓
   └─> e.RunFrames(1)
       └─> "tick" event fires
           └─> input.Step()      // STEP 1: Copy cur → prev (same as SDL!)
               └─> prev = cur    // Save current frame state for edge detection
           └─> HandleInput()     // STEP 2: Check input (prev = old, cur = new)
               └─> GSM.HandleInput()
                   └─> Lua _HANDLE_INPUT()
                       └─> rf.btnp(index)
                           └─> input.Btnp(index)  // Returns cur[i] && !prev[i]
                               └─> prev[i] = false (from last frame, now copied)
                               └─> cur[i] = true (just set by rf_set_btn)
                               └─> Returns TRUE ✓
           └─> Update()
           └─> Draw()
```

### Key Points (After Fix):
- `input.Step()` is called **BEFORE** `HandleInput()` (matches SDL)
- `prev` contains state from **previous frame** (copied at start)
- `cur` is updated **asynchronously** by `rf_set_btn()`
- When `HandleInput()` runs, `prev` is old, `cur` is new → `btnp()` works ✓

---

## THE FIX APPLIED:

✅ **Moved `input.Step()` to START of frame (before HandleInput)**
- Now matches SDL behavior exactly: `Step() → Set buttons → HandleInput()`
- Ensures `prev` and `cur` are in correct state for edge detection
- `btnp()` now works correctly in WASM! ✓

