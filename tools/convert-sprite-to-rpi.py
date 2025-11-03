#!/usr/bin/env python3
"""
Convert a sprite from sprites.json to .rpi (Raw Palette Indexed) format.

.rpi format:
- Header (8 bytes):
  - width: uint16 (2 bytes, little-endian)
  - height: uint16 (2 bytes, little-endian)
  - flags: uint16 (2 bytes) - reserved for future use
  - reserved: uint16 (2 bytes)
- Data: Packed 6-bit palette indices (0-49 for colors, 255 = transparent)
  - Each pixel uses 6 bits (supports 0-63, but we use 0-49 + 255 mapping)
  - Pixels are packed sequentially, row by row
  - Transparent is encoded as 63 (6 bits max), then mapped to -1 on decode

Usage:
  python3 convert-sprite-to-rpi.py <sprites.json> <sprite_name> <output.rpi>
"""

import json
import struct
import sys
import gzip

def convert_sprite_to_rpi(sprites_file, sprite_name, output_file, compress=True):
    # Load sprites.json
    with open(sprites_file, 'r') as f:
        sprites = json.load(f)
    
    if sprite_name not in sprites:
        raise ValueError(f"Sprite '{sprite_name}' not found in {sprites_file}")
    
    sprite = sprites[sprite_name]
    width = sprite['width']
    height = sprite['height']
    pixels = sprite['pixels']
    
    print(f"Converting {sprite_name}: {width}x{height}")
    
    # Pack header (8 bytes)
    # flags: bit 0 = landscape (0) / portrait (1), bits 1-15 reserved
    flags = 0
    if height > width:
        flags |= 1  # Portrait mode
    
    header = struct.pack('<HHHH', width, height, flags, 0)
    
    # Pack pixel data as 6-bit values
    # Map: -1 (transparent) -> 63, 0-49 -> 0-49
    # Pixels are stored row by row: pixels[height][width]
    encoded_pixels = []
    
    for y in range(height):
        if y >= len(pixels):
            # Row missing, fill with transparent
            encoded_pixels.extend([63] * width)
            continue
        row = pixels[y]
        for x in range(width):
            if x >= len(row):
                # Column missing, use transparent
                encoded = 63
            else:
                pixel = row[x]
                # Map transparent (-1) to 63, colors (0-49) stay as-is
                if pixel == -1:
                    encoded = 63
                elif 0 <= pixel <= 49:
                    encoded = pixel
                else:
                    # Invalid palette index, default to transparent
                    encoded = 63
            encoded_pixels.append(encoded)
    
    # Pack 6-bit values into bytes
    # Pack 4 pixels into 3 bytes: 4 * 6 = 24 bits = 3 bytes
    packed_data = bytearray()
    for i in range(0, len(encoded_pixels), 4):
        # Get up to 4 pixels
        p = encoded_pixels[i:i+4]
        while len(p) < 4:
            p.append(0)  # Pad with 0 if needed
        
        # Pack 4 pixels (24 bits) into 3 bytes:
        # Pixel 0: bits 0-5 -> Byte 0 bits 0-5
        # Pixel 1: bits 0-5 -> Byte 0 bits 6-7, Byte 1 bits 0-3
        # Pixel 2: bits 0-5 -> Byte 1 bits 4-7, Byte 2 bits 0-1
        # Pixel 3: bits 0-5 -> Byte 2 bits 2-7
        b0 = (p[0] & 0x3F) | ((p[1] & 0x03) << 6)
        b1 = ((p[1] >> 2) & 0x0F) | ((p[2] & 0x0F) << 4)
        b2 = ((p[2] >> 4) & 0x03) | ((p[3] & 0x3F) << 2)
        
        packed_data.extend([b0, b1, b2])
    
    # Combine header + data
    rpi_data = header + bytes(packed_data)
    
    # Compress if requested
    if compress:
        rpi_data = gzip.compress(rpi_data, compresslevel=9)
    
    # Write output
    with open(output_file, 'wb') as f:
        f.write(rpi_data)
    
    original_size = len(header) + (width * height * 6 + 7) // 8
    compressed_size = len(rpi_data)
    
    print(f"Original (uncompressed): {original_size:,} bytes ({original_size/1024:.1f} KB)")
    print(f"Compressed (.rpi): {compressed_size:,} bytes ({compressed_size/1024:.1f} KB)")
    print(f"Compression ratio: {compressed_size/original_size*100:.1f}%")
    print(f"Saved to: {output_file}")

if __name__ == '__main__':
    if len(sys.argv) != 4:
        print("Usage: python3 convert-sprite-to-rpi.py <sprites.json> <sprite_name> <output.rpi>")
        sys.exit(1)
    
    convert_sprite_to_rpi(sys.argv[1], sys.argv[2], sys.argv[3])

