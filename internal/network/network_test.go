package network

import (
	"testing"
	"time"
)

func TestNetworkManager_InitializeMultiplayer(t *testing.T) {
	nm := NewNetworkManager()

	// Test initialization
	err := nm.InitializeMultiplayer(
		"game-123",
		true, // isHost
		1,    // playerID
		3,    // playerCount
		nil,  // jsSendSignal
		nil,  // jsGetSignals
	)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	if !nm.IsMultiplayer() {
		t.Error("Expected IsMultiplayer() to return true")
	}
	if !nm.IsHost() {
		t.Error("Expected IsHost() to return true")
	}
	if nm.PlayerID() != 1 {
		t.Errorf("Expected PlayerID() to return 1, got %d", nm.PlayerID())
	}
	if nm.PlayerCount() != 3 {
		t.Errorf("Expected PlayerCount() to return 3, got %d", nm.PlayerCount())
	}
}

func TestNetworkManager_RegisterSyncedTable(t *testing.T) {
	nm := NewNetworkManager()
	err := nm.InitializeMultiplayer("game-123", true, 1, 2, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	initialState := map[string]interface{}{
		"x":  100.0,
		"y":  200.0,
		"vx": 5.0,
	}

	err = nm.RegisterSyncedTable("players.1", SyncTierFast, initialState)
	if err != nil {
		t.Fatalf("RegisterSyncedTable failed: %v", err)
	}

	// Try to register when not in multiplayer
	nm2 := NewNetworkManager()
	err = nm2.RegisterSyncedTable("players.1", SyncTierFast, initialState)
	if err == nil {
		t.Error("Expected error when registering in non-multiplayer mode")
	}
}

func TestNetworkManager_SetGetPlayerInput(t *testing.T) {
	nm := NewNetworkManager()
	err := nm.InitializeMultiplayer("game-123", true, 1, 3, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	// Host can set inputs for any player
	nm.SetPlayerInput(1, 0, true) // Player 1, button 0
	nm.SetPlayerInput(2, 0, true) // Player 2, button 0
	nm.SetPlayerInput(3, 1, true) // Player 3, button 1

	if !nm.GetPlayerInput(1, 0) {
		t.Error("Expected Player 1 button 0 to be true")
	}
	if !nm.GetPlayerInput(2, 0) {
		t.Error("Expected Player 2 button 0 to be true")
	}
	if !nm.GetPlayerInput(3, 1) {
		t.Error("Expected Player 3 button 1 to be true")
	}
	if nm.GetPlayerInput(1, 1) {
		t.Error("Expected Player 1 button 1 to be false")
	}

	// Non-host cannot get other players' inputs
	nm2 := NewNetworkManager()
	err = nm2.InitializeMultiplayer("game-123", false, 2, 3, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	if nm2.GetPlayerInput(1, 0) {
		t.Error("Non-host should not be able to get other players' inputs")
	}
}

func TestNetworkManager_ProcessReceivedInput(t *testing.T) {
	nm := NewNetworkManager()
	err := nm.InitializeMultiplayer("game-123", true, 1, 2, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	// Simulate input packet from player 2
	packet := InputPacket{
		Type:     PacketTypeInput,
		PlayerID: 2,
		Frame:    100,
		Buttons: map[string]bool{
			"0": true, // left
			"4": true, // jump
		},
		Timestamp: float64(time.Now().UnixNano()) / 1e9,
	}

	nm.ProcessReceivedInput(packet)

	// Verify input was recorded
	if !nm.GetPlayerInput(2, 0) {
		t.Error("Expected Player 2 button 0 (left) to be true")
	}
	if !nm.GetPlayerInput(2, 4) {
		t.Error("Expected Player 2 button 4 (jump) to be true")
	}
}

func TestNetworkManager_UpdateFrame(t *testing.T) {
	nm := NewNetworkManager()
	err := nm.InitializeMultiplayer("game-123", true, 1, 2, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	initialState := map[string]interface{}{
		"x": 100.0,
		"y": 200.0,
	}
	err = nm.RegisterSyncedTable("players.1", SyncTierSlow, initialState)
	if err != nil {
		t.Fatalf("RegisterSyncedTable failed: %v", err)
	}

	// Update frame multiple times (should trigger sync based on interval)
	for i := 0; i < 10; i++ {
		nm.UpdateFrame(time.Millisecond * 200) // 200ms per frame (slow tier)
		time.Sleep(time.Millisecond * 10)      // Small delay
	}
}

func TestSyncTier(t *testing.T) {
	// Test sync tier constants
	if SyncTierFast != "fast" {
		t.Errorf("Expected SyncTierFast to be 'fast', got '%s'", SyncTierFast)
	}
	if SyncTierModerate != "moderate" {
		t.Errorf("Expected SyncTierModerate to be 'moderate', got '%s'", SyncTierModerate)
	}
	if SyncTierSlow != "slow" {
		t.Errorf("Expected SyncTierSlow to be 'slow', got '%s'", SyncTierSlow)
	}
}

func TestNetworkManager_UnregisterSyncedTable(t *testing.T) {
	nm := NewNetworkManager()
	err := nm.InitializeMultiplayer("game-123", true, 1, 2, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	initialState := map[string]interface{}{
		"x": 100.0,
	}
	err = nm.RegisterSyncedTable("players.1", SyncTierFast, initialState)
	if err != nil {
		t.Fatalf("RegisterSyncedTable failed: %v", err)
	}

	nm.UnregisterSyncedTable("players.1")

	// Try to register again (should work)
	err = nm.RegisterSyncedTable("players.1", SyncTierFast, initialState)
	if err != nil {
		t.Errorf("Expected to be able to register table again, got error: %v", err)
	}
}

func TestNetworkManager_Close(t *testing.T) {
	nm := NewNetworkManager()
	err := nm.InitializeMultiplayer("game-123", true, 1, 2, nil, nil)
	if err != nil {
		t.Fatalf("InitializeMultiplayer failed: %v", err)
	}

	if !nm.IsMultiplayer() {
		t.Error("Expected to be in multiplayer mode")
	}

	nm.Close()

	if nm.IsMultiplayer() {
		t.Error("Expected to not be in multiplayer mode after Close()")
	}
}

func TestPacketTypes(t *testing.T) {
	// Test packet type constants
	if PacketTypeInput != "input" {
		t.Errorf("Expected PacketTypeInput to be 'input', got '%s'", PacketTypeInput)
	}
	if PacketTypeStateDelta != "state_delta" {
		t.Errorf("Expected PacketTypeStateDelta to be 'state_delta', got '%s'", PacketTypeStateDelta)
	}
	if PacketTypeFullState != "full_state" {
		t.Errorf("Expected PacketTypeFullState to be 'full_state', got '%s'", PacketTypeFullState)
	}
}
