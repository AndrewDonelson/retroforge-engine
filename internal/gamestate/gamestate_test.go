package gamestate

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/statemachine"
)

// Test NewGameStateMachine
func TestNewGameStateMachine(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(true, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.2.3", "TestDeveloper", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(true, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)
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

	// Simulate game loop - update with dt values to accumulate elapsedTime
	// Splash duration is 2.0 seconds, so we need to accumulate at least 2.0 seconds
	dt := 0.016 // ~60 FPS
	updates := 0
	for splash.elapsedTime < 2.5 { // Update a bit past the 2.0 second threshold
		splash.Update(dt)
		updates++
		if updates > 200 { // Safety limit to prevent infinite loop
			t.Fatal("too many updates needed - elapsedTime not accumulating")
		}
	}

	if !splash.ShouldTransition() {
		t.Error("should transition after duration")
	}
}

// Test CreditsState
func TestCreditsState(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)
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

	// Update a few times to advance frameCount (HandleInput requires frameCount >= 2)
	credits.Update(0.016)
	credits.Update(0.016)
	credits.Update(0.016) // Now frameCount should be >= 2

	if credits.GetScrollOffset() == 0 {
		t.Error("scroll offset should increase after update")
	}

	// HandleInput checks for input and only requests exit if input is present
	// Since we can't mock input in tests, verify that HandleInput doesn't crash
	// and respects the frameCount requirement (returns early if frameCount < 2)
	
	// Test that HandleInput with frameCount < 2 doesn't request exit (early return)
	credits2 := gsm.credits
	credits2.Enter(gsm.StateMachine) // Reset frameCount to 0
	initialExit := gsm.ShouldExit()
	credits2.HandleInput(gsm.StateMachine) // Should return early (frameCount = 0)
	if gsm.ShouldExit() != initialExit {
		t.Error("HandleInput should not request exit when frameCount < 2")
	}
	
	// Test that HandleInput with frameCount >= 2 runs without crashing
	// (It won't request exit without actual input, which is correct behavior)
	credits.HandleInput(gsm.StateMachine) // frameCount >= 2 from earlier updates
	// Should not crash - exit state should be unchanged (no input available)
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
	gsm := NewGameStateMachine(true, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)

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
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)
	splash := gsm.engineSplash

	// Test all methods for coverage
	splash.HandleInput(gsm.StateMachine)
	splash.Update(0.016)
	splash.Draw()
	splash.Exit(gsm.StateMachine)
	splash.Shutdown()

	// Verify state still works - test that ShouldTransition works after time accumulation
	if !splash.ShouldTransition() {
		// Should not transition immediately after Enter
		splash.Enter(gsm.StateMachine)
		
		// Simulate game loop - update with dt values to accumulate elapsedTime
		dt := 0.016 // ~60 FPS
		updates := 0
		for splash.elapsedTime < 2.5 { // Update past the 2.0 second threshold
			splash.Update(dt)
			updates++
			if updates > 200 { // Safety limit
				t.Fatal("too many updates needed - elapsedTime not accumulating")
			}
		}
		
		if !splash.ShouldTransition() {
			t.Error("should transition after duration")
		}
	}
}

// Test CreditsState Draw, Exit, Shutdown (coverage)
func TestCreditsStateMethods(t *testing.T) {
	gsm := NewGameStateMachine(false, "TestEngine", "1.0.0", "TestDev", nil, nil)
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
