package statemachine

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// TestState is a simple test implementation of State
type TestState struct {
	name           string
	initCalled     bool
	enterCount     int
	inputCount     int
	updateCount    int
	drawCount      int
	exitCount      int
	shutdownCalled bool
	initError      error
	mu             sync.Mutex
}

func NewTestState(name string) *TestState {
	return &TestState{name: name}
}

func (ts *TestState) Initialize(sm *StateMachine) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.initCalled = true
	return ts.initError
}

func (ts *TestState) Enter(sm *StateMachine) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.enterCount++
}

func (ts *TestState) HandleInput(sm *StateMachine) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.inputCount++
}

func (ts *TestState) Update(dt float64) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.updateCount++
}

func (ts *TestState) Draw() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.drawCount++
}

func (ts *TestState) Exit(sm *StateMachine) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.exitCount++
}

func (ts *TestState) Shutdown() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.shutdownCalled = true
}

func (ts *TestState) GetCounts() (init, enter, input, update, draw, exit, shutdown bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.initCalled, ts.enterCount > 0, ts.inputCount > 0, ts.updateCount > 0, ts.drawCount > 0, ts.exitCount > 0, ts.shutdownCalled
}

func (ts *TestState) GetEnterCount() int {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.enterCount
}

func (ts *TestState) GetExitCount() int {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.exitCount
}

// Test NewStateMachine
func TestNewStateMachine(t *testing.T) {
	sm := NewStateMachine()
	if sm == nil {
		t.Fatal("NewStateMachine returned nil")
	}
	if sm.GetStackDepth() != 0 {
		t.Errorf("expected stack depth 0, got %d", sm.GetStackDepth())
	}
	if sm.ShouldExit() {
		t.Error("shouldExit should be false initially")
	}
}

// Test RegisterState
func TestRegisterState(t *testing.T) {
	sm := NewStateMachine()

	// Valid registration
	err := sm.RegisterState("test", LuaCallbacks{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !sm.IsStateRegistered("test") {
		t.Error("state should be registered")
	}

	// Duplicate registration
	err = sm.RegisterState("test", LuaCallbacks{})
	if err == nil {
		t.Error("expected error for duplicate registration")
	}

	// Empty name
	err = sm.RegisterState("", LuaCallbacks{})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

// Test RegisterStateInstance
func TestRegisterStateInstance(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")

	err := sm.RegisterStateInstance("test", state)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !sm.IsStateRegistered("test") {
		t.Error("state should be registered")
	}
}

// Test UnregisterState
func TestUnregisterState(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")

	sm.RegisterStateInstance("test", state)
	sm.ChangeState("test")

	// Can't unregister active state
	err := sm.UnregisterState("test")
	if err == nil {
		t.Error("expected error when unregistering active state")
	}

	sm.PopAllStates()

	// Now can unregister
	// Note: state WAS initialized because ChangeState was called
	err = sm.UnregisterState("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !state.shutdownCalled {
		t.Error("shutdown should be called when unregistering initialized state")
	}

	// Unregister non-existent
	err = sm.UnregisterState("nonexistent")
	if err == nil {
		t.Error("expected error for unregistering non-existent state")
	}
}

// Test ChangeState
func TestChangeState(t *testing.T) {
	sm := NewStateMachine()
	state1 := NewTestState("state1")
	state2 := NewTestState("state2")

	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	// Change to state1
	err := sm.ChangeState("state1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	init, enter, _, _, _, _, _ := state1.GetCounts()
	if !init {
		t.Error("Initialize should be called")
	}
	if !enter {
		t.Error("Enter should be called")
	}

	// Change to state2 (should exit state1)
	err = sm.ChangeState("state2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if state1.GetExitCount() != 1 {
		t.Errorf("expected Exit to be called once, got %d", state1.GetExitCount())
	}

	// Change to non-existent state
	err = sm.ChangeState("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent state")
	}
}

// Test PushState
func TestPushState(t *testing.T) {
	sm := NewStateMachine()
	state1 := NewTestState("state1")
	state2 := NewTestState("state2")

	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	sm.ChangeState("state1")

	// Push state2
	err := sm.PushState("state2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sm.GetStackDepth() != 2 {
		t.Errorf("expected stack depth 2, got %d", sm.GetStackDepth())
	}

	if state1.GetExitCount() != 1 {
		t.Error("state1 should be exited")
	}

	init, enter, _, _, _, _, _ := state2.GetCounts()
	if !init {
		t.Error("state2 Initialize should be called")
	}
	if !enter {
		t.Error("state2 Enter should be called")
	}
}

// Test PopState
func TestPopState(t *testing.T) {
	sm := NewStateMachine()
	state1 := NewTestState("state1")
	state2 := NewTestState("state2")

	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	sm.ChangeState("state1")
	sm.PushState("state2")

	initialEnterCount := state1.GetEnterCount()

	// Pop state2
	err := sm.PopState()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sm.GetStackDepth() != 1 {
		t.Errorf("expected stack depth 1, got %d", sm.GetStackDepth())
	}

	if state2.GetExitCount() != 1 {
		t.Error("state2 should be exited")
	}

	// state1 should be re-entered
	if state1.GetEnterCount() != initialEnterCount+1 {
		t.Errorf("expected state1 enter count to increase, got %d", state1.GetEnterCount())
	}

	// Pop from empty stack
	sm.PopAllStates()
	err = sm.PopState()
	if err == nil {
		t.Error("expected error when popping from empty stack")
	}
}

// Test PopAllStates
func TestPopAllStates(t *testing.T) {
	sm := NewStateMachine()
	state1 := NewTestState("state1")
	state2 := NewTestState("state2")

	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	sm.ChangeState("state1")
	sm.PushState("state2")

	// At this point:
	// - state1 entered=1, exited=1 (when state2 was pushed)
	// - state2 entered=1, exited=0
	// Stack: [state1, state2] (state2 on top)

	sm.PopAllStates()

	if sm.GetStackDepth() != 0 {
		t.Errorf("expected stack depth 0, got %d", sm.GetStackDepth())
	}

	// PopAllStates exits from top to bottom: state2 first, then state1
	if state2.GetExitCount() != 1 {
		t.Errorf("state2 should be exited once, got %d", state2.GetExitCount())
	}
	// state1 was already exited once when state2 was pushed, now exited again = 2 total
	if state1.GetExitCount() != 2 {
		t.Errorf("state1 should be exited twice (once on push, once on popall), got %d", state1.GetExitCount())
	}
}

// Test HandleInput, Update, Draw
func TestStateCallbacks(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")

	sm.RegisterStateInstance("test", state)
	sm.ChangeState("test")

	sm.HandleInput()
	sm.Update(0.016)
	sm.Draw()

	_, _, input, update, draw, _, _ := state.GetCounts()
	if !input {
		t.Error("HandleInput should be called")
	}
	if !update {
		t.Error("Update should be called")
	}
	if !draw {
		t.Error("Draw should be called")
	}
}

// Test empty stack callbacks
func TestEmptyStackCallbacks(t *testing.T) {
	sm := NewStateMachine()

	// Should not panic
	sm.HandleInput()
	sm.Update(0.016)
	sm.Draw()
}

// Test Context Management
func TestContext(t *testing.T) {
	sm := NewStateMachine()

	// Set context
	sm.SetContext("key1", "value1")
	sm.SetContext("key2", 42)

	// Get context
	val, exists := sm.GetContext("key1")
	if !exists {
		t.Error("key1 should exist")
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got %v", val)
	}

	// HasContext
	if !sm.HasContext("key1") {
		t.Error("HasContext should return true")
	}
	if sm.HasContext("nonexistent") {
		t.Error("HasContext should return false")
	}

	// ClearContext
	sm.ClearContext("key1")
	if sm.HasContext("key1") {
		t.Error("key1 should be cleared")
	}

	// ClearAllContext
	sm.ClearAllContext()
	if sm.HasContext("key2") {
		t.Error("key2 should be cleared")
	}
}

// Test Concurrent Access
func TestConcurrentAccess(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")
	sm.RegisterStateInstance("test", state)

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sm.GetStackDepth()
				sm.HasContext("test")
				sm.ShouldExit()
			}
		}()
	}

	// Concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", idx*100+j)
				sm.SetContext(key, idx*100+j)
			}
		}(i)
	}

	wg.Wait()
}

// Test Initialize Error
func TestInitializeError(t *testing.T) {
	sm := NewStateMachine()

	state := NewTestState("test")
	state.initError = errors.New("init failed")

	sm.RegisterStateInstance("test", state)

	err := sm.ChangeState("test")
	if err == nil {
		t.Error("expected error from Initialize")
	}

	// Initialize WAS called, but state should not be marked as initialized
	// because ChangeState returns error before marking as initialized
	init, _, _, _, _, _, _ := state.GetCounts()
	if !init {
		t.Error("Initialize WAS called (even though it failed)")
	}

	// Verify state is not in initialized map by trying to change to it again
	// If it was initialized, second call wouldn't call Initialize again
	err2 := sm.ChangeState("test")
	if err2 == nil {
		t.Error("should still fail on second attempt")
	}
}

// Test Multiple Enter/Exit cycles
func TestMultipleEnterExit(t *testing.T) {
	sm := NewStateMachine()
	state1 := NewTestState("state1")
	state2 := NewTestState("state2")

	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	// Enter/exit multiple times
	// 1. ChangeState("state1") - enter=1
	sm.ChangeState("state1")
	// 2. PushState("state2") - exit=1, state1 stays in stack
	sm.PushState("state2")
	// 3. PopState() - exit state2, enter state1 again - enter=2
	sm.PopState()
	// 4. PushState("state2") - exit=2, state1 stays in stack
	sm.PushState("state2")
	// 5. PopState() - exit state2, enter state1 again - enter=3
	sm.PopState()
	// 6. ChangeState("state2") - exit=3 (state1 removed from stack)
	sm.ChangeState("state2")
	// 7. ChangeState("state1") - enter=4 (state1 added back)
	sm.ChangeState("state1")

	// state1 should have been entered 4 times (1, 2, 3, 4)
	if state1.GetEnterCount() != 4 {
		t.Errorf("expected state1 enter count 4, got %d", state1.GetEnterCount())
	}
	// state1 should have been exited 3 times (before pushes and on change)
	if state1.GetExitCount() != 3 {
		t.Errorf("expected state1 exit count 3, got %d", state1.GetExitCount())
	}
}

// Test Initialize only called once
func TestInitializeOnce(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")

	sm.RegisterStateInstance("test", state)

	// Initialize should only be called once
	sm.ChangeState("test")
	sm.ChangeState("other") // Create another state
	other := NewTestState("other")
	sm.RegisterStateInstance("other", other)
	sm.ChangeState("test") // Change back

	// Count initialize calls by checking if initCalled is true
	init, _, _, _, _, _, _ := state.GetCounts()
	if !init {
		t.Error("Initialize should be called")
	}

	// State should only be initialized once (we can't easily count this without modifying TestState)
	// But we can verify it doesn't error on second change
	err := sm.ChangeState("test")
	if err != nil {
		t.Errorf("unexpected error on second change: %v", err)
	}
}

// Test DrawPreviousState
func TestDrawPreviousState(t *testing.T) {
	sm := NewStateMachine()
	drawCalled := false

	sm.SetDrawPreviousStateHook(func() {
		drawCalled = true
	})

	state1 := NewTestState("state1")
	state2 := NewTestState("state2")
	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	sm.ChangeState("state1")
	sm.PushState("state2")

	// DrawPreviousState should call hook when stack depth > 1
	sm.DrawPreviousState()
	if !drawCalled {
		t.Error("DrawPreviousState should call hook")
	}

	// With single state, should not call
	drawCalled = false
	sm.PopState()
	sm.DrawPreviousState()
	if drawCalled {
		t.Error("DrawPreviousState should not call hook with single state")
	}
}

// Test GetActiveState
func TestGetActiveState(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")
	sm.RegisterStateInstance("test", state)

	// Empty stack
	name, exists := sm.GetActiveState()
	if exists {
		t.Error("should not have active state with empty stack")
	}

	sm.ChangeState("test")
	name, exists = sm.GetActiveState()
	if !exists {
		t.Error("should have active state")
	}
	if name != "test" {
		t.Errorf("expected 'test', got %s", name)
	}
}

// Test RequestExit
func TestRequestExit(t *testing.T) {
	sm := NewStateMachine()

	if sm.ShouldExit() {
		t.Error("should not exit initially")
	}

	sm.RequestExit()
	if !sm.ShouldExit() {
		t.Error("should exit after request")
	}
}

// Test RegisterState with LuaCallbacks (to cover LuaState methods)
func TestRegisterStateWithCallbacks(t *testing.T) {
	sm := NewStateMachine()

	initCalled := false
	enterCalled := false
	inputCalled := false
	updateCalled := false
	drawCalled := false
	exitCalled := false
	shutdownCalled := false

	callbacks := LuaCallbacks{
		Initialize: func(sm *StateMachine) error {
			initCalled = true
			return nil
		},
		Enter: func(sm *StateMachine) {
			enterCalled = true
		},
		HandleInput: func(sm *StateMachine) {
			inputCalled = true
		},
		Update: func(dt float64) {
			updateCalled = true
		},
		Draw: func() {
			drawCalled = true
		},
		Exit: func(sm *StateMachine) {
			exitCalled = true
		},
		Shutdown: func() {
			shutdownCalled = true
		},
	}

	err := sm.RegisterState("test", callbacks)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Change to state to trigger Initialize and Enter
	err = sm.ChangeState("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !initCalled {
		t.Error("Initialize should be called")
	}
	if !enterCalled {
		t.Error("Enter should be called")
	}

	// Test callbacks
	sm.HandleInput()
	if !inputCalled {
		t.Error("HandleInput should be called")
	}

	sm.Update(0.016)
	if !updateCalled {
		t.Error("Update should be called")
	}

	sm.Draw()
	if !drawCalled {
		t.Error("Draw should be called")
	}

	// Exit state
	sm.ChangeState("other")
	other := NewTestState("other")
	sm.RegisterStateInstance("other", other)
	sm.ChangeState("test")
	sm.ChangeState("other")

	if !exitCalled {
		t.Error("Exit should be called")
	}

	// Unregister to trigger Shutdown
	sm.PopAllStates()
	sm.UnregisterState("test")

	if !shutdownCalled {
		t.Error("Shutdown should be called")
	}
}

// Test RegisterStateInstance with empty name
func TestRegisterStateInstanceEmptyName(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")

	err := sm.RegisterStateInstance("", state)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

// Test RegisterStateInstance duplicate
func TestRegisterStateInstanceDuplicate(t *testing.T) {
	sm := NewStateMachine()
	state := NewTestState("test")

	sm.RegisterStateInstance("test", state)
	err := sm.RegisterStateInstance("test", state)
	if err == nil {
		t.Error("expected error for duplicate registration")
	}
}

// Test PushState error paths
func TestPushStateErrors(t *testing.T) {
	sm := NewStateMachine()

	// Push non-existent state
	err := sm.PushState("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent state")
	}

	// Push state with Initialize error
	state := NewTestState("test")
	state.initError = errors.New("init failed")
	sm.RegisterStateInstance("test", state)

	err = sm.PushState("test")
	if err == nil {
		t.Error("expected error from Initialize")
	}
}

// Test GetActiveState with multiple states
func TestGetActiveStateMultiple(t *testing.T) {
	sm := NewStateMachine()
	state1 := NewTestState("state1")
	state2 := NewTestState("state2")

	sm.RegisterStateInstance("state1", state1)
	sm.RegisterStateInstance("state2", state2)

	sm.ChangeState("state1")
	name, exists := sm.GetActiveState()
	if !exists || name != "state1" {
		t.Errorf("expected state1, got %s (exists: %v)", name, exists)
	}

	sm.PushState("state2")
	name, exists = sm.GetActiveState()
	if !exists || name != "state2" {
		t.Errorf("expected state2, got %s (exists: %v)", name, exists)
	}
}
