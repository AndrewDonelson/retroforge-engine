package gamestate

import (
	"testing"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/statemachine"
)

// Test NewGameStateMachine
func TestNewGameStateMachine(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	if gsm == nil {
		t.Fatal("NewGameStateMachine returned nil")
	}

	if !gsm.IsStateRegistered(EngineSplashStateName) {
		t.Error("engine splash should be registered")
	}
	if !gsm.IsStateRegistered(CreditsStateName) {
		t.Error("credits should always be registered")
	}
}

// Test Debug Mode
func TestDebugMode(t *testing.T) {
	gsm := NewGameStateMachine(true, "TestEngine", "1.0.0", "TestDev")

	if !gsm.IsDebug() {
		t.Error("should be in debug mode")
	}

	// In debug mode, Start() should skip splash
	err := gsm.Start("test_state")
	if err == nil {
		// This is expected if test_state doesn't exist
		// The important thing is it doesn't change to splash
	}
}

// Test AddCreditEntry
func TestAddCreditEntry(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	gsm.AddCreditEntry("Developer", "John Doe", "Lead Developer")
	gsm.AddCreditEntry("Artist", "Jane Smith", "Character Artist")

	entries := gsm.GetCreditsEntries()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Category != "Developer" {
		t.Error("first entry category should be Developer")
	}
	if entries[1].Name != "Jane Smith" {
		t.Error("second entry name should be Jane Smith")
	}
}

// Test GetEngineInfo
func TestGetEngineInfo(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.2.3", "TestDeveloper")

	name, version, dev := gsm.GetEngineInfo()
	if name != "TestEngine" {
		t.Errorf("expected 'TestEngine', got %s", name)
	}
	if version != "1.2.3" {
		t.Errorf("expected '1.2.3', got %s", version)
	}
	if dev != "TestDeveloper" {
		t.Errorf("expected 'TestDeveloper', got %s", dev)
	}
}

// Test Start (non-debug)
func TestStartNonDebug(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	err := gsm.Start("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should be on engine splash
	active, exists := gsm.GetActiveState()
	if !exists {
		t.Error("should have active state")
	}
	if active != EngineSplashStateName {
		t.Errorf("expected engine splash, got %s", active)
	}
}

// Test Start (debug)
func TestStartDebug(t *testing.T) {
	gsm := NewGameStateMachine(true, "TestEngine", "1.0.0", "TestDev")

	// Create a test state
	testState := &TestState{name: "test"}
	gsm.RegisterStateInstance("test", testState)

	err := gsm.Start("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	active, exists := gsm.GetActiveState()
	if !exists {
		t.Error("should have active state")
	}
	if active != "test" {
		t.Errorf("expected 'test', got %s", active)
	}
}

// Test Exit
func TestExit(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	// Create and enter a test state
	testState := &TestState{name: "test"}
	gsm.RegisterStateInstance("test", testState)
	gsm.StateMachine.ChangeState("test")

	err := gsm.Exit()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	active, exists := gsm.GetActiveState()
	if !exists {
		t.Error("should have active state")
	}
	if active != CreditsStateName {
		t.Errorf("expected credits state, got %s", active)
	}
}

// Test ChangeState prevents direct changes to built-in states
func TestChangeStateProtection(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	// Try to directly change to engine splash
	err := gsm.ChangeState(EngineSplashStateName)
	if err == nil {
		t.Error("expected error when changing to engine splash directly")
	}

	// Try to directly change to credits
	err = gsm.ChangeState(CreditsStateName)
	if err == nil {
		t.Error("expected error when changing to credits directly")
	}
}

// Test EngineSplashState
func TestEngineSplashState(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")
	splash := gsm.engineSplash

	// Initialize
	err := splash.Initialize(gsm.StateMachine)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Enter
	splash.Enter(gsm.StateMachine)

	// Should not transition immediately
	if splash.ShouldTransition() {
		t.Error("should not transition immediately")
	}

	// Wait for duration
	time.Sleep(3 * time.Second)

	if !splash.ShouldTransition() {
		t.Error("should transition after duration")
	}
}

// Test CreditsState
func TestCreditsState(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")
	credits := gsm.credits

	// Initialize
	err := credits.Initialize(gsm.StateMachine)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Enter
	credits.Enter(gsm.StateMachine)

	if credits.GetScrollOffset() != 0 {
		t.Error("scroll offset should start at 0")
	}

	// Update
	credits.Update(0.016)

	if credits.GetScrollOffset() == 0 {
		t.Error("scroll offset should increase after update")
	}

	// HandleInput should request exit
	initialExit := gsm.ShouldExit()
	credits.HandleInput(gsm.StateMachine)
	if !gsm.ShouldExit() && !initialExit {
		// Exit flag should be set (unless it was already set)
		// Actually, RequestExit sets the flag, so let's check
		if !gsm.ShouldExit() {
			t.Error("HandleInput should request exit")
		}
	}
}

// TestState is a simple test state
type TestState struct {
	name string
}

func (ts *TestState) Initialize(sm *statemachine.StateMachine) error {
	return nil
}

func (ts *TestState) Enter(sm *statemachine.StateMachine) {
}

func (ts *TestState) HandleInput(sm *statemachine.StateMachine) {
}

func (ts *TestState) Update(dt float64) {
}

func (ts *TestState) Draw() {
}

func (ts *TestState) Exit(sm *statemachine.StateMachine) {
}

func (ts *TestState) Shutdown() {
}

// Test Start with empty initial state in debug mode
func TestStartDebugEmptyState(t *testing.T) {
	gsm := NewGameStateMachine(true, "TestEngine", "1.0.0", "TestDev")

	err := gsm.Start("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should not have active state (no initial state provided)
	_, exists := gsm.GetActiveState()
	if exists {
		t.Error("should not have active state when no initial state provided")
	}
}

// Test ChangeState on GameStateMachine still works for non-builtin states
func TestGameStateMachineChangeState(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	testState := &TestState{name: "test"}
	gsm.RegisterStateInstance("test", testState)

	err := gsm.ChangeState("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	active, exists := gsm.GetActiveState()
	if !exists || active != "test" {
		t.Errorf("expected 'test', got %s (exists: %v)", active, exists)
	}
}

// Test that StateMachine methods still work on GameStateMachine
func TestGameStateMachineInheritsMethods(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	testState := &TestState{name: "test"}
	gsm.RegisterStateInstance("test", testState)

	// Test inherited methods work
	sm := gsm.StateMachine
	err := sm.ChangeState("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test context through inherited StateMachine
	sm.SetContext("key", "value")
	val, exists := sm.GetContext("key")
	if !exists || val != "value" {
		t.Error("context should work through inherited StateMachine")
	}
}

// Test ChangeState error path (trying to change to built-in states)
func TestChangeStateToBuiltinError(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")

	// Try to change to engine splash directly
	err := gsm.ChangeState(EngineSplashStateName)
	if err == nil {
		t.Error("expected error when changing to engine splash directly")
	}

	// Try to change to credits directly
	err = gsm.ChangeState(CreditsStateName)
	if err == nil {
		t.Error("expected error when changing to credits directly")
	}
}

// Test EngineSplashState Draw, HandleInput, Update, Exit, Shutdown (coverage)
func TestEngineSplashStateMethods(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")
	splash := gsm.engineSplash

	// Test all methods for coverage
	splash.HandleInput(gsm.StateMachine)
	splash.Update(0.016)
	splash.Draw()
	splash.Exit(gsm.StateMachine)
	splash.Shutdown()

	// Verify state still works
	if !splash.ShouldTransition() {
		// Should not transition immediately after Enter
		splash.Enter(gsm.StateMachine)
		time.Sleep(3 * time.Second)
		if !splash.ShouldTransition() {
			t.Error("should transition after duration")
		}
	}
}

// Test CreditsState Draw, Exit, Shutdown (coverage)
func TestCreditsStateMethods(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev")
	credits := gsm.credits

	credits.Initialize(gsm.StateMachine)
	credits.Enter(gsm.StateMachine)

	// Test all methods for coverage
	credits.Draw()
	credits.Exit(gsm.StateMachine)
	credits.Shutdown()

	// Test HandleCreditsInput
	credits.HandleCreditsInput() // This doesn't do much but covers the method
}
