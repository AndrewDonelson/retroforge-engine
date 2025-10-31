# Tron Light Cycles (RetroForge example)

A classic Tron-style light cycles game showcasing the RetroForge engine:
- Procedural enemy placement per level
- Increasing difficulty with speed and trail length scaling
- Simple AI for enemy light cycles
- Progressive enemy count (new enemy every 5 levels)

## Gameplay

You control the **Blue Light Cycle** starting from a random position at the bottom of the arena, always moving upward initially. Your goal is to force enemy light cycles to crash into your trail. If you crash into any trail (yours or theirs), you lose!

- **Enemy Colors**: Red, Green, Yellow, Orange, Cyan, Purple
- **Difficulty Scaling**: Every level increases speed and trail length
- **Enemy Progression**: A new enemy is added every 5 levels (up to 6 total)

## Controls

- **Arrow Keys**: Turn your light cycle (Left/Right/Up/Down)
- **O/X/Enter**: Menu selection
- **ESC**: Quit

## Build and Run

```bash
cd /home/andrew/Development/Golang/RetroForge/retroforge-engine
make pack-tron
make run CART=examples/tron-lightcycles.rf
```

## Bundle a Self-Contained Executable

```bash
make bundle CART=examples/tron-lightcycles.rf OUT=tron-lightcycles
./cart-tron-lightcycles -scale 3
```

## Notes

- Packed carts use the `.rf` extension
- The game uses a grid-based system for collision detection
- Enemy AI uses simple pathfinding with randomness for unpredictability
- Difficulty scales over 50 levels

