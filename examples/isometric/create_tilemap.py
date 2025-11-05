#!/usr/bin/env python3
"""
Generate a 10x10 tilemap for testing.
"""

import json
import random

def main():
    # Create a 10x10 tilemap
    width = 10
    height = 10
    
    # Tile names map: 0 = empty, 1 = earth, 2 = grass, 3 = desert, 4 = snow, 5 = water, 6 = lava, 7 = rock
    tile_names = ['', 'earth', 'grass', 'desert', 'snow', 'water', 'lava', 'rock']
    
    tiles = []
    
    # Create a simple pattern for testing
    for y in range(height):
        row = []
        for x in range(width):
            # Create a pattern: grass in center, water around edges, earth in corners
            if (x == 0 or x == width-1) and (y == 0 or y == height-1):
                tile_idx = 1  # earth
            elif x == 0 or x == width-1 or y == 0 or y == height-1:
                tile_idx = 5  # water
            elif x == 4 or x == 5 or y == 4 or y == 5:
                tile_idx = 2  # grass
            elif (x + y) % 3 == 0:
                tile_idx = 3  # desert
            elif (x + y) % 5 == 0:
                tile_idx = 7  # rock
            else:
                tile_idx = 2  # grass (default)
            
            # Store as tile name string (0 = empty string)
            row.append(tile_names[tile_idx] if tile_idx > 0 else "")
        tiles.append(row)
    
    tilemap = {
        "tilesetName": "terrain",
        "width": width,
        "height": height,
        "tiles": tiles
    }
    
    # Write tilemap file
    output_path = "assets/test_map.json"
    with open(output_path, 'w') as f:
        json.dump(tilemap, f, indent=2)
    
    print(f"✓ Created tilemap: {output_path}")
    print(f"  Size: {width}x{height}")
    print(f"  Tileset: terrain")

if __name__ == "__main__":
    main()

