package gamestate

import (
	"fmt"
	"image/color"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/app"
	"github.com/AndrewDonelson/retroforge-engine/internal/font"
	"github.com/AndrewDonelson/retroforge-engine/internal/graphics"
	"github.com/AndrewDonelson/retroforge-engine/internal/input"
	"github.com/AndrewDonelson/retroforge-engine/internal/pal"
	"github.com/AndrewDonelson/retroforge-engine/internal/statemachine"
)

const (
	EngineSplashStateName = "__engine_splash"
	CreditsStateName      = "__credits"
	// Logo sprite dimensions (28x28 from PNG, no scaling)
	logoWidth  = 28
	logoHeight = 28
)

// logoPixels contains the hardcoded RetroForge logo sprite data (28x28)
// Direct conversion from logo.png - no scaling, closest palette color matching
// Colors mapped to RetroForge 50 palette indices:
// -1=transparent, 0=black, 1=white, 16=dark blue (#0f172a), 19=cyan shadow (#009c8b),
// 22=sky blue shadow, 25=blue highlight, 28=indigo (#6f88ff), 43=brown shadow,
// 48=cyan blue (#38bdf8), 49=dark cyan blue (#0081bc)
var logoPixels = [][]int{
	{ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{ -1, -1, -1, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, -1, -1, -1},
	{ -1, -1, 49, 28, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 28, 49, -1, -1},
	{ -1, 16, 28,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 28, 16, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0, 43,  0, 43,  0, 43,  0, 43,  0, 43,  0, 43,  0, 43,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0, 25, 48, 48, 48, 48, 48, 25, 28,  0, 43, 43, 48, 48, 48, 48, 48, 48, 48, 48, 28,  0, 16, 19, -1},
	{ -1, 19, 16,  0, 43, 25, 48, 25, 28, 25, 48, 48, 43,  0,  0, 28, 48, 22, 28, 28, 28, 25, 48, 28,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0, 28, 48, 43,  0,  0, 43, 48, 28,  0,  0,  0, 22, 25,  0, 25, 28, 43, 22, 43,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0, 28, 48, 43,  0, 43, 25, 48, 16,  0,  0,  0, 22, 22, 28, 22, 25,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0, 28, 48, 48, 48, 48, 48, 28,  0,  0,  0,  0, 22, 48, 48, 48, 49,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0, 28, 48, 25, 28, 22, 48, 43,  0,  0,  0,  0, 22, 49,  0, 49, 28,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0, 19, 48, 43,  0, 43, 22, 22, 43,  0,  0,  0, 22, 49,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0, 43, 25, 48, 25, 43,  0, 43, 22, 22, 16,  0, 28, 22, 22, 28, 28,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0, 25, 48, 48, 48, 49,  0,  0, 28, 48, 25,  0, 22, 22, 22, 22, 22, 16,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 19, 16,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 16, 19, -1},
	{ -1, 16, 28,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 28, 16, -1},
	{ -1, -1, 49, 28, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 28, 49, -1, -1},
	{ -1, -1, -1, 16, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 19, 16, -1, -1, -1},
	{ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
}
// CreditEntry represents a single credit entry
type CreditEntry struct {
	Category string
	Name     string
	Role     string
}

// GameStateMachine extends StateMachine with built-in engine splash and credits states
type GameStateMachine struct {
	*statemachine.StateMachine

	isDebug        bool
	engineSplash   *EngineSplashState
	credits        *CreditsState
	creditsEntries []CreditEntry

	// Engine info for splash
	engineName      string
	engineVersion   string
	engineDeveloper string

	// Initial state to transition to after splash
	initialState string

	// Renderer and palette for drawing built-in states
	renderer graphics.Renderer
	palette  *pal.Manager
}

// NewGameStateMachine creates a new game state machine with built-in states
func NewGameStateMachine(isDebug bool, engineName, engineVersion, engineDeveloper string, renderer graphics.Renderer, palette *pal.Manager) *GameStateMachine {
	gsm := &GameStateMachine{
		StateMachine:    statemachine.NewStateMachine(),
		isDebug:         isDebug,
		creditsEntries:  make([]CreditEntry, 0),
		engineName:      engineName,
		engineVersion:   engineVersion,
		engineDeveloper: engineDeveloper,
		renderer:        renderer,
		palette:         palette,
	}

	// Create built-in states
	gsm.engineSplash = NewEngineSplashState(gsm)
	gsm.credits = NewCreditsState(gsm)

	// Register built-in states
	gsm.StateMachine.RegisterStateInstance(EngineSplashStateName, gsm.engineSplash)
	gsm.StateMachine.RegisterStateInstance(CreditsStateName, gsm.credits)

	return gsm
}

// AddCreditEntry adds a credit entry to be displayed in the credits state
func (gsm *GameStateMachine) AddCreditEntry(category, name, role string) {
	gsm.creditsEntries = append(gsm.creditsEntries, CreditEntry{
		Category: category,
		Name:     name,
		Role:     role,
	})
}

// GetCreditsEntries returns all credit entries
func (gsm *GameStateMachine) GetCreditsEntries() []CreditEntry {
	return gsm.creditsEntries
}

// IsDebug returns whether this is a debug build
func (gsm *GameStateMachine) IsDebug() bool {
	return gsm.isDebug
}

// GetEngineInfo returns engine information
func (gsm *GameStateMachine) GetEngineInfo() (name, version, developer string) {
	return gsm.engineName, gsm.engineVersion, gsm.engineDeveloper
}

// SetRenderer sets the renderer for drawing built-in states
func (gsm *GameStateMachine) SetRenderer(renderer graphics.Renderer) {
	gsm.renderer = renderer
}

// SetPalette sets the palette for drawing built-in states
func (gsm *GameStateMachine) SetPalette(palette *pal.Manager) {
	gsm.palette = palette
}

// Start begins the state machine, showing engine splash if not in debug mode
func (gsm *GameStateMachine) Start(initialState string) error {
	gsm.initialState = initialState // Store for splash transition
	if !gsm.isDebug {
		// Show engine splash first
		return gsm.StateMachine.ChangeState(EngineSplashStateName)
	}
	// In debug mode, skip splash and go directly to initial state
	if initialState != "" {
		return gsm.StateMachine.ChangeState(initialState)
	}
	return nil
}

// Exit transitions to credits state before exiting
func (gsm *GameStateMachine) Exit() error {
	// Transition to credits
	return gsm.StateMachine.ChangeState(CreditsStateName)
}

// Override ChangeState to prevent direct changes to built-in states from outside
// (except through Start() and Exit())
func (gsm *GameStateMachine) ChangeState(name string) error {
	if name == EngineSplashStateName || name == CreditsStateName {
		return fmt.Errorf("cannot directly change to built-in state '%s' (use Start() or Exit())", name)
	}
	return gsm.StateMachine.ChangeState(name)
}

// EngineSplashState displays engine branding
type EngineSplashState struct {
	gsm              *GameStateMachine
	elapsedTime      float64 // Accumulated time in seconds
	frameCount       int     // Frame counter as fallback (120 frames = ~2 seconds at 60 FPS)
	splashDuration   float64 // Duration in seconds (2.0 seconds)
	splashFrames     int     // Frame count fallback (120 frames at 60 FPS = 2 seconds)
	autoTransitioned bool
}

// NewEngineSplashState creates a new engine splash state
func NewEngineSplashState(gsm *GameStateMachine) *EngineSplashState {
	return &EngineSplashState{
		gsm:            gsm,
		splashDuration: 2.0, // 2 seconds as per spec
		splashFrames:  120,  // Fallback: 120 frames at 60 FPS = 2 seconds
	}
}

func (ess *EngineSplashState) Initialize(sm *statemachine.StateMachine) error {
	return nil
}

func (ess *EngineSplashState) Enter(sm *statemachine.StateMachine) {
	ess.elapsedTime = 0.0
	ess.frameCount = 0
	ess.autoTransitioned = false
}

// getTargetState determines the next state after engine splash
// Priority: 1) user-defined splash_state, 2) initial state from Start(), 3) menu
func (ess *EngineSplashState) getTargetState(sm *statemachine.StateMachine) string {
	// First, check if game has a splash_state module
	if sm.IsStateRegistered("splash") {
		return "splash"
	}
	// Second, use initial state if provided
	if ess.gsm.initialState != "" {
		return ess.gsm.initialState
	}
	// Default to menu
	return "menu"
}

func (ess *EngineSplashState) HandleInput(sm *statemachine.StateMachine) {
	// Debug logging removed to reduce console spam
	
	// Any input can skip the splash - transition immediately
	// Use Btn() instead of Btnp() to detect any currently pressed button
	// (Btnp requires edge detection which is timing-sensitive in WASM)
	hasInput := false
	for i := 0; i < 6; i++ {
		if input.Btn(i) {
			hasInput = true
			break
		}
	}

	if hasInput {
		// Skip splash on any input
		targetState := ess.getTargetState(sm)

		// Check if target state exists before trying to transition
		if sm.IsStateRegistered(targetState) {
			// State exists - transition
			ess.gsm.StateMachine.ChangeState(targetState)
		} else {
			// State doesn't exist - this game doesn't use state machine (old-style Lua game)
			// Pop splash screen to allow fallback to direct Lua calls
			ess.gsm.StateMachine.PopState()
		}
		// Mark as transitioned to prevent auto-transition from also firing
		ess.autoTransitioned = true
	}
}

func (ess *EngineSplashState) Update(dt float64) {
	// Increment frame counter as fallback (always works regardless of dt)
	ess.frameCount++
	
	// Accumulate elapsed time using frame delta (works in WASM where time.Since may be unreliable)
	// Ensure dt is valid (at least a small positive value to prevent accumulation issues)
	if dt <= 0 {
		dt = 0.016 // Default to ~60 FPS if dt is invalid
	}
	// Cap dt to prevent huge jumps (e.g., if tab was inactive)
	if dt > 0.1 {
		dt = 0.1 // Cap at 100ms max per frame
	}
	ess.elapsedTime += dt
	
	// Auto-transition if either condition is met:
	// 1. Time-based: elapsedTime >= splashDuration (2.0 seconds)
	// 2. Frame-based: frameCount >= splashFrames (120 frames at 60 FPS)
	// The frame counter ensures auto-advance works even if dt accumulation fails
	shouldTransition := !ess.autoTransitioned && (ess.elapsedTime >= ess.splashDuration || ess.frameCount >= ess.splashFrames)
	if shouldTransition {
		// Auto-transition after duration
		targetState := ess.getTargetState(ess.gsm.StateMachine)

		// Check if target state exists before trying to transition
		if ess.gsm.StateMachine.IsStateRegistered(targetState) {
			// State exists - transition now
			if err := ess.gsm.StateMachine.ChangeState(targetState); err != nil {
				// Transition failed - pop splash to allow fallback
				ess.gsm.StateMachine.PopState()
			}
			ess.autoTransitioned = true
		} else {
			// State doesn't exist - this game doesn't use state machine (old-style Lua game)
			// Pop splash screen to allow fallback to direct Lua calls
			// This allows games like moon-lander that use direct _UPDATE/_DRAW to work
			ess.gsm.StateMachine.PopState()
			ess.autoTransitioned = true
		}
	}
}

func (ess *EngineSplashState) Draw() {
	// Draw engine splash screen
	if ess.gsm.renderer == nil {
		return
	}
	if ess.gsm.palette == nil {
		return
	}

	// Reset camera and clipping to ensure clean drawing
	ess.gsm.renderer.SetCamera(0, 0)
	ess.gsm.renderer.SetClip(0, 0, 0, 0) // Disable clipping

	// Clear screen with black (index 0 is black)
	col := ess.gsm.palette.Color(0)
	ess.gsm.renderer.Clear(color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})

	// Calculate logo position (centered on screen)
	logoX := (ess.gsm.renderer.Width() - logoWidth) / 2
	logoY := (ess.gsm.renderer.Height() - logoHeight) / 2

	// Draw logo sprite pixel by pixel
	for y := 0; y < logoHeight; y++ {
		for x := 0; x < logoWidth; x++ {
			pixelColorIndex := logoPixels[y][x]
			// Skip transparent pixels
			if pixelColorIndex == -1 {
				continue
			}
			// Get color from palette
			pixelCol := ess.gsm.palette.Color(pixelColorIndex)
			// Draw pixel
			ess.gsm.renderer.PSet(logoX+x, logoY+y, color.RGBA{R: pixelCol.R, G: pixelCol.G, B: pixelCol.B, A: pixelCol.A})
		}
	}

	// Draw "Press any key" message at bottom
	msg := "Press any key..."
	msgCol := ess.gsm.palette.Color(1) // White (index 1)
	ess.gsm.renderer.Print(msg, (ess.gsm.renderer.Width()-len(msg)*font.Advance)/2,
		ess.gsm.renderer.Height()-20, color.RGBA{R: msgCol.R, G: msgCol.G, B: msgCol.B, A: msgCol.A})
}

func (ess *EngineSplashState) Exit(sm *statemachine.StateMachine) {
	// Cleanup if needed
}

func (ess *EngineSplashState) Shutdown() {
	// Final cleanup
}

// ShouldTransition checks if splash should transition
func (ess *EngineSplashState) ShouldTransition() bool {
	return ess.autoTransitioned || ess.elapsedTime >= ess.splashDuration
}

// CreditsState displays credits before exit
type CreditsState struct {
	gsm              *GameStateMachine
	scrollOffset     float64
	scrollSpeed      float64
	hasShown         bool          // Track if credits have been shown for at least one frame
	frameCount       int           // Track number of frames credits has been active
	startTime        time.Time     // Track when credits started
	creditsDuration  time.Duration // Duration before auto-exit
	autoTransitioned bool          // Track if auto-transition happened
}

// NewCreditsState creates a new credits state
func NewCreditsState(gsm *GameStateMachine) *CreditsState {
	return &CreditsState{
		gsm:             gsm,
		scrollSpeed:     30.0,            // pixels per second
		creditsDuration: 10 * time.Second, // 10 seconds auto-continue
	}
}

func (cs *CreditsState) Initialize(sm *statemachine.StateMachine) error {
	return nil
}

func (cs *CreditsState) Enter(sm *statemachine.StateMachine) {
	cs.scrollOffset = 0
	cs.hasShown = false      // Reset on enter
	cs.frameCount = 0        // Reset frame counter
	cs.startTime = time.Now() // Record start time
	cs.autoTransitioned = false
	// Reset input state to prevent immediate exit from button that triggered exit
	// (HandleInput will be called before Update, so we need hasShown protection)
}

func (cs *CreditsState) HandleInput(sm *statemachine.StateMachine) {
	// CRITICAL: Only exit if credits have been shown for at least 2 frames
	// This prevents immediate exit when credits state first enters
	// HandleInput is called BEFORE Update, so frameCount will be 0 on first frame
	// We need to wait 2 frames to ensure the button press that triggered exit is cleared
	if cs.frameCount < 2 {
		return // Ignore all input until credits have been displayed for at least 2 frames
	}

	// Only exit if there's actual user input (any button pressed)
	// Check all buttons - if any are currently pressed, exit
	hasInput := false
	for i := 0; i < 6; i++ {
		if input.Btnp(i) {
			hasInput = true
			break
		}
	}

	if hasInput {
		// Any input exits credits and requests engine exit
		cs.exitCredits(sm)
	}
}

func (cs *CreditsState) Update(dt float64) {
	cs.scrollOffset += cs.scrollSpeed * dt
	// Increment frame counter
	cs.frameCount++
	// Mark that credits have been shown for at least one frame
	if !cs.hasShown && cs.frameCount >= 1 {
		cs.hasShown = true
	}

	// Auto-exit after 10 seconds
	if !cs.autoTransitioned && time.Since(cs.startTime) >= cs.creditsDuration {
		// Auto-exit after duration (same as user pressing a key)
		cs.autoTransitioned = true
		// Get state machine to request exit
		if cs.gsm != nil && cs.gsm.StateMachine != nil {
			cs.exitCredits(cs.gsm.StateMachine)
		}
	}
}

// exitCredits handles the actual exit logic (shared between input and auto-exit)
func (cs *CreditsState) exitCredits(sm *statemachine.StateMachine) {
	// Request engine exit
	sm.RequestExit()
	// Also request app quit (which sdlrun checks)
	app.RequestQuit()
}

func (cs *CreditsState) Draw() {
	// Draw scrolling credits
	if cs.gsm.renderer == nil {
		return
	}
	if cs.gsm.palette == nil {
		return
	}

	// Reset camera and clipping to ensure clean drawing
	cs.gsm.renderer.SetCamera(0, 0)
	cs.gsm.renderer.SetClip(0, 0, 0, 0) // Disable clipping

	// Clear screen with black (index 0 is black)
	col := cs.gsm.palette.Color(0)
	cs.gsm.renderer.Clear(color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})

	// Draw title at top
	titleCol := cs.gsm.palette.Color(15) // White
	title := "CREDITS"
	cs.gsm.renderer.Print(title, (cs.gsm.renderer.Width()-len(title)*font.Advance)/2, 20,
		color.RGBA{R: titleCol.R, G: titleCol.G, B: titleCol.B, A: titleCol.A})

	// Draw engine credits (more compact)
	engineName, engineVersion, engineDev := cs.gsm.GetEngineInfo()
	engineCol := cs.gsm.palette.Color(11) // Light blue
	y := 50

	engineLine := engineName + " " + engineVersion
	cs.gsm.renderer.Print(engineLine, (cs.gsm.renderer.Width()-len(engineLine)*font.Advance)/2, y,
		color.RGBA{R: engineCol.R, G: engineCol.G, B: engineCol.B, A: engineCol.A})
	y += 12

	devLine := "Developed by " + engineDev
	cs.gsm.renderer.Print(devLine, (cs.gsm.renderer.Width()-len(devLine)*font.Advance)/2, y,
		color.RGBA{R: engineCol.R, G: engineCol.G, B: engineCol.B, A: engineCol.A})
	y += 20

	// Draw game credits (more compact, limit displayed)
	entryCol := cs.gsm.palette.Color(7)    // Light gray
	categoryCol := cs.gsm.palette.Color(6) // Gray

	// Group by category
	categories := make(map[string][]CreditEntry)
	for _, entry := range cs.gsm.creditsEntries {
		categories[entry.Category] = append(categories[entry.Category], entry)
	}

	// Sort category names for deterministic ordering (prevents flickering)
	catNames := make([]string, 0, len(categories))
	for cat := range categories {
		catNames = append(catNames, cat)
	}
	// Simple alphabetical sort
	for i := 0; i < len(catNames)-1; i++ {
		for j := i + 1; j < len(catNames); j++ {
			if catNames[i] > catNames[j] {
				catNames[i], catNames[j] = catNames[j], catNames[i]
			}
		}
	}

	// Draw credits by category (more compact layout)
	bottomMargin := 30 // Reserve space for "Press any key" message
	for _, cat := range catNames {
		entries := categories[cat]
		if y > cs.gsm.renderer.Height()-bottomMargin {
			break // Don't draw off screen
		}

		// Category header (compact)
		catText := cat + ":"
		cs.gsm.renderer.Print(catText, 20, y,
			color.RGBA{R: categoryCol.R, G: categoryCol.G, B: categoryCol.B, A: categoryCol.A})
		y += 10

		// Entries in this category (limit per category to prevent overflow)
		maxEntriesPerCat := 5
		for i, entry := range entries {
			if i >= maxEntriesPerCat || y > cs.gsm.renderer.Height()-bottomMargin {
				break
			}
			// More compact: just name, skip role if too long
			entryText := "  " + entry.Name
			if entry.Role != "" && len(entryText)+len(entry.Role) < 40 {
				entryText = entryText + " - " + entry.Role
			}
			cs.gsm.renderer.Print(entryText, 20, y,
				color.RGBA{R: entryCol.R, G: entryCol.G, B: entryCol.B, A: entryCol.A})
			y += 10 // Reduced spacing
		}
		y += 5 // Space between categories
	}

	// Draw "Press any key to exit" at bottom
	msgCol := cs.gsm.palette.Color(6) // Gray
	msg := "Press any key to exit"
	cs.gsm.renderer.Print(msg, (cs.gsm.renderer.Width()-len(msg)*font.Advance)/2,
		cs.gsm.renderer.Height()-15, color.RGBA{R: msgCol.R, G: msgCol.G, B: msgCol.B, A: msgCol.A})
}

func (cs *CreditsState) Exit(sm *statemachine.StateMachine) {
	// Cleanup
}

func (cs *CreditsState) Shutdown() {
	// Final cleanup
}

// GetScrollOffset returns the current scroll offset
func (cs *CreditsState) GetScrollOffset() float64 {
	return cs.scrollOffset
}

// HandleCreditsInput is called when input is detected in credits state
func (cs *CreditsState) HandleCreditsInput() {
	// This will be called from the game loop when input is detected
	// and will request exit
}
