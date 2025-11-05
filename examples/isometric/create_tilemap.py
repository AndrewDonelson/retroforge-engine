#!/usr/bin/env python3
"""
Generate a 16x16 tilemap for testing isometric tiles.
"""

import json

def main():
    # Create a 16x16 tilemap
    width = 16
    height = 16
    
    # Tile names from the PNG-generated tileset
    # Available tiles: dark_grass, deep_water, dirt, grass, gravel, ice, mud, sand, shallow_water, snow, stone, water
    tile_names = ['', 'grass', 'water', 'sand', 'snow', 'deep_water', 'stone', 'dirt']
    
    tiles = []
    
    # Create a more interesting pattern for 16x16 map
    for y in range(height):
        row = []
        for x in range(width):
            # Create varied terrain patterns:
            # - Water around edges (ocean)
            # - Earth/sand near edges (beach)
            # - Grass in center areas (plains)
            # - Desert patches
            # - Rock/mountain areas
            # - Snow patches
            # - Lava in corners
            
            # Corners: stone
            if (x == 0 or x == width-1) and (y == 0 or y == height-1):
                tile_idx = 6  # stone
            # Outer edges: deep_water
            elif x == 0 or x == width-1 or y == 0 or y == height-1:
                tile_idx = 5  # deep_water
            # Second ring: sand/dirt (beach)
            elif x == 1 or x == width-2 or y == 1 or y == height-2:
                if (x + y) % 2 == 0:
                    tile_idx = 7  # dirt
                else:
                    tile_idx = 3  # sand
            # Inner area: varied terrain
            elif x >= 4 and x <= 11 and y >= 4 and y <= 11:
                # Center area: mostly grass with some variation
                if (x + y) % 7 == 0:
                    tile_idx = 6  # stone
                elif (x + y) % 5 == 0:
                    tile_idx = 4  # snow
                elif (x * 2 + y) % 11 == 0:
                    tile_idx = 3  # sand
                else:
                    tile_idx = 1  # grass
            # Middle ring: mix of grass and sand
            elif (x + y) % 3 == 0:
                tile_idx = 3  # sand
            elif (x + y) % 5 == 0:
                tile_idx = 6  # stone
            else:
                tile_idx = 1  # grass (default)
            
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
    print(f"  Pattern: Ocean edges → Beach → Varied terrain center")

if __name__ == "__main__":
    main()
