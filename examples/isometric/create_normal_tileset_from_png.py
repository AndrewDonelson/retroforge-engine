#!/usr/bin/env python3
"""
Generate normal (orthogonal) tileset from PNG files in design/tiles.

This script:
1. Loads PNG files from design/tiles directory
2. Converts RGB colors to RetroForge 48 palette indices (game palette)
3. Creates a normal (non-isometric) tileset JSON
4. Includes all base terrain tiles (excludes transition tiles)
"""

import json
import os
import math
from PIL import Image

# Path to the tiles directory
def get_tiles_dir():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    # Go up: examples/isometric -> examples -> retroforge-engine -> RetroForge -> design/tiles
    engine_root = os.path.abspath(os.path.join(script_dir, '../..'))  # retroforge-engine
    retroforge_root = os.path.dirname(engine_root)   # RetroForge parent directory
    tiles_dir = os.path.join(retroforge_root, 'design', 'tiles')
    return tiles_dir

TILES_DIR = get_tiles_dir()

# Path to palette.json (same as terraingen uses)
def get_palette_path():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    engine_root = os.path.abspath(os.path.join(script_dir, '../..'))  # retroforge-engine
    palette_path = os.path.join(engine_root, 'palette.json')
    return palette_path

PALETTE_PATH = get_palette_path()

# Load RetroForge 48 palette from palette.json (matches engine's pal.GetPalette("RetroForge 48"))
def load_palette_from_json():
    """Load the RetroForge 48 game palette from palette.json to match the engine.
    This is the 48-color game palette (indices 16-63 in full 64-color system)."""
    try:
        with open(PALETTE_PATH, 'r') as f:
            palette_data = json.load(f)
        
        def hex_to_rgb(hex_str):
            hex_str = hex_str.lstrip('#')
            return tuple(int(hex_str[i:i+2], 16) for i in (0, 2, 4))
        
        # Convert hex colors to RGB tuples (game palette is 48 colors)
        palette = []
        for hex_color in palette_data['colors'][:48]:  # Take first 48 colors (game palette)
            palette.append(hex_to_rgb(hex_color))
        
        if len(palette) < 48:
            raise ValueError(f"palette.json has only {len(palette)} colors, expected 48")
        
        return palette
    except Exception as e:
        print(f"Warning: Could not load palette.json ({e}), using fallback palette")
        # Fallback to palette.json values as hardcoded (should match if palette.json is correct)
        return [
            (0, 0, 0), (255, 255, 255), (255, 137, 137), (255, 77, 77), (195, 17, 17),
            (255, 205, 137), (255, 145, 77), (195, 85, 17), (255, 255, 137), (255, 216, 77),
            (195, 156, 17), (242, 255, 137), (182, 255, 77), (122, 195, 17), (137, 255, 195),
            (77, 212, 135), (17, 152, 75), (114, 255, 255), (54, 216, 199), (0, 156, 139),
            (137, 255, 255), (77, 213, 255), (17, 153, 195), (162, 251, 255), (102, 191, 255),
            (42, 131, 195), (171, 196, 255), (111, 136, 255), (51, 76, 195), (198, 177, 255),
            (138, 117, 255), (78, 57, 195), (240, 180, 255), (180, 120, 255), (108, 24, 195),
            (255, 191, 217), (255, 111, 177), (195, 51, 81), (255, 191, 216), (255, 127, 160),
            (195, 67, 64), (228, 181, 150), (168, 121, 90), (108, 61, 30), (220, 237, 150),
            (160, 177, 90), (100, 117, 30), (116, 249, 255), (56, 189, 248), (0, 129, 188),
        ]

# RetroForge 48 palette - loaded from palette.json to match engine (game palette, indices 16-63)
RETROFORGE_PALETTE = load_palette_from_json()

def color_distance(r1, g1, b1, r2, g2, b2):
    """Calculate Euclidean distance between two RGB colors."""
    dr = r1 - r2
    dg = g1 - g2
    db = b1 - b2
    return math.sqrt(dr*dr + dg*dg + db*db)

def find_closest_palette_index(r, g, b):
    """Find the closest palette color index for an RGB color.
    First checks for exact matches, then falls back to closest match."""
    # First check for exact match (handles PNGs saved with exact palette colors)
    for idx, (pr, pg, pb) in enumerate(RETROFORGE_PALETTE):
        if (r, g, b) == (pr, pg, pb):
            return idx
    
    # Fall back to closest match if no exact match found
    min_dist = float('inf')
    best_idx = 0
    
    for idx, (pr, pg, pb) in enumerate(RETROFORGE_PALETTE):
        dist = color_distance(r, g, b, pr, pg, pb)
        if dist < min_dist:
            min_dist = dist
            best_idx = idx
    
    return best_idx

def load_png_to_pixels(png_path):
    """Load a PNG image and convert it to palette-indexed pixels."""
    try:
        img = Image.open(png_path)
        # Convert to RGB if needed
        if img.mode != 'RGB' and img.mode != 'RGBA':
            img = img.convert('RGB')
        
        width, height = img.size
        pixels = []
        
        for y in range(height):
            row = []
            for x in range(width):
                if img.mode == 'RGBA':
                    r, g, b, a = img.getpixel((x, y))
                    if a < 128:
                        row.append(-1)  # Transparent
                        continue
                else:
                    r, g, b = img.getpixel((x, y))
                
                # Find closest palette color
                palette_idx = find_closest_palette_index(r, g, b)
                row.append(palette_idx)
            pixels.append(row)
        
        return pixels
    except Exception as e:
        print(f"Error loading {png_path}: {e}")
        return None

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Check if tiles directory exists
    if not os.path.exists(TILES_DIR):
        print(f"ERROR: Tiles directory not found: {TILES_DIR}")
        return 1
    
    # Find all PNG files
    png_files = [f for f in os.listdir(TILES_DIR) if f.endswith('_16x16.png')]
    if not png_files:
        print(f"ERROR: No PNG files found in {TILES_DIR}")
        return 1
    
    print(f"Found {len(png_files)} PNG files")
    print()
    
    # Map of terrain names to their base names (excluding transitions)
    terrain_tiles = {}
    
    for png_file in sorted(png_files):
        # Extract tile name (remove _16x16.png suffix)
        tile_name = png_file.replace('_16x16.png', '')
        png_path = os.path.join(TILES_DIR, png_file)
        
        # Skip transition tiles for normal tileset (we only want base terrains)
        if '_to_' in tile_name:
            continue
        
        terrain_tiles[tile_name] = png_path
    
    print(f"Loading {len(terrain_tiles)} base terrain tiles...")
    print()
    
    # Create tileset
    tileset = {}
    for tile_name in sorted(terrain_tiles.keys()):
        png_path = terrain_tiles[tile_name]
        print(f"  Loading {tile_name}...")
        
        pixels = load_png_to_pixels(png_path)
        if pixels is None:
            print(f"    ✗ Failed to load {tile_name}")
            continue
        
        # Create tile data
        tile_data = {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': pixels,
            'useCollision': False,
            'mountPoints': []
        }
        
        tileset[tile_name] = tile_data
        print(f"    ✓ Loaded {tile_name} (16x16)")
    
    # Create final tileset with isISO=false (normal tiles)
    output_tileset = {
        'isISO': False,
        'tiles': tileset
    }
    
    # Save tileset
    output_path = os.path.join(script_dir, 'assets', 'terrain_normal_tiles.json')
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    
    with open(output_path, 'w') as f:
        json.dump(output_tileset, f, indent=2)
    
    print()
    print(f"✓ Created normal tileset: {output_path}")
    print(f"  Tiles: {', '.join(sorted(tileset.keys()))}")
    print(f"  Palette: RetroForge 48 (game palette, indices 16-63 in full system)")
    print(f"  Isometric: False (normal orthogonal tiles)")
    print(f"  Source: PNG files from {TILES_DIR}")
    print()
    
    return 0

if __name__ == '__main__':
    exit(main())

