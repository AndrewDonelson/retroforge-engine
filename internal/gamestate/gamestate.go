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
)

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
	startTime        time.Time
	splashDuration   time.Duration
	autoTransitioned bool
}

// NewEngineSplashState creates a new engine splash state
func NewEngineSplashState(gsm *GameStateMachine) *EngineSplashState {
	return &EngineSplashState{
		gsm:            gsm,
		splashDuration: 2 * time.Second, // 2-3 seconds as per spec
	}
}

func (ess *EngineSplashState) Initialize(sm *statemachine.StateMachine) error {
	return nil
}

func (ess *EngineSplashState) Enter(sm *statemachine.StateMachine) {
	ess.startTime = time.Now()
	ess.autoTransitioned = false
}

func (ess *EngineSplashState) HandleInput(sm *statemachine.StateMachine) {
	// Any input can skip the splash - transition immediately
	// Check if any button is pressed
	hasInput := false
	for i := 0; i < 6; i++ {
		if input.Btnp(i) {
			hasInput = true
			break
		}
	}

	if hasInput {
		// Skip splash on any input
		if ess.gsm.initialState != "" {
			ess.gsm.StateMachine.ChangeState(ess.gsm.initialState)
		} else {
			// Default to "menu" if no initial state
			ess.gsm.StateMachine.ChangeState("menu")
		}
	}
}

func (ess *EngineSplashState) Update(dt float64) {
	if !ess.autoTransitioned && time.Since(ess.startTime) >= ess.splashDuration {
		// Auto-transition after duration
		if ess.gsm.initialState != "" {
			ess.gsm.StateMachine.ChangeState(ess.gsm.initialState)
		} else {
			// Default to "menu" if no initial state
			ess.gsm.StateMachine.ChangeState("menu")
		}
		ess.autoTransitioned = true
	}
}

func (ess *EngineSplashState) Draw() {
	// Draw engine splash screen
	if ess.gsm.renderer == nil {
		return
	}

	// Clear screen (dark blue/black)
	col := ess.gsm.palette.Color(1)
	ess.gsm.renderer.Clear(color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})

	// Draw engine name
	nameCol := ess.gsm.palette.Color(15) // White
	name := ess.gsm.engineName
	ess.gsm.renderer.Print(name, (ess.gsm.renderer.Width()-len(name)*font.Advance)/2,
		ess.gsm.renderer.Height()/2-20, color.RGBA{R: nameCol.R, G: nameCol.G, B: nameCol.B, A: nameCol.A})

	// Draw version
	version := "v" + ess.gsm.engineVersion
	versionCol := ess.gsm.palette.Color(7) // Light gray
	ess.gsm.renderer.Print(version, (ess.gsm.renderer.Width()-len(version)*font.Advance)/2,
		ess.gsm.renderer.Height()/2, color.RGBA{R: versionCol.R, G: versionCol.G, B: versionCol.B, A: versionCol.A})

	// Draw "Press any key" message at bottom
	msg := "Press any key..."
	msgCol := ess.gsm.palette.Color(6) // Gray
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
	return ess.autoTransitioned || time.Since(ess.startTime) >= ess.splashDuration
}

// CreditsState displays credits before exit
type CreditsState struct {
	gsm          *GameStateMachine
	scrollOffset float64
	scrollSpeed  float64
	hasShown     bool // Track if credits have been shown for at least one frame
	frameCount   int  // Track number of frames credits has been active
}

// NewCreditsState creates a new credits state
func NewCreditsState(gsm *GameStateMachine) *CreditsState {
	return &CreditsState{
		gsm:         gsm,
		scrollSpeed: 30.0, // pixels per second
	}
}

func (cs *CreditsState) Initialize(sm *statemachine.StateMachine) error {
	return nil
}

func (cs *CreditsState) Enter(sm *statemachine.StateMachine) {
	cs.scrollOffset = 0
	cs.hasShown = false // Reset on enter
	cs.frameCount = 0   // Reset frame counter
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
		sm.RequestExit()
		// Also request app quit (which sdlrun checks)
		app.RequestQuit()
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
}

func (cs *CreditsState) Draw() {
	// Draw scrolling credits
	if cs.gsm.renderer == nil {
		return
	}

	// Clear screen (dark blue/black)
	col := cs.gsm.palette.Color(1)
	cs.gsm.renderer.Clear(color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})

	// Draw title at top
	titleCol := cs.gsm.palette.Color(15) // White
	title := "CREDITS"
	cs.gsm.renderer.Print(title, (cs.gsm.renderer.Width()-len(title)*font.Advance)/2, 20,
		color.RGBA{R: titleCol.R, G: titleCol.G, B: titleCol.B, A: titleCol.A})

	// Draw engine credits
	engineName, engineVersion, engineDev := cs.gsm.GetEngineInfo()
	engineCol := cs.gsm.palette.Color(11) // Light blue
	y := 60

	engineLine := engineName + " " + engineVersion
	cs.gsm.renderer.Print(engineLine, (cs.gsm.renderer.Width()-len(engineLine)*font.Advance)/2, y,
		color.RGBA{R: engineCol.R, G: engineCol.G, B: engineCol.B, A: engineCol.A})
	y += 15

	devLine := "Developed by " + engineDev
	cs.gsm.renderer.Print(devLine, (cs.gsm.renderer.Width()-len(devLine)*font.Advance)/2, y,
		color.RGBA{R: engineCol.R, G: engineCol.G, B: engineCol.B, A: engineCol.A})
	y += 30

	// Draw game credits
	entryCol := cs.gsm.palette.Color(7)    // Light gray
	categoryCol := cs.gsm.palette.Color(6) // Gray

	// Group by category
	categories := make(map[string][]CreditEntry)
	for _, entry := range cs.gsm.creditsEntries {
		categories[entry.Category] = append(categories[entry.Category], entry)
	}

	// Draw credits by category
	for cat, entries := range categories {
		if y > cs.gsm.renderer.Height()-40 {
			break // Don't draw off screen
		}

		// Category header
		catText := cat + ":"
		cs.gsm.renderer.Print(catText, 20, y,
			color.RGBA{R: categoryCol.R, G: categoryCol.G, B: categoryCol.B, A: categoryCol.A})
		y += 15

		// Entries in this category
		for _, entry := range entries {
			if y > cs.gsm.renderer.Height()-40 {
				break
			}
			entryText := "  " + entry.Name
			if entry.Role != "" {
				entryText = entryText + " - " + entry.Role
			}
			cs.gsm.renderer.Print(entryText, 20, y,
				color.RGBA{R: entryCol.R, G: entryCol.G, B: entryCol.B, A: entryCol.A})
			y += 12
		}
		y += 10 // Space between categories
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
