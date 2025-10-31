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
	sm.mu.Lock()
	defer sm.mu.Unlock()

	state, exists := sm.stateRegistry[name]
	if !exists {
		return fmt.Errorf("state '%s' is not registered", name)
	}

	// Exit and remove all states from stack
	for len(sm.stateStack) > 0 {
		top := sm.stateStack[len(sm.stateStack)-1]
		top.Exit(sm)
		sm.stateStack = sm.stateStack[:len(sm.stateStack)-1]
	}

	// Initialize if first time
	if !sm.initialized[name] {
		if err := state.Initialize(sm); err != nil {
			return fmt.Errorf("failed to initialize state '%s': %w", name, err)
		}
		sm.initialized[name] = true
	}

	// Add to stack and enter
	sm.stateStack = append(sm.stateStack, state)
	state.Enter(sm)

	return nil
}

// PushState adds a new state on top of the stack
func (sm *StateMachine) PushState(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	state, exists := sm.stateRegistry[name]
	if !exists {
		return fmt.Errorf("state '%s' is not registered", name)
	}

	// Exit current top state (but keep on stack)
	if len(sm.stateStack) > 0 {
		top := sm.stateStack[len(sm.stateStack)-1]
		top.Exit(sm)
	}

	// Initialize if first time
	if !sm.initialized[name] {
		if err := state.Initialize(sm); err != nil {
			return fmt.Errorf("failed to initialize state '%s': %w", name, err)
		}
		sm.initialized[name] = true
	}

	// Add to stack and enter
	sm.stateStack = append(sm.stateStack, state)
	state.Enter(sm)

	return nil
}

// PopState removes the top state from the stack
func (sm *StateMachine) PopState() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if len(sm.stateStack) == 0 {
		return fmt.Errorf("cannot pop state from empty stack")
	}

	// Exit and remove top state
	top := sm.stateStack[len(sm.stateStack)-1]
	top.Exit(sm)
	sm.stateStack = sm.stateStack[:len(sm.stateStack)-1]

	// Re-enter previous state if any
	if len(sm.stateStack) > 0 {
		previous := sm.stateStack[len(sm.stateStack)-1]
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
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.stateStack) > 0 {
		top := sm.stateStack[len(sm.stateStack)-1]
		top.HandleInput(sm)
	}
}

// Update calls Update on the top state
func (sm *StateMachine) Update(dt float64) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.stateStack) > 0 {
		top := sm.stateStack[len(sm.stateStack)-1]
		top.Update(dt)
	}
}

// Draw calls Draw on the top state
func (sm *StateMachine) Draw() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.stateStack) > 0 {
		top := sm.stateStack[len(sm.stateStack)-1]
		top.Draw()
	}
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
