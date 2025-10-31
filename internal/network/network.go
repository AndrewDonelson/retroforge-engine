package network

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SyncTier represents the synchronization frequency tier
type SyncTier string

const (
	SyncTierFast     SyncTier = "fast"     // 30-60/sec
	SyncTierModerate SyncTier = "moderate" // 15/sec
	SyncTierSlow     SyncTier = "slow"     // 5/sec
)

// PacketType represents the type of network packet
type PacketType string

const (
	PacketTypeInput        PacketType = "input"
	PacketTypeStateDelta   PacketType = "state_delta"
	PacketTypeFullState    PacketType = "full_state"
	PacketTypePlayerJoined PacketType = "player_joined"
	PacketTypePlayerLeft   PacketType = "player_left"
	PacketTypeGameOver     PacketType = "game_over"
)

// InputPacket represents player input data
type InputPacket struct {
	Type      PacketType      `json:"type"`
	PlayerID  int             `json:"player_id"`
	Frame     uint64          `json:"frame"`
	Buttons   map[string]bool `json:"buttons"`
	Timestamp float64         `json:"timestamp"`
}

// StateDeltaPacket represents incremental state changes
type StateDeltaPacket struct {
	Type      PacketType             `json:"type"`
	Frame     uint64                 `json:"frame"`
	Tier      SyncTier               `json:"tier"`
	Changes   map[string]interface{} `json:"changes"`
	Timestamp float64                `json:"timestamp"`
}

// FullStatePacket represents complete game state snapshot
type FullStatePacket struct {
	Type      PacketType             `json:"type"`
	Frame     uint64                 `json:"frame"`
	State     map[string]interface{} `json:"state"`
	Timestamp float64                `json:"timestamp"`
}

// Connection represents a connection to another player
type Connection interface {
	Send(data []byte) error
	Close() error
	IsConnected() bool
	PlayerID() int
}

// NetworkManager manages multiplayer networking
type NetworkManager struct {
	mu sync.RWMutex

	// Multiplayer state
	isMultiplayer bool
	isHost        bool
	playerID      int
	playerCount   int

	// Connections (only host has connections to other players)
	connections map[int]Connection // playerID -> Connection

	// Input state (host stores inputs from all players)
	playerInputs map[int]map[int]bool // playerID -> buttonID -> pressed

	// State synchronization
	syncedTables map[string]*SyncedTable // tablePath -> SyncedTable
	lastSyncTime map[string]time.Time    // tier -> last sync time

	// Frame tracking
	frame uint64

	// Callbacks
	onPlayerJoined func(playerID int)
	onPlayerLeft   func(playerID int)
	onGameOver     func(results map[string]interface{})

	// JavaScript interop (WASM only)
	jsSendSignal func(signalType string, toPlayerID int, data interface{}) error
	jsGetSignals func(forPlayerID int) ([]Signal, error)

	// Sync intervals
	fastInterval     time.Duration // ~16-33ms (30-60/sec)
	moderateInterval time.Duration // ~66ms (15/sec)
	slowInterval     time.Duration // ~200ms (5/sec)
}

// SyncedTable tracks a table that should be synchronized
type SyncedTable struct {
	Path      string
	Tier      SyncTier
	LastState map[string]interface{}
	LuaTable  interface{} // Reference to Lua table (for WASM, we'll use a different approach)
}

// Signal represents a WebRTC signaling message
type Signal struct {
	Type      string      `json:"type"` // "offer", "answer", "ice-candidate"
	Data      interface{} `json:"data"`
	FromID    int         `json:"from_id"`
	ToID      int         `json:"to_id"`
	Timestamp float64     `json:"timestamp"`
}

// NewNetworkManager creates a new network manager
func NewNetworkManager() *NetworkManager {
	return &NetworkManager{
		connections:      make(map[int]Connection),
		playerInputs:     make(map[int]map[int]bool),
		syncedTables:     make(map[string]*SyncedTable),
		lastSyncTime:     make(map[string]time.Time),
		isMultiplayer:    false,
		isHost:           false,
		playerID:         1,
		playerCount:      1,
		frame:            0,
		fastInterval:     time.Millisecond * 16,  // ~60/sec
		moderateInterval: time.Millisecond * 66,  // ~15/sec
		slowInterval:     time.Millisecond * 200, // ~5/sec
	}
}

// InitializeMultiplayer initializes multiplayer mode
// gameInstanceID: Convex game instance ID
// isHost: whether this player is the host
// playerID: this player's ID (1-6)
// playerCount: total number of players
// jsSendSignal: JavaScript function to send WebRTC signals
// jsGetSignals: JavaScript function to get WebRTC signals
func (nm *NetworkManager) InitializeMultiplayer(
	gameInstanceID string,
	isHost bool,
	playerID int,
	playerCount int,
	jsSendSignal func(string, int, interface{}) error,
	jsGetSignals func(int) ([]Signal, error),
) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	nm.isMultiplayer = true
	nm.isHost = isHost
	nm.playerID = playerID
	nm.playerCount = playerCount
	nm.jsSendSignal = jsSendSignal
	nm.jsGetSignals = jsGetSignals

	// Initialize player inputs
	for i := 1; i <= playerCount; i++ {
		nm.playerInputs[i] = make(map[int]bool)
	}

	return nil
}

// IsMultiplayer returns true if in multiplayer mode
func (nm *NetworkManager) IsMultiplayer() bool {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return nm.isMultiplayer
}

// IsHost returns true if this player is the host
func (nm *NetworkManager) IsHost() bool {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return nm.isHost
}

// PlayerID returns this player's ID (1-6)
func (nm *NetworkManager) PlayerID() int {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return nm.playerID
}

// PlayerCount returns total number of players (1-6)
func (nm *NetworkManager) PlayerCount() int {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return nm.playerCount
}

// RegisterSyncedTable registers a Lua table for automatic synchronization
func (nm *NetworkManager) RegisterSyncedTable(tablePath string, tier SyncTier, initialState map[string]interface{}) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if !nm.isMultiplayer {
		return fmt.Errorf("not in multiplayer mode")
	}

	nm.syncedTables[tablePath] = &SyncedTable{
		Path:      tablePath,
		Tier:      tier,
		LastState: make(map[string]interface{}),
	}

	// Deep copy initial state
	for k, v := range initialState {
		nm.syncedTables[tablePath].LastState[k] = v
	}

	return nil
}

// UnregisterSyncedTable removes a table from synchronization
func (nm *NetworkManager) UnregisterSyncedTable(tablePath string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	delete(nm.syncedTables, tablePath)
}

// UpdateSyncedTable updates the state of a synced table
// This should be called from Lua when the table changes
func (nm *NetworkManager) UpdateSyncedTable(tablePath string, newState map[string]interface{}) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	synced, ok := nm.syncedTables[tablePath]
	if !ok {
		return // Table not registered for sync
	}

	// Store new state
	synced.LastState = make(map[string]interface{})
	for k, v := range newState {
		synced.LastState[k] = v
	}
}

// SetPlayerInput sets input state for a player (host only, or local player)
func (nm *NetworkManager) SetPlayerInput(playerID int, buttonID int, pressed bool) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if !nm.isMultiplayer {
		return
	}

	// Only host can set inputs for other players
	// Non-hosts can only set their own input
	if !nm.isHost && playerID != nm.playerID {
		return
	}

	if nm.playerInputs[playerID] == nil {
		nm.playerInputs[playerID] = make(map[int]bool)
	}
	nm.playerInputs[playerID][buttonID] = pressed
}

// GetPlayerInput gets input state for a player (host only)
func (nm *NetworkManager) GetPlayerInput(playerID int, buttonID int) bool {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	if !nm.isHost {
		return false
	}

	inputs, ok := nm.playerInputs[playerID]
	if !ok {
		return false
	}
	return inputs[buttonID]
}

// SendInput sends local player's input to host (non-host only)
func (nm *NetworkManager) SendInput(buttons map[int]bool) error {
	nm.mu.RLock()
	isHost := nm.isHost
	isMultiplayer := nm.isMultiplayer
	playerID := nm.playerID
	frame := nm.frame
	nm.mu.RUnlock()

	if !isMultiplayer || isHost {
		return nil // Host doesn't send inputs, or not multiplayer
	}

	// Convert button map
	buttonMap := make(map[string]bool)
	for btnID, pressed := range buttons {
		buttonMap[fmt.Sprintf("%d", btnID)] = pressed
	}

	packet := InputPacket{
		Type:      PacketTypeInput,
		PlayerID:  playerID,
		Frame:     frame,
		Buttons:   buttonMap,
		Timestamp: float64(time.Now().UnixNano()) / 1e9,
	}

	data, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	// Send to host via connection (will be implemented in WASM layer)
	// For now, this is a placeholder
	return nm.sendToHost(data)
}

// sendToHost sends data to host (non-host only)
func (nm *NetworkManager) sendToHost(data []byte) error {
	// This will be implemented in WASM layer via WebRTC
	// For now, just a placeholder
	return nil
}

// ProcessReceivedInput processes an input packet from a player (host only)
func (nm *NetworkManager) ProcessReceivedInput(packet InputPacket) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if !nm.isHost {
		return
	}

	// Update player input state
	if nm.playerInputs[packet.PlayerID] == nil {
		nm.playerInputs[packet.PlayerID] = make(map[int]bool)
	}

	for btnStr, pressed := range packet.Buttons {
		var btnID int
		fmt.Sscanf(btnStr, "%d", &btnID)
		nm.playerInputs[packet.PlayerID][btnID] = pressed
	}
}

// UpdateFrame increments frame counter and processes synchronization
func (nm *NetworkManager) UpdateFrame(dt time.Duration) {
	// Fast path: completely skip all work if not in multiplayer mode
	nm.mu.RLock()
	isMultiplayer := nm.isMultiplayer
	nm.mu.RUnlock()

	if !isMultiplayer {
		return // Skip all work in solo mode for maximum performance
	}

	nm.mu.Lock()
	nm.frame++
	frame := nm.frame
	isHost := nm.isHost
	syncedTables := make(map[string]*SyncedTable)
	for k, v := range nm.syncedTables {
		syncedTables[k] = v
	}
	lastSyncTime := make(map[string]time.Time)
	for k, v := range nm.lastSyncTime {
		lastSyncTime[k] = v
	}
	nm.mu.Unlock()

	if !isHost {
		return // Only host syncs state
	}

	now := time.Now()

	// Check each sync tier and send updates if needed
	for tablePath, synced := range syncedTables {
		var interval time.Duration
		switch synced.Tier {
		case SyncTierFast:
			interval = nm.fastInterval
		case SyncTierModerate:
			interval = nm.moderateInterval
		case SyncTierSlow:
			interval = nm.slowInterval
		}

		lastTime, ok := lastSyncTime[string(synced.Tier)]
		if !ok || now.Sub(lastTime) >= interval {
			nm.sendStateDelta(tablePath, synced, frame)
			nm.mu.Lock()
			nm.lastSyncTime[string(synced.Tier)] = now
			nm.mu.Unlock()
		}
	}
}

// sendStateDelta sends state delta to all connected players
func (nm *NetworkManager) sendStateDelta(tablePath string, synced *SyncedTable, frame uint64) {
	// Generate delta (for now, send full state; optimize later)
	delta := StateDeltaPacket{
		Type:      PacketTypeStateDelta,
		Frame:     frame,
		Tier:      synced.Tier,
		Changes:   synced.LastState,
		Timestamp: float64(time.Now().UnixNano()) / 1e9,
	}

	data, err := json.Marshal(delta)
	if err != nil {
		return
	}

	// Send to all connections (implemented in WASM layer)
	nm.broadcastToAll(data)
}

// broadcastToAll broadcasts data to all connected players
func (nm *NetworkManager) broadcastToAll(data []byte) {
	// This will be implemented in WASM layer
	// For now, just a placeholder
}

// ProcessReceivedState processes a state delta from host (non-host only)
func (nm *NetworkManager) ProcessReceivedState(packet StateDeltaPacket) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if nm.isHost {
		return // Host doesn't receive state
	}

	// Update synced table state
	// This will be passed to Lua to update the actual table
	// For now, just store it
}

// Close closes all connections
func (nm *NetworkManager) Close() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for _, conn := range nm.connections {
		_ = conn.Close()
	}
	nm.connections = make(map[int]Connection)
	nm.isMultiplayer = false
}
