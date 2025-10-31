# RetroForge Multiplayer Game Engine - Complete Design Document

**Version:** 2.0 (Multiplayer Edition)  
**Date:** October 31, 2025  
**Status:** FINAL - Ready for Implementation

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [System Architecture](#system-architecture)
3. [Platform Roles](#platform-roles)
4. [Network Architecture](#network-architecture)
5. [Multiplayer API Specification](#multiplayer-api-specification)
6. [Data Models](#data-models)
7. [Game Flow](#game-flow)
8. [Connection Management](#connection-management)
9. [Network Protocol](#network-protocol)
10. [Example Implementation](#example-implementation)
11. [Security & Performance](#security--performance)
12. [Development & Testing](#development--testing)

---

## Executive Summary

### Project Vision

**RetroForge** is a fantasy console (like PICO-8) with built-in multiplayer support, enabling developers to create retro-style 2D games that can be played solo or with up to 6 players online.

### Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│         Next.js 16 Web Application                      │
│         (Matchmaking Platform Only)                     │
│  - Browse game library                                  │
│  - Create/join lobbies                                  │
│  - Player profiles & leaderboards                       │
│  - WebRTC signaling server (via Convex)                │
└─────────────┬───────────────────────────────────────────┘
              │
              │ Click "Play" → Loads cart
              ↓
┌─────────────────────────────────────────────────────────┐
│    RetroForge Engine (Go → WASM, ~2MB)                  │
│    (Runs entirely in browser)                           │
│                                                          │
│  ┌────────────────────────────────────────────────┐   │
│  │     Engine Core                                │   │
│  │  - Lua VM (gopher-lua)                         │   │
│  │  - Graphics (Canvas/WebGL)                     │   │
│  │  - Audio (chip-tune synthesis)                 │   │
│  │  - Input (keyboard/gamepad/touch)              │   │
│  │  - Physics (Box2D)                             │   │
│  │  - Node system (Godot-style)                   │   │
│  │  - WebRTC networking (BUILT-IN)                │   │
│  └────────────────────────────────────────────────┘   │
│                                                          │
│  ┌────────────────────────────────────────────────┐   │
│  │     Game Cart (.rfs file, ~100KB)             │   │
│  │  - game.lua (16,384 tokens max)                │   │
│  │  - sprites.json (256 slots)                    │   │
│  │  - map.json (8 layers)                         │   │
│  │  - sfx.json (sound effects)                    │   │
│  │  - music.json (music tracks)                   │   │
│  │  - manifest.json (metadata)                    │   │
│  └────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Clear Separation:** Next.js/Convex = platform, RetroForge = game engine
2. **Developer Simplicity:** Minimal networking API, automatic synchronization
3. **Bandwidth Efficiency:** 3-tier sync system (fast/moderate/slow)
4. **Host Authority:** One player controls game logic, prevents conflicts
5. **PICO-8 Familiarity:** Developers write normal game code, networking "just works"

---

## System Architecture

### Component Responsibilities

#### Next.js/Convex (Matchmaking Platform)

**What it DOES handle:**
- User authentication (Clerk)
- Game library/browser
- Lobby creation and management
- Player profiles and statistics
- Leaderboards
- WebRTC signaling (offer/answer/ICE exchange)
- Match history storage

**What it DOES NOT handle:**
- Game rendering
- Game logic
- Physics simulation
- Player input processing
- Game state synchronization
- Audio/visual effects

**Think of it as:** Steam or Epic Games Store - a platform to find and launch games.

---

#### RetroForge Engine (Go → WASM)

**What it DOES handle:**
- Complete game runtime
- All rendering (480×270 resolution)
- Physics simulation (Box2D)
- Audio synthesis (5 waveforms, 8 channels)
- Input management
- Lua code execution
- Node system (Godot-style)
- **WebRTC networking (built-in)**
- **Automatic state synchronization**
- **Connection management**

**What it DOES NOT handle:**
- Matchmaking
- Lobby management
- User profiles
- Persistent storage (handled by Convex)

**Think of it as:** Unity/Godot - a complete game engine that happens to run in the browser.

---

### Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Platform** | Next.js 16 (App Router) | Web app framework |
| **Backend** | Convex | Real-time serverless backend |
| **Auth** | Clerk | User authentication |
| **Engine** | Go 1.23 + TinyGo | Compiles to WASM |
| **Scripting** | gopher-lua | Lua VM in Go |
| **Physics** | Box2D-go | Rigid body physics |
| **Graphics** | Canvas 2D / WebGL | Browser rendering |
| **Networking** | WebRTC (pion/webrtc) | P2P data channels |
| **UI** | React 19 + Tailwind | Modern UI |

---

## Platform Roles

### Next.js/Convex Platform

#### Features

**Game Library**
- Browse available RetroForge carts
- Filter by genre, multiplayer support, player count
- Search by title, author, tags
- Featured games, trending, new releases

**Lobby System**
- Create lobby (select cart, set player limit)
- Browse active lobbies
- Join lobby (if not full, not started)
- Lobby chat (optional)
- Ready/unready status
- Host can start game when ready

**User Management**
- Profile creation (username, avatar)
- Match history
- Win/loss statistics
- Leaderboards (global, per-game)
- Friend system (optional)

**WebRTC Signaling**
- Exchange WebRTC offers/answers between players
- ICE candidate relay
- Connection status tracking
- Automatic cleanup on disconnect

---

### RetroForge Engine

#### Core Systems

**Lua Runtime**
- 16,384 token limit per cart
- PICO-8-compatible API
- Extended with node system
- Automatic garbage collection

**Graphics System**
- 480×270 resolution (16:9)
- 50-color palettes
- Sprite rendering (4×4 to 32×32)
- Tilemap support (8 layers)
- Shape primitives
- Camera system

**Audio System**
- 8 simultaneous channels
- 5 waveforms (square, triangle, sawtooth, sine, noise)
- ADSR envelope
- Effects (vibrato, arpeggio)
- 3-tier system (SoundManager, AudioPlayer, MusicPlayer)

**Physics System**
- Box2D rigid body simulation
- Static, dynamic, kinematic bodies
- Circle and box shapes
- Collision detection
- Forces and impulses

**Input System**
- Keyboard support (16 buttons)
- Gamepad support (standard mapping)
- Touch support (mobile)
- Button press/release detection

**Node System**
- Godot-style scene graph
- 15+ node types
- Parent-child relationships
- Signals and callbacks

**Network System** (NEW!)
- WebRTC data channels
- Automatic state synchronization
- Host authority model
- 3-tier sync frequencies
- Connection lifecycle management

---

## Network Architecture

### Topology: Star (Hub and Spoke)

```
        Player 2
            ↑
            |
Player 3 ← Host (Player 1) → Player 4
            |
            ↓
        Player 5
            |
            ↓
        Player 6

Total Connections: 5 (host to each player)
```

**Rationale:**
- Matches host authority model
- Simpler than full mesh (15 connections)
- Less bandwidth than mesh
- Host has best connection quality
- Single source of truth

**Connection Flow:**
1. All players connect ONLY to host via WebRTC
2. Host receives inputs from all players
3. Host runs game logic (60 FPS)
4. Host broadcasts state updates to all players
5. Players render received state

---

### Authority Model: Host Authority

**Host Responsibilities:**
- Runs complete game logic in `_update()`
- Receives all player inputs (5/sec)
- Applies inputs to game state
- Broadcasts state changes to all players
- Handles shared game objects (enemies, powerups, etc.)

**Non-Host Responsibilities:**
- Send inputs to host (5/sec)
- Receive state updates from host
- Render game state in `_draw()`
- Optional: Run prediction for local player (v2 feature)

**Why Host Authority:**
- Simple to implement and understand
- No conflicts (single source of truth)
- Good for casual 2-6 player games
- Prevents most cheating (host validates)
- Matches RetroForge's simplicity goal

**Trade-offs:**
- Host has 0ms latency advantage
- Non-hosts have ~200ms input latency (acceptable for casual games)
- If host disconnects, game ends

---

### Sync System: 3-Tier Frequencies

Developers register tables for automatic synchronization with one of three priority levels:

#### Fast Tier (30-60/sec)
**Use for:** Visually smooth movement and fast-changing data
```lua
rf.network_sync(players, "fast")      -- player positions
rf.network_sync(bullets, "fast")      -- projectiles
rf.network_sync(particles, "fast")    -- visual effects
```

#### Moderate Tier (15/sec)
**Use for:** Important but not visually critical data
```lua
rf.network_sync(powerups, "moderate") -- pickup states
rf.network_sync(enemies, "moderate")  -- NPC positions
rf.network_sync(animations, "moderate") -- animation states
```

#### Slow Tier (5/sec)
**Use for:** Rarely changing data, numbers, UI
```lua
rf.network_sync(score, "slow")        -- score updates
rf.network_sync(level, "slow")        -- level number
rf.network_sync(timer, "slow")        -- countdown timer
rf.network_sync(health, "slow")       -- health bars
```

**Input Frequency:** 5/sec (slow tier)
- Humans cannot react faster than ~200ms
- Holding a button doesn't need 60 updates/sec
- Frees bandwidth for visual state updates

**Bandwidth Calculation:**
```
Fast tier:   1KB × 60/sec = 60 KB/sec
Moderate:    0.5KB × 15/sec = 7.5 KB/sec
Slow:        0.1KB × 5/sec = 0.5 KB/sec
Player input: 0.05KB × 5/sec × 5 players = 1.25 KB/sec

Total per connection: ~70 KB/sec
Host with 5 connections: ~350 KB/sec upload (2.8 Mbps)
```

---

## Multiplayer API Specification

### Core Network Functions

```lua
-- Connection Information
rf.is_multiplayer()     → boolean
  -- Returns true if game is in multiplayer mode
  -- Returns false for solo play

rf.player_count()       → number (1-6)
  -- Returns total number of connected players
  -- Always >= 1 (includes local player)

rf.my_player_id()       → number (1-6)
  -- Returns the local player's ID
  -- Host is always player_id = 1

rf.is_host()           → boolean
  -- Returns true if local player is the host
  -- Only host runs game logic in _update()
```

### State Synchronization

```lua
-- Register table for automatic synchronization
rf.network_sync(table, frequency)
  -- table: Lua table to synchronize
  -- frequency: "fast" | "moderate" | "slow"
  -- Example: rf.network_sync(players, "fast")

-- Unregister table from synchronization
rf.network_unsync(table)
  -- Stops syncing the table
  -- Useful for cleanup or dynamic objects
```

**Sync Behavior:**
- Engine automatically detects changes to registered tables
- Only sends deltas (changed values), not full table
- Enforces ownership rules (see below)
- Batches multiple changes into single packet

**Ownership Rules:**
- Tables keyed by player_id: only that player can modify their entry
  ```lua
  players[1].x = 100  -- Only Player 1 can do this
  players[2].x = 100  -- Only Player 2 can do this
  players[3].x = 100  -- Player 1 CANNOT modify this (ignored)
  ```
- Other tables: only host can modify
  ```lua
  enemies[1].x = 100  -- Only host can modify
  powerups[1].taken = true  -- Only host can modify
  ```

### Input Handling (Host Only)

```lua
-- Check another player's button state (host only)
rf.btn(player_id, button) → boolean
  -- player_id: 1-6
  -- button: 0-15
  -- Returns true if that player is pressing the button
  -- Only works in host's code
  -- Non-hosts use normal btn() for local player

-- Example (host code):
if rf.btn(2, 0) then  -- Check if Player 2 is pressing left
  players[2].vx = -2
end
```

### Network Statistics (Debug)

```lua
-- Only available in development mode
rf.network_ping(player_id) → number
  -- Returns latency to player in milliseconds
  -- Only works in dev mode

rf.network_bandwidth() → number
  -- Returns current bandwidth usage in KB/sec
  -- Only works in dev mode

rf.network_packet_loss() → number
  -- Returns packet loss percentage (0.0-1.0)
  -- Only works in dev mode
```

### Debug Tools (Development Only)

```lua
-- Simulate network conditions
rf.debug_set_latency(ms)
  -- Add artificial latency (milliseconds)
  -- Example: rf.debug_set_latency(100)

rf.debug_set_packet_loss(percentage)
  -- Simulate packet loss (0.0-1.0)
  -- Example: rf.debug_set_packet_loss(0.1)  -- 10% loss
```

---

## Data Models

### Convex Schema

```typescript
// convex/schema.ts
import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  // User profiles
  users: defineTable({
    clerkId: v.string(),
    username: v.string(),
    displayName: v.string(),
    avatar: v.optional(v.string()),
    wins: v.number(),
    losses: v.number(),
    gamesPlayed: v.number(),
    createdAt: v.number(),
  })
    .index("by_clerk_id", ["clerkId"])
    .index("by_username", ["username"]),

  // Game carts in library
  carts: defineTable({
    title: v.string(),
    author: v.id("users"),
    description: v.string(),
    version: v.string(),
    cartUrl: v.string(),  // Storage URL for .rfs file
    thumbnailUrl: v.optional(v.string()),
    
    // Multiplayer metadata (from manifest.json)
    multiplayerEnabled: v.boolean(),
    minPlayers: v.optional(v.number()),  // 2-6
    maxPlayers: v.optional(v.number()),  // 2-6
    supportsSolo: v.boolean(),
    
    // Stats
    plays: v.number(),
    favorites: v.number(),
    
    tags: v.array(v.string()),
    createdAt: v.number(),
    updatedAt: v.number(),
  })
    .index("by_author", ["author"])
    .index("by_multiplayer", ["multiplayerEnabled"]),

  // Game lobbies
  lobbies: defineTable({
    hostId: v.id("users"),
    cartId: v.id("carts"),
    name: v.string(),
    
    // Current players in lobby
    players: v.array(v.object({
      userId: v.id("users"),
      username: v.string(),
      isReady: v.boolean(),
    })),
    
    maxPlayers: v.number(),  // From cart manifest
    status: v.union(
      v.literal("waiting"),    // Accepting players
      v.literal("starting"),   // Countdown/setup
      v.literal("in_progress"), // Game active (locked)
      v.literal("completed")   // Game ended
    ),
    
    createdAt: v.number(),
    startedAt: v.optional(v.number()),
  })
    .index("by_status", ["status"])
    .index("by_cart", ["cartId"])
    .index("by_host", ["hostId"]),

  // Active game instances
  gameInstances: defineTable({
    lobbyId: v.id("lobbies"),
    cartId: v.id("carts"),
    
    // Players in game
    players: v.array(v.object({
      userId: v.id("users"),
      playerId: v.number(),  // 1-6 (in-game ID)
      isHost: v.boolean(),
    })),
    
    status: v.union(
      v.literal("initializing"), // WebRTC connecting
      v.literal("running"),      // Game active
      v.literal("ended")         // Game completed
    ),
    
    createdAt: v.number(),
    endedAt: v.optional(v.number()),
  })
    .index("by_lobby", ["lobbyId"])
    .index("by_status", ["status"]),

  // WebRTC signaling messages
  webrtcSignals: defineTable({
    gameInstanceId: v.id("gameInstances"),
    fromPlayerId: v.number(),  // 1-6
    toPlayerId: v.number(),    // 1-6 (always host for star topology)
    
    signalType: v.union(
      v.literal("offer"),
      v.literal("answer"),
      v.literal("ice-candidate")
    ),
    
    signalData: v.any(),  // SDP or ICE candidate JSON
    processed: v.boolean(),
    createdAt: v.number(),
  })
    .index("by_game_and_receiver", ["gameInstanceId", "toPlayerId", "processed"])
    .index("by_game", ["gameInstanceId"]),

  // Match results and statistics
  matchResults: defineTable({
    gameInstanceId: v.id("gameInstances"),
    cartId: v.id("carts"),
    
    // Per-player results
    players: v.array(v.object({
      userId: v.id("users"),
      playerId: v.number(),
      score: v.number(),
      placement: v.number(),  // 1st, 2nd, 3rd, etc.
      customStats: v.any(),   // Game-specific stats
    })),
    
    duration: v.number(),  // milliseconds
    createdAt: v.number(),
  })
    .index("by_game", ["gameInstanceId"])
    .index("by_cart", ["cartId"]),
});
```

---

### Cart Manifest (manifest.json)

Every RetroForge cart includes a `manifest.json` file with metadata:

```json
{
  "version": "1.0",
  "title": "Space Battle Royale",
  "author": "CoolDev",
  "description": "Fast-paced space combat with power-ups",
  "thumbnail": "thumbnail.png",
  
  "multiplayer": {
    "enabled": true,
    "minPlayers": 2,
    "maxPlayers": 6,
    "supportsSolo": false,
    "description": "Team up or battle royale mode"
  },
  
  "tags": ["action", "arcade", "space", "competitive"],
  "cartVersion": "1.2.0",
  "engineVersion": "2.0.0"
}
```

**Field Descriptions:**

- `multiplayer.enabled`: If true, cart can be played in multiplayer
- `multiplayer.minPlayers`: Minimum players needed to start (2-6)
- `multiplayer.maxPlayers`: Maximum players allowed (2-6)
- `multiplayer.supportsSolo`: If true, can also be played alone
- `tags`: Categories for filtering in game library

**Next.js uses this to:**
- Show "Multiplayer" badge
- Filter multiplayer-capable games
- Enforce player limits in lobbies
- Display player count range

---

## Game Flow

### Complete Player Journey

```
┌─────────────────────────────────────────────────────────┐
│ 1. AUTHENTICATION                                        │
│    User opens app → Clerk login → Profile loaded        │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 2. GAME LIBRARY                                          │
│    Browse carts → Filter by multiplayer → Select cart   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 3. GAME MODE SELECTION                                   │
│    If multiplayer enabled:                               │
│      Option A: Play Solo (if supportsSolo = true)       │
│      Option B: Create Lobby                             │
│      Option C: Join Existing Lobby                      │
│    If multiplayer not enabled:                           │
│      Launch directly to game                            │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 4. LOBBY PHASE (Multiplayer Only)                       │
│    - Real-time player list via Convex subscription      │
│    - Players mark ready                                 │
│    - Host clicks "Start Game" when ready                │
│    - Convex locks lobby (status: "starting")            │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 5. GAME INITIALIZATION                                   │
│    A. Convex creates gameInstance                       │
│    B. Convex sets lobby status: "in_progress"           │
│    C. All browsers load cart .rfs file                  │
│    D. RetroForge WASM engine initializes                │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 6. WEBRTC CONNECTION ESTABLISHMENT                       │
│    Star Topology Setup:                                  │
│                                                          │
│    For each non-host player:                            │
│      1. Create RTCPeerConnection to host                │
│      2. Create data channel                             │
│      3. Generate SDP offer                              │
│      4. Send offer to Convex                            │
│                                                          │
│    Host:                                                 │
│      1. Subscribe to incoming offers from Convex        │
│      2. For each offer, create RTCPeerConnection        │
│      3. Generate SDP answer                             │
│      4. Send answer to Convex                           │
│                                                          │
│    All players:                                          │
│      1. Exchange ICE candidates via Convex              │
│      2. Establish direct connections                    │
│      3. Retry up to 3 times (max 3 sec each)           │
│      4. If fails, kick player from game                 │
│                                                          │
│    Result: 5 direct WebRTC connections (host ↔ players) │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 7. INITIAL STATE SYNC                                    │
│    A. Host's Lua _init() executes                       │
│    B. Host creates game state:                          │
│         players = {}                                    │
│         rf.network_sync(players, "fast")                │
│         enemies = {}                                    │
│         rf.network_sync(enemies, "moderate")            │
│    C. Engine serializes all sync'd tables               │
│    D. Host sends full state snapshot to all players     │
│    E. Non-hosts' engines populate their tables          │
│    F. All players ready to start                        │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 8. ACTIVE GAMEPLAY                                       │
│                                                          │
│    HOST (60 FPS):                                        │
│      _update() {                                        │
│        1. Receive inputs from all players (5/sec)       │
│        2. Apply inputs to game state                    │
│        3. Run game logic (physics, AI, collisions)      │
│        4. Engine auto-syncs registered tables:          │
│           - Fast tier: 30-60/sec                        │
│           - Moderate: 15/sec                            │
│           - Slow: 5/sec                                 │
│      }                                                   │
│      _draw() {                                          │
│        Render local game state                          │
│      }                                                   │
│                                                          │
│    NON-HOST PLAYERS (60 FPS):                           │
│      _update() {                                        │
│        1. Capture local input                           │
│        2. Engine sends input to host (5/sec)            │
│        3. Engine receives state updates from host       │
│        4. Engine updates local copy of game state       │
│      }                                                   │
│      _draw() {                                          │
│        Render received game state                       │
│      }                                                   │
│                                                          │
│    Loose Sync: Host doesn't wait for slow players       │
│    Lagging players get older state but game continues   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 9. DISCONNECTION HANDLING                                │
│                                                          │
│    Non-host disconnects:                                 │
│      - Remove player from game immediately              │
│      - Notify remaining players                         │
│      - Game continues with remaining players            │
│                                                          │
│    Host disconnects:                                     │
│      - Game ends for all players                        │
│      - Return all to lobby/menu                         │
│      - Display "Host disconnected" message              │
│      - No host migration (too complex for v1)           │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 10. GAME END                                             │
│    A. Win condition met (cart-specific)                 │
│    B. Host detects game over                            │
│    C. Host calculates final scores/placements           │
│    D. Host calls: rf.game_over(results)                 │
│    E. Engine sends results to Convex                    │
│    F. Convex stores match results                       │
│    G. Convex updates player statistics                  │
│    H. All players see results screen                    │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 11. POST-GAME                                            │
│    - Display results (scores, placements, stats)        │
│    - Show MVP, achievements, etc.                       │
│    - Options:                                           │
│        • Return to Menu                                 │
│        • Browse Games                                   │
│        • Create New Lobby (for rematch)                 │
│    - Lobby is closed (must create new for rematch)      │
└─────────────────────────────────────────────────────────┘
```

---

### State Transitions

**Lobby States:**
```
waiting → starting → in_progress → completed
```

**Game Instance States:**
```
initializing → running → ended
```

**Connection States:**
```
connecting → connected → disconnected
```

---

## Connection Management

### WebRTC Connection Flow (Detailed)

#### Star Topology - Host as Central Hub

```
Non-Host Player Flow:
  1. Create RTCPeerConnection
  2. Create data channel "retroforge-game"
  3. Set up event handlers:
     - onicecandidate → send to Convex
     - ondatachannel → (N/A for initiator)
     - onconnectionstatechange → monitor
  4. Create SDP offer
  5. setLocalDescription(offer)
  6. Send offer to Convex (signaling server)
  7. Wait for SDP answer from host via Convex
  8. setRemoteDescription(answer)
  9. Exchange ICE candidates via Convex
  10. Data channel opens → ready to play

Host Flow:
  1. Subscribe to incoming offers via Convex
  2. For each offer received:
     a. Create RTCPeerConnection
     b. Set up event handlers
     c. setRemoteDescription(offer)
     d. Create SDP answer
     e. setLocalDescription(answer)
     f. Send answer to Convex
     g. ondatachannel event → store data channel
     h. Exchange ICE candidates via Convex
     i. Data channel opens → player connected
  3. Once all N-1 connections open → start game
```

#### Connection Retry Logic

```javascript
// Retry parameters
const MAX_RETRIES = 3;
const RETRY_TIMEOUT = 3000; // 3 seconds

async function connectToHost(hostId, gameId) {
  for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
    try {
      const connection = await establishWebRTC(hostId, gameId);
      
      // Wait for connection with timeout
      await waitForConnection(connection, RETRY_TIMEOUT);
      
      return connection; // Success!
      
    } catch (error) {
      console.log(`Connection attempt ${attempt} failed:`, error);
      
      if (attempt === MAX_RETRIES) {
        // Final attempt failed - kick player
        throw new Error("Failed to establish connection after 3 attempts");
      }
      
      // Clean up failed connection before retry
      connection.close();
    }
  }
}

// In RetroForge engine (Go):
func (nm *NetworkManager) ConnectToHost() error {
    for attempt := 1; attempt <= 3; attempt++ {
        conn, err := nm.establishWebRTC()
        if err == nil {
            return nil // Success
        }
        
        log.Printf("Attempt %d failed: %v", attempt, err)
        time.Sleep(time.Second) // Brief pause between retries
    }
    
    return fmt.Errorf("connection failed after 3 attempts")
}
```

### Convex Signaling Functions

```typescript
// convex/webrtc.ts

export const sendSignal = mutation({
  args: {
    gameInstanceId: v.id("gameInstances"),
    fromPlayerId: v.number(),
    toPlayerId: v.number(),
    signalType: v.union(
      v.literal("offer"),
      v.literal("answer"),
      v.literal("ice-candidate")
    ),
    signalData: v.any(),
  },
  handler: async (ctx, args) => {
    // Store signal for recipient to fetch
    await ctx.db.insert("webrtcSignals", {
      gameInstanceId: args.gameInstanceId,
      fromPlayerId: args.fromPlayerId,
      toPlayerId: args.toPlayerId,
      signalType: args.signalType,
      signalData: args.signalData,
      processed: false,
      createdAt: Date.now(),
    });
  },
});

export const getSignals = query({
  args: {
    gameInstanceId: v.id("gameInstances"),
    forPlayerId: v.number(),
  },
  handler: async (ctx, args) => {
    // Get unprocessed signals for this player
    const signals = await ctx.db
      .query("webrtcSignals")
      .withIndex("by_game_and_receiver", (q) =>
        q
          .eq("gameInstanceId", args.gameInstanceId)
          .eq("toPlayerId", args.forPlayerId)
          .eq("processed", false)
      )
      .collect();
    
    // Mark as processed
    await Promise.all(
      signals.map((s) => ctx.db.patch(s._id, { processed: true }))
    );
    
    return signals;
  },
});
```

---

## Network Protocol

### Packet Types

#### 1. Input Packet (Non-Host → Host, 5/sec)

```json
{
  "type": "input",
  "player_id": 2,
  "frame": 1234,
  "buttons": {
    "0": true,   // left
    "1": false,  // right
    "2": false,  // up
    "3": false,  // down
    "4": true,   // button A (jump)
    "5": false   // button B
  },
  "timestamp": 1698765432.123
}
```

**Size:** ~50 bytes  
**Frequency:** 5/sec  
**Bandwidth:** 250 bytes/sec per player

#### 2. State Delta Packet (Host → Non-Host, variable)

```json
{
  "type": "state_delta",
  "frame": 1234,
  "tier": "fast",  // or "moderate" or "slow"
  "changes": {
    "players.2.x": 150.5,
    "players.2.y": 100.0,
    "players.2.vx": 5.0,
    "players.3.x": 200.0,
    "bullets.1.x": 300.0,
    "bullets.1.y": 150.0
  },
  "timestamp": 1698765432.123
}
```

**Size:** Variable (100 bytes - 5KB depending on changes)  
**Frequency:** 
- Fast tier: 30-60/sec
- Moderate: 15/sec
- Slow: 5/sec

#### 3. Full State Snapshot (Host → Non-Host, on connect)

```json
{
  "type": "full_state",
  "frame": 0,
  "state": {
    "players": {
      "1": {"x": 100, "y": 100, "health": 100, "sprite": 1},
      "2": {"x": 200, "y": 100, "health": 100, "sprite": 1},
      "3": {"x": 300, "y": 100, "health": 100, "sprite": 1}
    },
    "enemies": [],
    "powerups": [],
    "score": {"1": 0, "2": 0, "3": 0},
    "level": 1
  }
}
```

**Size:** Variable (1-10KB)  
**Frequency:** Once on initial connection

#### 4. Player Joined (Host → All, event)

```json
{
  "type": "player_joined",
  "player_id": 4,
  "username": "NewPlayer",
  "timestamp": 1698765432.123
}
```

#### 5. Player Left (Host → All, event)

```json
{
  "type": "player_left",
  "player_id": 3,
  "reason": "disconnected",
  "timestamp": 1698765432.123
}
```

#### 6. Game Over (Host → All, event)

```json
{
  "type": "game_over",
  "results": {
    "1": {"score": 1500, "placement": 1},
    "2": {"score": 1200, "placement": 2},
    "3": {"score": 800, "placement": 3}
  },
  "duration": 180000,  // milliseconds
  "timestamp": 1698765432.123
}
```

---

### Delta Encoding Algorithm

**Efficiency Strategy:**

1. **Track Previous State**
   - Engine keeps copy of last sent state per player
   - Compare current state with previous state
   - Only send changed values

2. **Compress Keys**
   - Use short path notation: `players.2.x` instead of full JSON
   - Binary encoding for numbers (future optimization)

3. **Batch Updates**
   - Group all changes in sync tier into single packet
   - Send one packet per tier per interval
   - Don't send separate packets for each table

4. **Omit Unchanged Values**
   - If `players[2].x` hasn't changed, don't send it
   - Especially important for slow-changing data (score, health)

**Example:**

```lua
-- Previous state (frame 1233)
players[2] = {x=100, y=100, vx=5, vy=0, health=100, sprite=1}

-- Current state (frame 1234)
players[2] = {x=105, y=100, vx=5, vy=0, health=100, sprite=1}

-- Delta sent (only x changed)
{
  "changes": {
    "players.2.x": 105
  }
}
-- Saved: 80% bandwidth vs sending full object
```

---

## Example Implementation

### Example Multiplayer Cart

Complete working example of a simple multiplayer game:

```lua
-- game.lua
-- Simple multiplayer platformer example

-- Game state tables
players = {}
platforms = {}
score = {}

function _init()
  -- Create platforms (host only)
  if rf.is_host() then
    platforms = {
      {x=0, y=250, w=480, h=20},
      {x=100, y=200, w=100, h=20},
      {x=300, y=150, w=100, h=20}
    }
    rf.network_sync(platforms, "slow")  -- Platforms never move
  end
  
  -- Create player for each connected player
  local player_count = rf.player_count()
  for i = 1, player_count do
    players[i] = {
      x = 50 + (i-1) * 100,
      y = 100,
      vx = 0,
      vy = 0,
      sprite = i,
      health = 100,
      on_ground = false
    }
    score[i] = 0
  end
  
  -- Register for automatic synchronization
  rf.network_sync(players, "fast")     -- Player positions update smoothly
  rf.network_sync(score, "slow")       -- Score updates don't need high frequency
end

function _update()
  if rf.is_host() then
    -- HOST: Run game logic for all players
    update_host()
  else
    -- NON-HOST: Just send inputs (engine handles automatically)
    -- Normal btn() calls work for local player
  end
end

function update_host()
  -- Update each player based on their inputs
  for id = 1, rf.player_count() do
    local p = players[id]
    
    -- Apply inputs (5/sec, but smooth with interpolation)
    if rf.btn(id, 0) then p.vx = -3 end    -- left
    if rf.btn(id, 1) then p.vx = 3 end     -- right
    if rf.btn(id, 4) and p.on_ground then  -- jump
      p.vy = -10
      p.on_ground = false
    end
    
    -- Apply physics
    p.vy = p.vy + 0.5  -- gravity
    p.x = p.x + p.vx
    p.y = p.y + p.vy
    
    -- Apply friction
    p.vx = p.vx * 0.9
    
    -- Check platform collisions
    p.on_ground = false
    for plat in all(platforms) do
      if check_collision(p, plat) then
        p.y = plat.y - 16
        p.vy = 0
        p.on_ground = true
      end
    end
    
    -- Wrap screen
    if p.x < 0 then p.x = 480 end
    if p.x > 480 then p.x = 0 end
    
    -- Check if player fell off bottom
    if p.y > 270 then
      p.y = 100
      p.x = 50 + (id-1) * 100
      score[id] = score[id] - 10  -- Penalty
    end
  end
  
  -- Engine automatically syncs players and score tables!
end

function _draw()
  -- Both host and non-host render the same way
  cls(0)  -- Clear to black
  
  -- Draw platforms
  for plat in all(platforms) do
    rectfill(plat.x, plat.y, plat.x + plat.w, plat.y + plat.h, 7)
  end
  
  -- Draw all players
  for id, p in pairs(players) do
    spr(p.sprite, p.x, p.y)
    
    -- Show player ID above sprite
    local color = (id == rf.my_player_id()) and 7 or 6
    print(id, p.x + 4, p.y - 8, color)
  end
  
  -- Draw scores
  for id, s in pairs(score) do
    print("P" .. id .. ": " .. s, 10, 10 + (id-1) * 10, 7)
  end
  
  -- Show connection info
  if rf.is_multiplayer() then
    print("Players: " .. rf.player_count(), 10, 250, 7)
    if rf.is_host() then
      print("HOST", 10, 260, 10)
    end
  end
end

function check_collision(player, platform)
  -- Simple AABB collision
  return player.x < platform.x + platform.w and
         player.x + 16 > platform.x and
         player.y < platform.y + platform.h and
         player.y + 16 > platform.y and
         player.vy >= 0  -- Only collide when falling
end
```

**manifest.json:**
```json
{
  "version": "1.0",
  "title": "Multiplayer Platformer Demo",
  "author": "RetroForge Team",
  "description": "Simple platformer demonstrating multiplayer sync",
  
  "multiplayer": {
    "enabled": true,
    "minPlayers": 2,
    "maxPlayers": 6,
    "supportsSolo": true,
    "description": "Race to the top, avoid falling!"
  },
  
  "tags": ["demo", "platformer", "multiplayer"],
  "cartVersion": "1.0.0",
  "engineVersion": "2.0.0"
}
```

**Key Points:**
- Game developer writes almost normal PICO-8 style code
- Only differences:
  1. Check `rf.is_host()` to know who runs logic
  2. Use `rf.btn(player_id, button)` to read other players' inputs
  3. Call `rf.network_sync()` to register tables
- Engine handles all networking automatically
- No manual packet sending/receiving
- No WebRTC code visible to developer

---

## Security & Performance

### Security Measures

#### Host Authority Prevents Cheating

**Problem:** Player modifies their own score
```lua
-- Malicious player's code (won't work!)
score[rf.my_player_id()] = 99999
```

**Solution:** Engine enforces ownership
- Non-host player can only modify `players[their_id]`
- Attempts to modify other keys are ignored
- Host is source of truth for `score` table
- Host validates all state changes

#### Bandwidth Limiting

**Problem:** Malicious cart syncs huge tables
```lua
-- Bad cart code
huge_table = {}
for i = 1, 100000 do
  huge_table[i] = {x=i, y=i, data="spam"}
end
rf.network_sync(huge_table, "fast")  -- Would use 60MB/sec!
```

**Solution:** Engine warns and rate-limits
```go
// In RetroForge engine
const MaxSyncTableSize = 10 * 1024 // 10KB per table

func (nm *NetworkManager) RegisterSync(table LuaTable, freq string) error {
    size := calculateTableSize(table)
    if size > MaxSyncTableSize {
        log.Printf("WARNING: Table size %d exceeds %d bytes", size, MaxSyncTableSize)
        log.Printf("This may cause network lag. Consider using slower sync tier.")
        
        // Still allow, but warn developer
    }
    // ... register table
}
```

#### Disconnect Protection

**Problem:** Host rage-quits when losing

**Solution:** Game ends for all players
- No advantage to host disconnecting
- Match results still saved (based on last known state)
- Players returned to lobby gracefully

---

### Performance Optimization

#### Client-Side Techniques

**1. Interpolation (Smooth Movement)**
```go
// When receiving position updates at 60/sec
// Interpolate between old and new positions over ~16ms

type PlayerState struct {
    Current  Vector2
    Previous Vector2
    Target   Vector2
    LerpTime float64
}

func (ps *PlayerState) Update(dt float64) {
    ps.LerpTime += dt
    t := ps.LerpTime / 0.016 // 16ms interpolation window
    
    ps.Current.X = lerp(ps.Previous.X, ps.Target.X, t)
    ps.Current.Y = lerp(ps.Previous.Y, ps.Target.Y, t)
}
```

**2. Dead Reckoning (Predict Movement)**
```go
// If no update received, predict based on velocity
func (ps *PlayerState) Predict(dt float64) {
    if timeSinceLastUpdate > 0.1 { // 100ms no update
        ps.Current.X += ps.Velocity.X * dt
        ps.Current.Y += ps.Velocity.Y * dt
    }
}
```

**3. Object Pooling (Reduce GC)**
```lua
-- Reuse bullet objects instead of creating new ones
bullet_pool = {}

function spawn_bullet(x, y)
  local bullet = table.remove(bullet_pool) or {}
  bullet.x = x
  bullet.y = y
  bullet.active = true
  return bullet
end

function despawn_bullet(bullet)
  bullet.active = false
  table.insert(bullet_pool, bullet)
end
```

#### Host-Side Techniques

**1. Delta Compression**
```go
// Only send changed values
func (nm *NetworkManager) GenerateDelta(prev, current State) Delta {
    delta := Delta{}
    
    for key, value := range current {
        if prev[key] != value {
            delta[key] = value  // Only include if changed
        }
    }
    
    return delta
}
```

**2. Batch Updates**
```go
// Collect all changes in a frame, send once
type BatchedUpdate struct {
    Frame   uint64
    Changes map[string]interface{}
}

func (nm *NetworkManager) SendBatchedUpdate() {
    batch := nm.CollectChanges()
    
    // Send to all connections
    for _, conn := range nm.connections {
        conn.Send(batch)
    }
}
```

**3. Priority Queue**
```go
// Send critical updates immediately, defer less important
type UpdatePriority int

const (
    PriorityHigh   UpdatePriority = 0  // Player positions
    PriorityMedium UpdatePriority = 1  // Powerups
    PriorityLow    UpdatePriority = 2  // Scores
)

func (nm *NetworkManager) QueueUpdate(priority UpdatePriority, data interface{}) {
    queue := nm.priorityQueues[priority]
    queue.Push(data)
}
```

---

### Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Frame Rate | 60 FPS | Consistent, no drops |
| Input Latency (host) | <16ms | Local button to action |
| Input Latency (non-host) | <200ms | Button to visible action |
| State Sync Latency | <50ms | Host to non-host updates |
| Initial Connection Time | <5 sec | WebRTC establishment |
| Bandwidth (host) | <500 KB/sec | Upload with 5 players |
| Bandwidth (non-host) | <100 KB/sec | Download + input upload |
| Memory Usage | <100 MB | WASM heap |
| Cart Load Time | <2 sec | .rfs file to playable |

---

## Development & Testing

### Local Testing Workflow

#### Option 1: Multiple Browser Tabs

```bash
# Terminal 1: Start Next.js dev server
npm run dev

# Browser:
# Tab 1 (Host): http://localhost:3000/game/mycart?player=1
# Tab 2: http://localhost:3000/game/mycart?player=2
# Tab 3: http://localhost:3000/game/mycart?player=3
```

**Engine detects `?player=N` in dev mode:**
- Simulates multiple players on localhost
- Each tab gets different player_id
- WebRTC still used (loopback connections)
- Can test lobby/matchmaking flow

#### Option 2: Fake Multiplayer Mode

```lua
-- In cart code (development only)
function _init()
  if rf.stat then  -- Only available in dev mode
    -- Simulate 4 players locally
    rf.debug_fake_multiplayer(4)
  end
  
  -- Rest of init code works normally
  for i = 1, rf.player_count() do
    players[i] = {x=100*i, y=100}
  end
end
```

**Engine behavior in fake mode:**
- `rf.player_count()` returns 4
- `rf.my_player_id()` returns 1 (always host)
- `rf.btn(2, 0)` can be controlled by developer
- No actual network traffic
- Perfect for testing game logic without network

#### Option 3: Network Simulation

```lua
-- Simulate poor network conditions
rf.debug_set_latency(150)      -- 150ms delay
rf.debug_set_packet_loss(0.05) -- 5% packet loss
rf.debug_set_jitter(20)         -- ±20ms jitter

-- Test how game behaves with lag
```

---

### Developer Tools

#### Network Debug Panel (Dev Mode Only)

```lua
function _draw()
  -- Normal game rendering
  render_game()
  
  -- Dev mode: show network stats
  if rf.stat then
    local y = 10
    print("=== Network Debug ===", 10, y, 7)
    y += 10
    
    print("Players: " .. rf.player_count(), 10, y, 7)
    y += 8
    
    for i = 1, rf.player_count() do
      local ping = rf.network_ping(i)
      local color = ping < 50 and 11 or ping < 100 and 9 or 8
      print("P" .. i .. " ping: " .. ping .. "ms", 10, y, color)
      y += 8
    end
    
    local bw = rf.network_bandwidth()
    print("Bandwidth: " .. bw .. " KB/s", 10, y, 7)
    y += 8
    
    local loss = rf.network_packet_loss()
    print("Packet loss: " .. (loss * 100) .. "%", 10, y, 7)
  end
end
```

#### Console Logging (Dev Mode)

```lua
-- Only works when running from folder (not packed .rfe)
function _update()
  if rf.stat then
    rf.printh("Frame: " .. rf.stat(1))
    rf.printh("FPS: " .. rf.stat(0))
    rf.printh("Memory: " .. rf.stat(2) .. " bytes")
  end
end
```

---

### Testing Checklist

**Connection Tests:**
- [ ] 2 players connect successfully
- [ ] 6 players connect successfully
- [ ] Player joins after lobby created
- [ ] Player leaves before game starts
- [ ] Player disconnects during game
- [ ] Host disconnects during game
- [ ] Connection retry on failure (3 attempts)
- [ ] Kick player after 3 failed attempts

**Sync Tests:**
- [ ] Fast tier updates 30-60 times/sec
- [ ] Moderate tier updates ~15 times/sec
- [ ] Slow tier updates ~5 times/sec
- [ ] Delta encoding reduces bandwidth
- [ ] Large tables trigger warning
- [ ] Ownership rules enforced (non-host can't modify others)

**Gameplay Tests:**
- [ ] Input latency acceptable (~200ms)
- [ ] Visual movement smooth (interpolation works)
- [ ] Physics deterministic (all players see same result)
- [ ] Score updates correctly
- [ ] Game over triggers for all players
- [ ] Results saved to Convex

**Performance Tests:**
- [ ] 60 FPS maintained with 6 players
- [ ] Bandwidth under 500 KB/sec for host
- [ ] Memory usage under 100 MB
- [ ] No memory leaks over 10 minute game
- [ ] Works on low-end hardware

**Edge Cases:**
- [ ] Handle rapid button mashing
- [ ] Handle player AFKing
- [ ] Handle network spike (500ms+ latency)
- [ ] Handle packet burst loss
- [ ] Handle browser tab backgrounding
- [ ] Handle mobile screen lock

---

## Implementation Notes

### Phase 1: Core Networking (Foundation)

**Goal:** Get two players connected and seeing each other move

**Tasks:**
1. Implement WebRTC connection in Go (pion/webrtc)
2. Compile to WASM with TinyGo
3. Create signaling server in Convex
4. Establish star topology connections
5. Implement basic state sync (no deltas yet)
6. Test with simple demo cart

**Deliverable:** Two players can join and see each other's positions update

---

### Phase 2: API & Sync System

**Goal:** Complete multiplayer API with 3-tier sync

**Tasks:**
1. Implement `rf.network_sync(table, tier)` in Lua binding
2. Add delta encoding for bandwidth efficiency
3. Implement ownership enforcement
4. Add input collection at 5/sec
5. Create example multiplayer carts
6. Write API documentation

**Deliverable:** Full API working, developers can create multiplayer games

---

### Phase 3: Polish & Optimization

**Goal:** Production-ready performance and UX

**Tasks:**
1. Add interpolation for smooth movement
2. Implement dead reckoning
3. Add bandwidth warnings
4. Create debug tools (network panel, fake multiplayer)
5. Optimize packet sizes
6. Add connection retry logic
7. Handle disconnections gracefully

**Deliverable:** Smooth, polished multiplayer experience

---

### Phase 4: Next.js Integration

**Goal:** Complete platform with matchmaking

**Tasks:**
1. Build game library UI
2. Implement lobby system in Convex
3. Create lobby browser/creator
4. Add match history and statistics
5. Build leaderboards
6. Integrate cart uploads
7. Add user profiles

**Deliverable:** Full platform like itch.io but for RetroForge games

---

## Appendices

### A. Network Packet Size Reference

```
Input Packet:        ~50 bytes
State Delta (fast):  100-1000 bytes (depends on changes)
State Delta (mod):   50-500 bytes
State Delta (slow):  20-100 bytes
Full Snapshot:       1-10 KB (one-time on connect)
Game Over:           ~200 bytes
Player Join/Leave:   ~100 bytes
```

### B. WebRTC Configuration

```json
{
  "iceServers": [
    {
      "urls": "stun:stun.l.google.com:19302"
    },
    {
      "urls": "turn:turn.retroforge.dev:3478",
      "username": "retroforge",
      "credential": "shared_secret"
    }
  ],
  "iceTransportPolicy": "all",
  "bundlePolicy": "balanced",
  "rtcpMuxPolicy": "require"
}
```

**Notes:**
- STUN server for NAT traversal (free, Google's)
- TURN server for relay (fallback, requires hosting)
- Most connections work with STUN only
- TURN only needed for strict corporate firewalls (~5% of users)

### C. Bandwidth Budget Breakdown

**6-Player Game Example:**

**Host Upload:**
```
Fast tier (60/sec):   5 players × 1KB × 60 = 300 KB/sec
Moderate (15/sec):    5 players × 0.5KB × 15 = 37.5 KB/sec
Slow (5/sec):         5 players × 0.1KB × 5 = 2.5 KB/sec
Total: 340 KB/sec = 2.7 Mbps upload
```

**Host Download (inputs):**
```
5 players × 50 bytes × 5/sec = 1.25 KB/sec = negligible
```

**Non-Host Upload (inputs):**
```
50 bytes × 5/sec = 250 bytes/sec = negligible
```

**Non-Host Download:**
```
Fast: 1KB × 60/sec = 60 KB/sec
Moderate: 0.5KB × 15/sec = 7.5 KB/sec
Slow: 0.1KB × 5/sec = 0.5 KB/sec
Total: 68 KB/sec = 544 Kbps download
```

**Conclusion:**
- Host needs decent upload (3+ Mbps)
- Non-hosts need minimal upload, decent download (1+ Mbps)
- Mobile hotspot may struggle as host
- Wi-Fi/Ethernet recommended for host

### D. Glossary

- **Cart:** RetroForge game file (.rfs or .rfe)
- **Delta Encoding:** Sending only changed values, not full state
- **Host Authority:** One player controls game logic, others are clients
- **Interpolation:** Smoothing between network updates
- **Latency:** Time delay for data to travel over network
- **Lobby:** Pre-game waiting room for players to gather
- **Packet:** Single network message
- **Peer:** Another player in the game
- **Signaling:** Exchange of connection info (SDP/ICE) before WebRTC starts
- **Star Topology:** Hub-and-spoke network (host in center)
- **State Sync:** Synchronizing game state across players
- **STUN/TURN:** Protocols for NAT traversal (connecting through routers)
- **WebRTC:** Browser technology for real-time P2P communication

---

## Summary

This design document specifies a complete multiplayer system for RetroForge:

**✅ Platform (Next.js/Convex):**
- Matchmaking and lobby management
- WebRTC signaling server
- User profiles and statistics
- Game library

**✅ Engine (Go/WASM):**
- Built-in WebRTC networking
- Automatic state synchronization
- 3-tier sync system (fast/moderate/slow)
- Minimal developer-facing API
- Host authority model
- Star network topology

**✅ Developer Experience:**
- Write normal game code
- Add 2-3 API calls for multiplayer
- Engine handles all networking
- Test locally with multiple tabs
- Debug tools included

**✅ Performance:**
- 60 FPS maintained
- <200ms input latency
- <500 KB/sec host bandwidth
- Supports 2-6 players

**Ready for implementation!** 🚀

---

**Document Version:** 2.0 (Multiplayer Edition)  
**Last Updated:** October 31, 2025  
**Status:** FINAL - Ready for Development