#!/usr/bin/env python3
"""
Generate tileset for isometric example.
Creates tiles: earth, grass, desert, snow, water, lava, rock
All tiles are 16x16 static sprites.
"""

import json
import os

# Palette colors (using Pastel 50 palette indices)
# We'll use a subset of colors for terrain tiles
PALETTE_INDICES = {
    'earth': 12,      # Brown
    'grass': 17,      # Green
    'desert': 15,     # Yellow
    'snow': 0,        # White
    'water': 22,      # Blue
    'lava': 23,       # Red/Orange
    'rock': 7,        # Gray
}

def create_tile_pixels(tile_name, color_idx):
    """Create a simple 16x16 tile filled with the color."""
    pixels = []
    for y in range(16):
        row = []
        for x in range(16):
            row.append(color_idx)
        pixels.append(row)
    return pixels

def create_pattern_tile(tile_name, base_color_idx, pattern_color_idx=None):
    """Create a tile with a simple pattern."""
    pixels = []
    for y in range(16):
        row = []
        for x in range(16):
            # Simple checkerboard pattern for grass
            if tile_name == 'grass' and (x + y) % 4 < 2:
                row.append(base_color_idx + 1 if base_color_idx < 49 else base_color_idx)
            # Sand pattern for desert
            elif tile_name == 'desert' and (x * 3 + y * 2) % 7 < 3:
                row.append(base_color_idx - 1 if base_color_idx > 0 else base_color_idx)
            # Water pattern (waves)
            elif tile_name == 'water' and (x + y * 2) % 6 < 2:
                row.append(base_color_idx + 2 if base_color_idx < 48 else base_color_idx)
            # Rock pattern (random speckles)
            elif tile_name == 'rock' and (x * 5 + y * 7) % 11 < 3:
                row.append(base_color_idx - 2 if base_color_idx > 1 else base_color_idx)
            else:
                row.append(base_color_idx)
        pixels.append(row)
    return pixels

def main():
    tileset = {}
    
    # Create each tile type
    tiles = ['earth', 'grass', 'desert', 'snow', 'water', 'lava', 'rock']
    
    for tile_name in tiles:
        color_idx = PALETTE_INDICES[tile_name]
        
        # Create pixels with pattern for visual interest
        if tile_name in ['grass', 'desert', 'water', 'rock']:
            pixels = create_pattern_tile(tile_name, color_idx)
        else:
            pixels = create_tile_pixels(tile_name, color_idx)
        
        tileset[tile_name] = {
            "width": 16,
            "height": 16,
            "type": "static",
            "pixels": pixels,
            "useCollision": False,
            "mountPoints": [],
            "isISO": True  # Mark as isometric tile
        }
    
    # Write tileset file
    output_path = "assets/terrain_tiles.json"
    with open(output_path, 'w') as f:
        json.dump(tileset, f, indent=2)
    
    print(f"✓ Created tileset with {len(tiles)} tiles: {output_path}")
    print(f"  Tiles: {', '.join(tiles)}")

if __name__ == "__main__":
    main()

