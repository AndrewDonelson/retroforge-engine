package gamestate

import (
	"fmt"
	"time"

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
}

// NewGameStateMachine creates a new game state machine with built-in states
func NewGameStateMachine(isDebug bool, engineName, engineVersion, engineDeveloper string) *GameStateMachine {
	gsm := &GameStateMachine{
		StateMachine:    statemachine.NewStateMachine(),
		isDebug:         isDebug,
		creditsEntries:  make([]CreditEntry, 0),
		engineName:      engineName,
		engineVersion:   engineVersion,
		engineDeveloper: engineDeveloper,
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

// Start begins the state machine, showing engine splash if not in debug mode
func (gsm *GameStateMachine) Start(initialState string) error {
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
	// Any input can skip the splash
	// For now, we'll transition automatically after duration
}

func (ess *EngineSplashState) Update(dt float64) {
	if !ess.autoTransitioned && time.Since(ess.startTime) >= ess.splashDuration {
		// Auto-transition - but we need to know where to go
		// For now, we'll just mark as transitioned
		// The actual transition should be handled by the game
		ess.autoTransitioned = true
	}
}

func (ess *EngineSplashState) Draw() {
	// Draw engine logo/name
	// This will be handled by Lua bindings
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
}

func (cs *CreditsState) HandleInput(sm *statemachine.StateMachine) {
	// Any input exits credits and requests engine exit
	// The StateMachine interface has RequestExit() method
	sm.RequestExit()
}

func (cs *CreditsState) Update(dt float64) {
	cs.scrollOffset += cs.scrollSpeed * dt
}

func (cs *CreditsState) Draw() {
	// Draw scrolling credits
	// This will be handled by Lua bindings
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
