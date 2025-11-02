#!/usr/bin/env python3
"""
Generate RetroForge logo sprite data from SVG description.
The icon is a 48x48 box with:
- Dark blue background (#0f172a) = palette index 16
- Cyan border (#0369a1, closest to #0081bc) = palette index 49  
- Dotted white line at y=12
- Cyan "RF" text centered (palette index 48)
- Solid white line at y=42
"""

WIDTH = 48
HEIGHT = 48

# Palette indices
TRANSPARENT = -1
BLACK = 0
WHITE = 1
DARK_BLUE = 16  # #0f172a
CYAN_BLUE = 48  # #38bdf8
DARK_CYAN = 49  # #0081bc (closest to #0369a1 border)

# Create empty sprite (all transparent)
sprite = [[TRANSPARENT for _ in range(WIDTH)] for _ in range(HEIGHT)]

# Icon box: centered, 24x24 pixels (half size)
# Offset to center: (48-24)/2 = 12
box_x = 12
box_y = 12
box_w = 24
box_h = 24

# Draw border (outer 2 pixels)
for y in range(box_y, box_y + box_h):
    for x in range(box_x, box_x + box_w):
        if y == box_y or y == box_y + box_h - 1 or x == box_x or x == box_x + box_w - 1:
            sprite[y][x] = DARK_CYAN  # Border

# Draw background
for y in range(box_y + 2, box_y + box_h - 2):
    for x in range(box_x + 2, box_x + box_w - 2):
        sprite[y][x] = DARK_BLUE

# Draw dotted line at y=12 (top line, inside box)
for x in range(box_x + 2, box_x + box_w - 2):
    if (x - box_x - 2) % 4 < 2:  # Dotted pattern
        sprite[box_y + 2][x] = WHITE

# Draw "RF" text (simplified pixel art)
# R starts at x=14, y=18, width=8, height=12
# F starts at x=24, y=18, width=7, height=12
rx = box_x + 2
ry = box_y + 6
rw = 7
rh = 10

# Letter R
for y in range(ry, ry + rh):
    for x in range(rx, rx + rw):
        # Vertical left
        if x == rx:
            sprite[y][x] = CYAN_BLUE
        # Top horizontal
        elif y == ry and x < rx + rw - 1:
            sprite[y][x] = CYAN_BLUE
        # Middle horizontal
        elif y == ry + rh//2 - 1 and x < rx + rw - 1:
            sprite[y][x] = CYAN_BLUE
        # Vertical right (top half)
        elif x == rx + rw - 2 and y < ry + rh//2:
            sprite[y][x] = CYAN_BLUE
        # Diagonal (bottom right)
        elif x >= rx + rw - 3 and y >= ry + rh//2 and (x - rx - rw + 3) == (y - ry - rh//2):
            sprite[y][x] = CYAN_BLUE

# Letter F
fx = rx + rw + 1
fy = ry
fw = 6
fh = 10

for y in range(fy, fy + fh):
    for x in range(fx, fx + fw):
        # Vertical left
        if x == fx:
            sprite[y][x] = CYAN_BLUE
        # Top horizontal
        elif y == fy and x < fx + fw - 1:
            sprite[y][x] = CYAN_BLUE
        # Middle horizontal
        elif y == fy + fh//2 - 1 and x < fx + fw - 2:
            sprite[y][x] = CYAN_BLUE

# Draw solid line at y=42 (bottom line, inside box)
for x in range(box_x + 2, box_x + box_w - 2):
    sprite[box_y + box_h - 3][x] = WHITE

# Print Go code
print("var logoPixels = [][]int{")
for y in range(HEIGHT):
    row = ", ".join(f"{pixel:3d}" for pixel in sprite[y])
    print(f"\t{{ {row} }},")
print("}")

