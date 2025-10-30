# Moon Lander (RetroForge example)

A tiny lunar-landing demo showcasing the RetroForge engine:
- Deterministic terrain per level
- Pad placement and scoring
- Stars, HUD, and simple SFX/music (tokens)

Controls
- Left/Right: rotate
- Up: thrust
- O/X/Enter: menu select
- ESC: quit

Build and run
```
cd /home/andrew/Development/Golang/RetroForge/retroforge-engine
make pack-moon
make run CART=examples/moon-lander.rf
```

Bundle a self-contained executable
```
make bundle CART=examples/moon-lander.rf OUT=moon-lander
./cart-moon-lander -scale 3
```

Notes
- Packed carts now use the `.rf` extension.
- You can compose simple tunes with `rf.music({"4C1","4E1","4G2"}, 120, 0.3)`.
