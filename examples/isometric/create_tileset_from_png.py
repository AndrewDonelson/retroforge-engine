#!/usr/bin/env python3
"""
Generate isometric tileset from PNG files using tile2iso converter.

This script:
1. Loads PNG files from design/tiles directory
2. Converts RGB colors to RetroForge 50 palette indices
3. Creates three sprites (top, left, right) for each terrain type
4. Uses the tile2iso CLI tool to convert them into 2.5D isometric tiles
5. Outputs a tileset with isISO=true flag
"""

import json
import os
import subprocess
import tempfile
import shutil
from PIL import Image
import math

# Path to the tiles directory
# Script is in: retroforge-engine/examples/isometric/
# Tiles are in: design/tiles/ (sibling directory to retroforge-engine)
def get_tiles_dir():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    # Go up: examples/isometric -> examples -> retroforge-engine -> RetroForge -> design/tiles
    engine_root = os.path.abspath(os.path.join(script_dir, '../..'))  # retroforge-engine
    retroforge_root = os.path.dirname(engine_root)   # RetroForge parent directory
    tiles_dir = os.path.join(retroforge_root, 'design', 'tiles')
    return tiles_dir

TILES_DIR = get_tiles_dir()

# RetroForge 50 palette (matching the Go terraingen tool)
RETROFORGE_PALETTE = [
    (0, 0, 0),           # 0: Black
    (29, 43, 83),        # 1: Dark blue
    (126, 37, 83),       # 2: Dark purple
    (0, 135, 81),        # 3: Dark green
    (171, 82, 54),       # 4: Brown
    (95, 87, 79),        # 5: Dark gray
    (194, 195, 199),     # 6: Light gray
    (255, 241, 232),     # 7: White
    (255, 0, 77),        # 8: Red
    (255, 163, 0),       # 9: Orange
    (255, 236, 39),      # 10: Yellow
    (0, 228, 54),        # 11: Green
    (41, 173, 255),      # 12: Blue
    (131, 118, 156),     # 13: Lavender
    (255, 119, 168),     # 14: Pink
    (255, 204, 170),     # 15: Peach
    (34, 32, 52),        # 16: Very dark blue
    (69, 40, 60),        # 17: Dark purple-brown
    (102, 57, 49),       # 18: Dark brown
    (143, 86, 59),       # 19: Medium brown
    (223, 113, 38),      # 20: Light brown/orange
    (217, 160, 102),     # 21: Tan
    (238, 195, 154),     # 22: Light tan
    (251, 242, 54),      # 23: Bright yellow
    (153, 229, 80),      # 24: Light green
    (106, 190, 48),      # 25: Medium green
    (55, 148, 110),      # 26: Teal green
    (75, 105, 47),       # 27: Dark olive
    (82, 75, 36),        # 28: Olive brown
    (50, 60, 57),        # 29: Dark teal
    (63, 63, 116),       # 30: Dark blue
    (48, 96, 130),       # 31: Ocean blue
    (91, 110, 225),      # 32: Bright blue
    (99, 155, 255),      # 33: Sky blue
    (95, 205, 228),      # 34: Light blue
    (203, 219, 252),     # 35: Very light blue
    (155, 173, 183),     # 36: Blue-gray
    (132, 126, 135),     # 37: Medium gray
    (105, 106, 106),     # 38: Dark gray
    (89, 86, 82),        # 39: Very dark gray
    (118, 66, 138),      # 40: Purple
    (172, 50, 50),       # 41: Dark red
    (217, 87, 99),       # 42: Pink-red
    (215, 123, 186),     # 43: Pink
    (143, 151, 74),      # 44: Yellow-green
    (138, 111, 48),      # 45: Gold-brown
    (194, 195, 199),     # 46: Light gray
    (255, 255, 255),     # 47: Pure white
    (0, 0, 0),           # 48: Black
    (0, 0, 0),           # 49: Black
]

def color_distance(r1, g1, b1, r2, g2, b2):
    """Calculate Euclidean distance between two RGB colors."""
    dr = r1 - r2
    dg = g1 - g2
    db = b1 - b2
    return math.sqrt(dr*dr + dg*dg + db*db)

def find_closest_palette_index(r, g, b):
    """Find the closest palette color index for an RGB color."""
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
        if img.mode != 'RGB':
            img = img.convert('RGB')
        
        width, height = img.size
        pixels = []
        
        for y in range(height):
            row = []
            for x in range(width):
                r, g, b = img.getpixel((x, y))
                # Check for transparency (alpha channel if present)
                if img.mode == 'RGBA':
                    r, g, b, a = img.getpixel((x, y))
                    if a < 128:
                        row.append(-1)  # Transparent
                        continue
                
                # Find closest palette color
                palette_idx = find_closest_palette_index(r, g, b)
                row.append(palette_idx)
            pixels.append(row)
        
        return pixels
    except Exception as e:
        print(f"Error loading {png_path}: {e}")
        return None

def create_sprites_from_png(png_path, tile_name):
    """Create three sprites (top, left, right) from a single PNG tile."""
    pixels = load_png_to_pixels(png_path)
    if pixels is None:
        return None
    
    # All three sprites use the same top-down view
    # The tile2iso converter will handle the side face generation
    sprites = {
        f"{tile_name}_top": {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': pixels,
            'useCollision': False,
            'isUI': False,
            'mountPoints': [],
            'lifetime': 0,
            'maxSpawn': 0,
        },
        f"{tile_name}_left": {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': pixels,  # Same as top for now
            'useCollision': False,
            'isUI': False,
            'mountPoints': [],
            'lifetime': 0,
            'maxSpawn': 0,
        },
        f"{tile_name}_right": {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': pixels,  # Same as top for now
            'useCollision': False,
            'isUI': False,
            'mountPoints': [],
            'lifetime': 0,
            'maxSpawn': 0,
        },
    }
    
    return sprites

def build_tile2iso():
    """Build the tile2iso CLI tool if needed."""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    engine_root = os.path.join(script_dir, '../..')
    tile2iso_path = os.path.join(engine_root, 'bin', 'tile2iso')
    
    if not os.path.exists(tile2iso_path):
        print(f"Building tile2iso tool...")
        result = subprocess.run(
            ['go', 'build', '-o', tile2iso_path, './cmd/tile2iso'],
            cwd=engine_root,
            capture_output=True,
            text=True
        )
        if result.returncode != 0:
            print(f"ERROR: Failed to build tile2iso: {result.stderr}")
            return None
    
    return tile2iso_path

def create_palette_json(palette_path):
    """Create a palette.json file with RetroForge 50 colors."""
    colors = []
    for r, g, b in RETROFORGE_PALETTE:
        hex_color = f"#{r:02x}{g:02x}{b:02x}"
        colors.append(hex_color)
    
    palette_data = {'colors': colors}
    
    with open(palette_path, 'w') as f:
        json.dump(palette_data, f, indent=2)
    
    return palette_path

def convert_tile_to_isometric(tile2iso_path, sprites_path, palette_path, tile_name, temp_dir):
    """Convert a tile to isometric using tile2iso."""
    top_name = f"{tile_name}_top"
    left_name = f"{tile_name}_left"
    right_name = f"{tile_name}_right"
    
    output_path = os.path.join(temp_dir, f"{tile_name}_iso.json")
    
    cmd = [
        tile2iso_path,
        'convert',
        '--sprites', sprites_path,
        '--palette', palette_path,
        '--top', top_name,
        '--left', left_name,
        '--right', right_name,
        '--height', '8',  # Side face height for 32×24 tiles
        '--lighting', 'gradient',
        '--tile-width', '32',  # Final tile width for 32×24 tiles
        '--tile-height', '16',  # Top diamond height for 32×24 tiles
        # --show-outline flag removed (set to false/default)
        '--name', tile_name,
        '--output', output_path
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True)
    
    if result.returncode != 0:
        raise RuntimeError(f"tile2iso failed for {tile_name}: {result.stderr}")
    
    # Load the generated isometric tile
    with open(output_path, 'r') as f:
        output_data = json.load(f)
    
    return output_data[tile_name]

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    engine_root = os.path.join(script_dir, '../..')
    
    # Build tile2iso if needed
    tile2iso_path = build_tile2iso()
    if not tile2iso_path:
        return 1
    
    print(f"Using tile2iso: {tile2iso_path}")
    print()
    
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
    
    # Create temporary directory for intermediate files
    temp_dir = tempfile.mkdtemp()
    try:
        # Create palette.json
        palette_path = os.path.join(temp_dir, 'palette.json')
        create_palette_json(palette_path)
        
        # Map of terrain names to their base names (excluding transitions)
        terrain_tiles = {}
        transition_tiles = {}
        
        for png_file in sorted(png_files):
            # Extract tile name (remove _16x16.png suffix)
            tile_name = png_file.replace('_16x16.png', '')
            png_path = os.path.join(TILES_DIR, png_file)
            
            # Skip transition tiles for now (we'll use base tiles for isometric conversion)
            if '_to_' in tile_name:
                transition_tiles[tile_name] = png_path
                continue
            
            terrain_tiles[tile_name] = png_path
        
        print("Converting tiles to isometric 2.5D...")
        print()
        
        # Create sprites.json with all terrain sprites
        all_sprites = {}
        for tile_name, png_path in terrain_tiles.items():
            sprites = create_sprites_from_png(png_path, tile_name)
            if sprites:
                all_sprites.update(sprites)
        
        # Save sprites.json
        sprites_path = os.path.join(temp_dir, 'sprites.json')
        with open(sprites_path, 'w') as f:
            json.dump(all_sprites, f, indent=2)
        
        # Convert each terrain tile to isometric
        tileset = {}
        for tile_name in sorted(terrain_tiles.keys()):
            print(f"  Converting {tile_name}...")
            try:
                iso_tile = convert_tile_to_isometric(
                    tile2iso_path, sprites_path, palette_path,
                    tile_name, temp_dir
                )
                
                # Convert to TileData format (remove isUI, lifetime, maxSpawn)
                tile_data = {
                    'width': iso_tile['width'],
                    'height': iso_tile['height'],
                    'type': iso_tile['type'],
                    'pixels': iso_tile.get('pixels', []),
                    'frames': iso_tile.get('frames', []),
                    'animations': iso_tile.get('animations', []),
                    'useCollision': iso_tile.get('useCollision', False),
                    'mountPoints': iso_tile.get('mountPoints', [])
                }
                
                tileset[tile_name] = tile_data
                print(f"    ✓ Created {tile_name} ({iso_tile['width']}x{iso_tile['height']})")
                
            except Exception as e:
                print(f"    ✗ Failed to convert {tile_name}: {e}")
                return 1
        
        # Create final tileset with isISO flag
        output_tileset = {
            'isISO': True,
            'tiles': tileset
        }
        
        # Save tileset
        output_path = os.path.join(script_dir, 'assets', 'terrain_tiles.json')
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        
        with open(output_path, 'w') as f:
            json.dump(output_tileset, f, indent=2)
        
        print()
        print(f"✓ Created isometric tileset: {output_path}")
        print(f"  Tiles: {', '.join(sorted(tileset.keys()))}")
        print(f"  Palette: RetroForge 50")
        print(f"  Isometric: True (tileset-level flag)")
        print(f"  All tiles are 2.5D isometric tiles with side faces")
        print()
        
        return 0
        
    finally:
        # Clean up temporary directory
        shutil.rmtree(temp_dir, ignore_errors=True)

if __name__ == '__main__':
    exit(main())

