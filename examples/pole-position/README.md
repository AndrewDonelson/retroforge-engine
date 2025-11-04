# Pole Position

A classic arcade racing game built with RetroForge Engine.

## Features

- **10 Procedural Tracks**: Each track is procedurally generated with unique curves, widths, and difficulty
- **3D Perspective Racing**: Classic arcade-style perspective rendering
- **Traffic AI**: Race against 7 AI opponents
- **3-Lap Races**: Complete 3 laps to finish
- **Full State System**: Splash, Menu, Track Select, Play, Pause, and Results states
- **Dynamic Car Rendering**: Cars rendered as 3x3 composite sprites using primitives
- **Music & SFX**: Background music and sound effects

## Controls

- **UP**: Accelerate
- **DOWN**: Brake
- **LEFT/RIGHT**: Steer
- **SELECT**: Pause (during race)
- **START/A**: Confirm selection

## Gameplay

- Race through 10 different tracks with increasing difficulty
- Avoid collisions with traffic cars
- Complete 3 laps as fast as possible
- Track your best lap times

## Technical Details

- **Procedural Track Generation**: Each track uses seeded random generation for consistent replayability
- **Perspective Rendering**: Classic arcade-style 3D perspective using 2D primitives
- **Composite Sprites**: Cars are rendered as 9 sprites (3x3 grid) using drawing primitives
- **Physics**: Speed, acceleration, steering, and collision detection
- **Traffic AI**: Simple AI for opponent cars with lane changes and speed variation

## Assets

- **sprites.json**: Car sprite definitions (rendered with primitives)
- **sfx.json**: Engine loop, crash, shift, checkpoint, and countdown sounds
- **music.json**: Menu music and 3 different race music tracks

