package cartio

// MountPoint represents a point within a sprite where projectiles/thrusters originate
type MountPoint struct {
	X    int    `json:"x"`              // X coordinate within sprite bounds
	Y    int    `json:"y"`              // Y coordinate within sprite bounds
	Name string `json:"name,omitempty"` // Optional name for accessing by name in Lua
}

// SpriteData represents a single sprite
type SpriteData struct {
	Width        int          `json:"width"`        // Sprite width in pixels
	Height       int          `json:"height"`       // Sprite height in pixels
	Pixels       [][]int      `json:"pixels"`       // 2D array of color indices (0-49, -1 for transparent)
	UseCollision bool         `json:"useCollision"` // Enable collision detection with other sprites
	MountPoints  []MountPoint `json:"mountPoints"`  // Array of mount points (e.g., for bullets, thrusters)
	IsUI         bool         `json:"isUI"`         // If true, sprite is UI element and not affected by physics
	Lifetime     int          `json:"lifetime"`     // Lifetime in milliseconds (0 = no lifetime limit)
	MaxSpawn     int          `json:"maxSpawn"`     // Maximum instances that can be spawned simultaneously (0 = no limit)
}

// SpriteMap maps sprite names to their data
type SpriteMap map[string]SpriteData
