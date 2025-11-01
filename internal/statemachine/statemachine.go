package statemachine

import (
	"fmt"
	"sync"
)

// State represents a state in the state machine with its lifecycle methods
type State interface {
	// Initialize is called once when the state is first created
	Initialize(sm *StateMachine) error

	// Enter is called every time the state becomes active
	Enter(sm *StateMachine)

	// HandleInput processes user input while the state is active
	HandleInput(sm *StateMachine)

	// Update updates game logic each frame while the state is active
	Update(dt float64)

	// Draw renders the state
	Draw()

	// Exit is called when leaving the state (but state may still be in memory)
	Exit(sm *StateMachine)

	// Shutdown is called once when the state is destroyed
	Shutdown()
}

// LuaCallbacks holds Lua callback functions for a state
type LuaCallbacks struct {
	Initialize  func(*StateMachine) error
	Enter       func(*StateMachine)
	HandleInput func(*StateMachine)
	Update      func(float64)
	Draw        func()
	Exit        func(*StateMachine)
	Shutdown    func()
}

// LuaState wraps Lua callbacks to implement the State interface
type LuaState struct {
	name      string
	callbacks LuaCallbacks
}

// Initialize calls the Lua initialize callback
func (ls *LuaState) Initialize(sm *StateMachine) error {
	if ls.callbacks.Initialize != nil {
		return ls.callbacks.Initialize(sm)
	}
	return nil
}

// Enter calls the Lua enter callback
func (ls *LuaState) Enter(sm *StateMachine) {
	if ls.callbacks.Enter != nil {
		ls.callbacks.Enter(sm)
	}
}

// HandleInput calls the Lua handleInput callback
func (ls *LuaState) HandleInput(sm *StateMachine) {
	if ls.callbacks.HandleInput != nil {
		ls.callbacks.HandleInput(sm)
	}
}

// Update calls the Lua update callback
func (ls *LuaState) Update(dt float64) {
	if ls.callbacks.Update != nil {
		// Use a defer/recover to catch any panics from Update callbacks
		defer func() {
			if r := recover(); r != nil {
				// Panic in Update - log but don't crash
				// State machine should continue running
				_ = r // Silently recover - game should keep running
			}
		}()
		ls.callbacks.Update(dt)
	}
}

// Draw calls the Lua draw callback
func (ls *LuaState) Draw() {
	if ls.callbacks.Draw != nil {
		ls.callbacks.Draw()
	}
}

// Exit calls the Lua exit callback
func (ls *LuaState) Exit(sm *StateMachine) {
	if ls.callbacks.Exit != nil {
		ls.callbacks.Exit(sm)
	}
}

// Shutdown calls the Lua shutdown callback
func (ls *LuaState) Shutdown() {
	if ls.callbacks.Shutdown != nil {
		ls.callbacks.Shutdown()
	}
}

// StateMachine manages game states with support for state stacking and shared context
type StateMachine struct {
	mu sync.RWMutex

	stateStack    []State                // Stack of active states
	stateRegistry map[string]State       // All registered states
	initialized   map[string]bool        // Track which states are initialized
	context       map[string]interface{} // Shared context for data passing
	shouldExit    bool                   // Flag to request exit

	// Optional: hook for drawing previous state
	drawPreviousStateHook func()

	// Pending state changes (to avoid deadlocks when called from HandleInput/Update/Draw)
	pendingChangeState string // Queue a ChangeState operation
	pendingPushState   string // Queue a PushState operation
	pendingPopState    bool   // Queue a PopState operation
	inCallback         bool   // Flag to detect if we're in a callback (HandleInput/Update/Draw)
}

// NewStateMachine creates a new generic state machine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		stateStack:    make([]State, 0),
		stateRegistry: make(map[string]State),
		initialized:   make(map[string]bool),
		context:       make(map[string]interface{}),
		shouldExit:    false,
	}
}

// RegisterState registers a new state with the given name
func (sm *StateMachine) RegisterState(name string, callbacks LuaCallbacks) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if name == "" {
		return fmt.Errorf("state name cannot be empty")
	}

	if _, exists := sm.stateRegistry[name]; exists {
		return fmt.Errorf("state '%s' is already registered", name)
	}

	sm.stateRegistry[name] = &LuaState{
		name:      name,
		callbacks: callbacks,
	}

	return nil
}

// RegisterStateInstance registers a State instance directly
func (sm *StateMachine) RegisterStateInstance(name string, state State) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if name == "" {
		return fmt.Errorf("state name cannot be empty")
	}

	if _, exists := sm.stateRegistry[name]; exists {
		return fmt.Errorf("state '%s' is already registered", name)
	}

	sm.stateRegistry[name] = state

	return nil
}

// UnregisterState removes a state from the registry
func (sm *StateMachine) UnregisterState(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	state, exists := sm.stateRegistry[name]
	if !exists {
		return fmt.Errorf("state '%s' is not registered", name)
	}

	// If state is in stack, we can't unregister it
	for _, s := range sm.stateStack {
		if s == state {
			return fmt.Errorf("cannot unregister state '%s' while it's in the stack", name)
		}
	}

	// Shutdown if initialized
	if sm.initialized[name] {
		state.Shutdown()
		delete(sm.initialized, name)
	}

	delete(sm.stateRegistry, name)
	return nil
}

// ChangeState replaces all states in the stack with a new state
func (sm *StateMachine) ChangeState(name string) error {
	sm.mu.RLock()
	inCallback := sm.inCallback
	sm.mu.RUnlock()

	// If called from within HandleInput/Update/Draw, queue it
	if inCallback {
		sm.mu.Lock()
		sm.pendingChangeState = name
		sm.pendingPushState = ""
		sm.pendingPopState = false
		sm.mu.Unlock()
		return nil
	}

	// Otherwise execute immediately
	return sm.doChangeState(name)
}

// doChangeState performs the actual state change (must be called without locks or with write lock)
func (sm *StateMachine) doChangeState(name string) error {
	sm.mu.Lock()
	state, exists := sm.stateRegistry[name]
	if !exists {
		sm.mu.Unlock()
		return fmt.Errorf("state '%s' is not registered", name)
	}

	// Collect states to exit (while holding lock)
	statesToExit := make([]State, len(sm.stateStack))
	copy(statesToExit, sm.stateStack)
	sm.stateStack = sm.stateStack[:0] // Clear stack
	sm.mu.Unlock()

	// Exit all states (outside lock to avoid deadlock)
	for _, top := range statesToExit {
		top.Exit(sm)
	}

	sm.mu.Lock()
	// Initialize if first time
	if !sm.initialized[name] {
		sm.mu.Unlock()
		if err := state.Initialize(sm); err != nil {
			return fmt.Errorf("failed to initialize state '%s': %w", name, err)
		}
		sm.mu.Lock()
		sm.initialized[name] = true
	}

	// Add to stack
	sm.stateStack = append(sm.stateStack, state)
	sm.mu.Unlock()

	// Enter state (outside lock to avoid deadlock)
	state.Enter(sm)

	return nil
}

// PushState adds a new state on top of the stack
func (sm *StateMachine) PushState(name string) error {
	sm.mu.RLock()
	inCallback := sm.inCallback
	sm.mu.RUnlock()

	// If called from within HandleInput/Update/Draw, queue it
	if inCallback {
		sm.mu.Lock()
		sm.pendingPushState = name
		sm.pendingChangeState = ""
		sm.pendingPopState = false
		sm.mu.Unlock()
		return nil
	}

	// Otherwise execute immediately
	return sm.doPushState(name)
}

// doPushState performs the actual push (must be called without locks or with write lock)
func (sm *StateMachine) doPushState(name string) error {
	sm.mu.Lock()
	state, exists := sm.stateRegistry[name]
	if !exists {
		sm.mu.Unlock()
		return fmt.Errorf("state '%s' is not registered", name)
	}

	var topToExit State
	if len(sm.stateStack) > 0 {
		topToExit = sm.stateStack[len(sm.stateStack)-1]
	}
	sm.mu.Unlock()

	// Exit current top state (but keep on stack) - outside lock to avoid deadlock
	if topToExit != nil {
		topToExit.Exit(sm)
	}

	sm.mu.Lock()
	// Initialize if first time
	if !sm.initialized[name] {
		sm.mu.Unlock()
		if err := state.Initialize(sm); err != nil {
			return fmt.Errorf("failed to initialize state '%s': %w", name, err)
		}
		sm.mu.Lock()
		sm.initialized[name] = true
	}

	// Add to stack
	sm.stateStack = append(sm.stateStack, state)
	sm.mu.Unlock()

	// Enter state (outside lock to avoid deadlock)
	state.Enter(sm)

	return nil
}

// PopState removes the top state from the stack
func (sm *StateMachine) PopState() error {
	sm.mu.RLock()
	inCallback := sm.inCallback
	sm.mu.RUnlock()

	// If called from within HandleInput/Update/Draw, queue it
	if inCallback {
		sm.mu.Lock()
		sm.pendingPopState = true
		sm.pendingChangeState = ""
		sm.pendingPushState = ""
		sm.mu.Unlock()
		return nil
	}

	// Otherwise execute immediately
	return sm.doPopState()
}

// doPopState performs the actual pop (must be called without locks or with write lock)
func (sm *StateMachine) doPopState() error {
	sm.mu.Lock()
	if len(sm.stateStack) == 0 {
		sm.mu.Unlock()
		return fmt.Errorf("cannot pop state from empty stack")
	}

	// Get states (while holding lock)
	top := sm.stateStack[len(sm.stateStack)-1]
	sm.stateStack = sm.stateStack[:len(sm.stateStack)-1]
	var previous State
	if len(sm.stateStack) > 0 {
		previous = sm.stateStack[len(sm.stateStack)-1]
	}
	sm.mu.Unlock()

	// Exit top state (outside lock to avoid deadlock)
	top.Exit(sm)

	// Re-enter previous state if any (outside lock to avoid deadlock)
	if previous != nil {
		previous.Enter(sm)
	}

	return nil
}

// PopAllStates removes all states from the stack
func (sm *StateMachine) PopAllStates() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for len(sm.stateStack) > 0 {
		top := sm.stateStack[len(sm.stateStack)-1]
		top.Exit(sm)
		sm.stateStack = sm.stateStack[:len(sm.stateStack)-1]
	}
}

// HandleInput calls HandleInput on the top state
func (sm *StateMachine) HandleInput() {
	sm.mu.Lock()
	sm.inCallback = true
	sm.mu.Unlock()

	sm.mu.RLock()
	var top State
	if len(sm.stateStack) > 0 {
		top = sm.stateStack[len(sm.stateStack)-1]
	}
	sm.mu.RUnlock()

	if top != nil {
		top.HandleInput(sm)
	}

	sm.mu.Lock()
	sm.inCallback = false
	sm.mu.Unlock()
	// Don't process state changes here - defer to Update() to avoid deadlocks
}

// Update calls Update on the top state
func (sm *StateMachine) Update(dt float64) {
	// Process pending state changes FIRST (before calling Update on current state)
	// This ensures state changes from previous frame are handled before the new state's Update runs
	sm.mu.Lock()
	pendingChange := sm.pendingChangeState != ""
	pendingPush := sm.pendingPushState != ""
	pendingPop := sm.pendingPopState

	var changeStateName string
	var pushStateName string

	if pendingChange {
		changeStateName = sm.pendingChangeState
		sm.pendingChangeState = ""
		sm.pendingPushState = ""
		sm.pendingPopState = false
	} else if pendingPush {
		pushStateName = sm.pendingPushState
		sm.pendingPushState = ""
		sm.pendingChangeState = ""
		sm.pendingPopState = false
	} else if pendingPop {
		sm.pendingPopState = false
		sm.pendingChangeState = ""
		sm.pendingPushState = ""
	}
	sm.mu.Unlock()

	// Execute pending state changes outside lock to avoid deadlock
	if pendingChange {
		if err := sm.doChangeState(changeStateName); err != nil {
			// State change failed - don't call Update this frame to avoid issues
			sm.mu.Lock()
			sm.inCallback = false
			sm.mu.Unlock()
			return
		}
		// After state change, the new state is now on top, so it will get Update called below
	} else if pendingPush {
		if err := sm.doPushState(pushStateName); err != nil {
			// State push failed - don't call Update this frame
			sm.mu.Lock()
			sm.inCallback = false
			sm.mu.Unlock()
			return
		}
		// After push, the new state is now on top, so it will get Update called below
	} else if pendingPop {
		if err := sm.doPopState(); err != nil {
			// State pop failed - don't call Update this frame
			sm.mu.Lock()
			sm.inCallback = false
			sm.mu.Unlock()
			return
		}
		// After pop, the previous state is now on top, so it will get Update called below
	}

	// Now call Update on the current top state (which may have just changed)
	// Re-read the stack after state changes to ensure we have the latest state
	sm.mu.Lock()
	sm.inCallback = true
	sm.mu.Unlock()

	sm.mu.RLock()
	var top State
	if len(sm.stateStack) > 0 {
		top = sm.stateStack[len(sm.stateStack)-1]
	}
	sm.mu.RUnlock()

	if top != nil {
		top.Update(dt)
	}

	sm.mu.Lock()
	sm.inCallback = false
	sm.mu.Unlock()
}

// Draw calls Draw on the top state
func (sm *StateMachine) Draw() {
	sm.mu.Lock()
	sm.inCallback = true
	sm.mu.Unlock()

	sm.mu.RLock()
	var top State
	if len(sm.stateStack) > 0 {
		top = sm.stateStack[len(sm.stateStack)-1]
	}
	sm.mu.RUnlock()

	if top != nil {
		top.Draw()
	}

	sm.mu.Lock()
	sm.inCallback = false
	sm.mu.Unlock()
	// Don't process state changes here - defer to Update() to avoid deadlocks
}

// DrawPreviousState sets a hook function that will be called to draw the previous state
func (sm *StateMachine) SetDrawPreviousStateHook(hook func()) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.drawPreviousStateHook = hook
}

// DrawPreviousState calls the hook to draw the previous state (if available)
func (sm *StateMachine) DrawPreviousState() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.drawPreviousStateHook != nil && len(sm.stateStack) > 1 {
		sm.drawPreviousStateHook()
	}
}

// Context Management

// SetContext sets a value in the shared context
func (sm *StateMachine) SetContext(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.context[key] = value
}

// GetContext gets a value from the shared context
func (sm *StateMachine) GetContext(key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.context[key]
	return value, exists
}

// HasContext checks if a key exists in the context
func (sm *StateMachine) HasContext(key string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.context[key]
	return exists
}

// ClearContext removes a specific key from the context
func (sm *StateMachine) ClearContext(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.context, key)
}

// ClearAllContext clears all context data
func (sm *StateMachine) ClearAllContext() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.context = make(map[string]interface{})
}

// Control

// RequestExit sets the exit flag
func (sm *StateMachine) RequestExit() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.shouldExit = true
}

// ShouldExit returns whether exit was requested
func (sm *StateMachine) ShouldExit() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.shouldExit
}

// Utility Methods

// GetActiveState returns the name of the currently active state (if registered)
func (sm *StateMachine) GetActiveState() (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.stateStack) == 0 {
		return "", false
	}

	active := sm.stateStack[len(sm.stateStack)-1]

	// Try to find the name
	for name, state := range sm.stateRegistry {
		if state == active {
			return name, true
		}
	}

	return "", false
}

// GetStackDepth returns the number of states in the stack
func (sm *StateMachine) GetStackDepth() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.stateStack)
}

// IsStateRegistered checks if a state is registered
func (sm *StateMachine) IsStateRegistered(name string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.stateRegistry[name]
	return exists
}
