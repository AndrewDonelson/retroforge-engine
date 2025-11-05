#!/usr/bin/env python3
"""
Generate isometric tileset using tile2iso converter.

This script:
1. Creates three 2D sprites (top, left, right) for each terrain type
2. Uses the tile2iso CLI tool to convert them into 2.5D isometric tiles
3. Outputs a tileset with isISO=true flag

The tile2iso converter performs these steps:
1. Transform top face: Rotate 45° CW + scale Y by 50% → diamond shape
2. Scale side faces: Match desired height (16px)
3. Apply lighting: Gradient lighting to rectangular side faces (before skewing)
4. Transform side faces: Shear into parallelograms (left: +shear, right: -shear)
5. Composite: Layer left side (back) → right side (middle) → top face (front)
"""

import json
import os
import subprocess
import tempfile
import shutil

# RetroForge 50 palette indices for terrain tiles
PALETTE_INDICES = {
    'earth': 7,       # Orange shadow (brownish)
    'grass': 15,      # Green base
    'desert': 9,      # Yellow base
    'snow': 1,        # White
    'water': 22,      # Sky blue shadow
    'lava': 3,        # Red base
    'rock': 4,        # Red shadow (dark)
}

def create_sprite_pixels(width, height, color_idx):
    """Create a simple sprite filled with a color."""
    pixels = []
    for y in range(height):
        row = []
        for x in range(width):
            row.append(color_idx)
        pixels.append(row)
    return pixels

def create_pattern_sprite(width, height, base_color_idx, pattern_func=None):
    """Create a sprite with optional pattern."""
    pixels = []
    for y in range(height):
        row = []
        for x in range(width):
            if pattern_func:
                row.append(pattern_func(x, y, base_color_idx))
            else:
                row.append(base_color_idx)
        pixels.append(row)
    return pixels

def create_dithered_sprite(width, height, base_color_idx, dither_intensity=1):
    """Create a sprite with dithering pattern for more realistic appearance."""
    import random
    pixels = []
    # Create a seed based on base color for consistent dithering
    random.seed(base_color_idx * 1000)
    
    # Define nearby palette indices for dithering (within +/- 2 indices)
    darker_idx = max(0, base_color_idx - dither_intensity)
    lighter_idx = min(49, base_color_idx + dither_intensity)
    
    for y in range(height):
        row = []
        for x in range(width):
            # Use checkerboard dithering with variation
            if (x + y) % 2 == 0:
                # Checkerboard pattern - use base or slightly different shade
                if random.random() < 0.4:
                    row.append(darker_idx if random.random() < 0.5 else lighter_idx)
                else:
                    row.append(base_color_idx)
            else:
                # Alternate squares - sometimes use different shade
                if random.random() < 0.3:
                    row.append(lighter_idx if random.random() < 0.5 else darker_idx)
                else:
                    row.append(base_color_idx)
        pixels.append(row)
    return pixels

def create_terrain_sprites(tile_name, color_idx):
    """Create three 2D sprites (top, left, right) for a terrain type with dithering."""
    # Top face: 16x16 square (top-down view) with dithering
    top_pixels = create_dithered_sprite(16, 16, color_idx, dither_intensity=1)
    
    # Add material-specific patterns for more realistic appearance
    if tile_name == 'grass':
        # Add grass texture with darker spots
        for y in range(16):
            for x in range(16):
                pattern = (x * 3 + y * 7) % 11
                if pattern == 0:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Lighter green spots
                elif pattern == 1:
                    top_pixels[y][x] = max(0, color_idx - 1)  # Darker green spots
    elif tile_name == 'water':
        # Add wave-like pattern with more variation
        for y in range(16):
            for x in range(16):
                wave = int((x * 2 + y * 3) / 3) % 3
                if wave == 0:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Lighter water
                elif wave == 2:
                    top_pixels[y][x] = max(0, color_idx - 1)  # Darker water
    elif tile_name == 'rock':
        # Add rocky texture with cracks
        for y in range(16):
            for x in range(16):
                pattern = (x * 5 + y * 3) % 7
                if pattern < 2:
                    top_pixels[y][x] = max(0, color_idx - 1)  # Dark cracks
                elif pattern == 3:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Light highlights
    elif tile_name == 'desert':
        # Add sand texture with wind patterns
        for y in range(16):
            for x in range(16):
                pattern = (x + y * 2) % 5
                if pattern == 0:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Lighter sand
                elif pattern == 4:
                    top_pixels[y][x] = max(0, color_idx - 1)  # Darker sand
    elif tile_name == 'snow':
        # Add snow texture with subtle variations
        for y in range(16):
            for x in range(16):
                if (x * 7 + y * 11) % 13 == 0:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Slight highlights
    elif tile_name == 'lava':
        # Add lava texture with glowing spots
        for y in range(16):
            for x in range(16):
                pattern = (x * 3 + y * 5) % 8
                if pattern < 2:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Bright spots
                elif pattern == 3:
                    top_pixels[y][x] = max(0, color_idx - 1)  # Darker areas
    elif tile_name == 'earth':
        # Add earth texture with dirt clumps
        for y in range(16):
            for x in range(16):
                pattern = (x * 2 + y * 3) % 6
                if pattern == 0:
                    top_pixels[y][x] = min(49, color_idx + 1)  # Lighter dirt
                elif pattern == 5:
                    top_pixels[y][x] = max(0, color_idx - 1)  # Darker clumps
    
    # Left side: 16x16 rectangle (side view, will be scaled to height)
    # Use slightly darker shade for side with dithering
    left_color = max(0, color_idx - 1) if color_idx > 0 else color_idx
    left_pixels = create_dithered_sprite(16, 16, left_color, dither_intensity=1)
    
    # Right side: 16x16 rectangle (side view, will be scaled to height)
    right_color = max(0, color_idx - 1) if color_idx > 0 else color_idx
    right_pixels = create_dithered_sprite(16, 16, right_color, dither_intensity=1)
    
    return {
        'top': {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': top_pixels,
            'useCollision': False,
            'mountPoints': [],
            'isUI': False,
            'lifetime': 0,
            'maxSpawn': 0
        },
        'left': {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': left_pixels,
            'useCollision': False,
            'mountPoints': [],
            'isUI': False,
            'lifetime': 0,
            'maxSpawn': 0
        },
        'right': {
            'width': 16,
            'height': 16,
            'type': 'static',
            'pixels': right_pixels,
            'useCollision': False,
            'mountPoints': [],
            'isUI': False,
            'lifetime': 0,
            'maxSpawn': 0
        }
    }

def find_tile2iso_binary():
    """Find the tile2iso binary."""
    # Try common locations
    possible_paths = [
        os.path.join(os.path.dirname(__file__), '..', '..', 'bin', 'tile2iso'),
        os.path.join(os.path.dirname(__file__), '..', '..', 'tile2iso'),
        shutil.which('tile2iso'),
    ]
    
    for path in possible_paths:
        if path and os.path.exists(path) and os.access(path, os.X_OK):
            return path
    
    # Try building it
    engine_dir = os.path.join(os.path.dirname(__file__), '..', '..')
    if os.path.exists(os.path.join(engine_dir, 'cmd', 'tile2iso')):
        print("Building tile2iso tool...")
        result = subprocess.run(
            ['go', 'build', '-o', 'bin/tile2iso', './cmd/tile2iso'],
            cwd=engine_dir,
            capture_output=True,
            text=True
        )
        if result.returncode == 0:
            return os.path.join(engine_dir, 'bin', 'tile2iso')
        else:
            print(f"Warning: Failed to build tile2iso: {result.stderr}")
    
    return None

def create_palette_json(palette_path):
    """Create a palette.json file for RetroForge 50."""
    # Load the actual palette.json from the engine root
    engine_dir = os.path.join(os.path.dirname(__file__), '..', '..')
    actual_palette_path = os.path.join(engine_dir, 'palette.json')
    
    if os.path.exists(actual_palette_path):
        # Copy the actual palette file
        shutil.copy(actual_palette_path, palette_path)
        return
    
    # Fallback: RetroForge 50 palette colors (exactly 50 colors)
    palette = {
        "colors": [
            "#000000", "#ffffff", "#ff8989", "#ff4d4d", "#c31111",  # 0-4
            "#ffcd89", "#ff914d", "#c35511", "#ffff89", "#ffd84d",  # 5-9
            "#c39c11", "#f2ff89", "#b6ff4d", "#7ac311", "#89ffc3",  # 10-14
            "#4dd487", "#11984b", "#72ffff", "#36d8c7", "#009c8b",  # 15-19
            "#89ffff", "#4dd5ff", "#1199c3", "#a2fbff", "#66bfff",  # 20-24
            "#2a83c3", "#abc4ff", "#6f88ff", "#334cc3", "#c6b1ff",  # 25-29
            "#8a75ff", "#4e39c3", "#f0b4ff", "#b478ff", "#783cc3",  # 30-34
            "#ffabed", "#ff6fb1", "#c33375", "#ffbbdc", "#ff7fa0",  # 35-39
            "#c34364", "#e4b596", "#a8795a", "#6c3d1e", "#dced96",  # 40-44
            "#a0b15a", "#64751e", "#74f9ff", "#38bdf8", "#0081bc"   # 45-49
        ]
    }
    
    with open(palette_path, 'w') as f:
        json.dump(palette, f, indent=2)

def convert_tile_to_isometric(tile2iso_path, sprites_path, palette_path, tile_name, top_name, left_name, right_name, output_dir):
    """Convert three sprites to an isometric tile using tile2iso CLI."""
    output_path = os.path.join(output_dir, f"{tile_name}_iso.json")
    
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
    
    # Read the output JSON
    with open(output_path, 'r') as f:
        output_data = json.load(f)
    
    # Extract the sprite data for this tile
    if tile_name in output_data:
        return output_data[tile_name]
    else:
        raise RuntimeError(f"Tile {tile_name} not found in output: {output_data.keys()}")

def main():
    tileset_dir = os.path.dirname(__file__)
    assets_dir = os.path.join(tileset_dir, 'assets')
    
    # Find or build tile2iso binary
    tile2iso_path = find_tile2iso_binary()
    if not tile2iso_path:
        print("ERROR: Could not find or build tile2iso binary")
        print("Please build it manually: cd retroforge-engine && go build -o bin/tile2iso ./cmd/tile2iso")
        return 1
    
    print(f"Using tile2iso: {tile2iso_path}")
    
    # Create temporary directory for conversion
    with tempfile.TemporaryDirectory() as temp_dir:
        sprites_path = os.path.join(temp_dir, 'temp_sprites.json')
        palette_path = os.path.join(temp_dir, 'palette.json')
        
        # Create palette.json
        create_palette_json(palette_path)
        
        # Create tileset
        tileset = {}
        tiles = ['earth', 'grass', 'desert', 'snow', 'water', 'lava', 'rock']
        
        # Step 1: Create all input sprites
        all_sprites = {}
        for tile_name in tiles:
            color_idx = PALETTE_INDICES[tile_name]
            sprites = create_terrain_sprites(tile_name, color_idx)
            
            # Add to sprite map with unique names
            all_sprites[f"{tile_name}_top"] = sprites['top']
            all_sprites[f"{tile_name}_left"] = sprites['left']
            all_sprites[f"{tile_name}_right"] = sprites['right']
        
        # Save temporary sprites.json
        with open(sprites_path, 'w') as f:
            json.dump(all_sprites, f, indent=2)
        
        # Step 2: Convert each tile using tile2iso
        print("\nConverting tiles to isometric 2.5D...")
        for tile_name in tiles:
            print(f"  Converting {tile_name}...")
            try:
                top_name = f"{tile_name}_top"
                left_name = f"{tile_name}_left"
                right_name = f"{tile_name}_right"
                
                # Call tile2iso converter
                iso_tile = convert_tile_to_isometric(
                    tile2iso_path, sprites_path, palette_path,
                    tile_name, top_name, left_name, right_name, temp_dir
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
        
        # Step 3: Create tileset data structure
        tileset_data = {
            'isISO': True,
            'tiles': tileset
        }
        
        # Write tileset file
        output_path = os.path.join(assets_dir, 'terrain_tiles.json')
        with open(output_path, 'w') as f:
            json.dump(tileset_data, f, indent=2)
        
        print(f"\n✓ Created isometric tileset: {output_path}")
        print(f"  Tiles: {', '.join(tiles)}")
        print(f"  Palette: RetroForge 50")
        print(f"  Isometric: True (tileset-level flag)")
        print(f"  All tiles are 2.5D isometric tiles with side faces")
    
    return 0

if __name__ == "__main__":
    exit(main())
