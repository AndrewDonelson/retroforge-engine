# RetroForge - Quick Reference

**Project Status:** Design Complete ✅ | Ready for Development 🚀  
**Date:** October 29, 2025

---

## 📦 What is RetroForge?

A modern fantasy console for creating retro-style games with:
- 480×270 resolution (16:9 landscape) or 270×480 (9:16 portrait)
- 50-color palettes (16 base colors × 3 shades + B&W)
- 16,384 Lua code tokens
- 256 sprites (4×4 to 32×32)
- 8-layer maps with parallax
- 8-channel chip-tune audio
- Pure synthesis (square, triangle, sawtooth, sine, noise)
- Box2D physics engine
- Godot-style node system

**Platforms:**
- Windows, macOS, Linux (desktop)
- Android (mobile)
- Web browsers (WASM)
- iOS not supported

---

## 🎯 File Extensions

- **`.rfs`** - RetroForge Source (development)
- **`.rfe`** - RetroForge Executable (distribution)

---

## 🛠️ Tech Stack

**Engine:**
- Go 1.23
- gopher-lua (Lua 5.1)
- SDL2 + OpenGL
- TinyGo for WASM

**Web App:**
- Next.js 16
- TypeScript
- TailwindCSS
- react-ace (code editor)

---

## 📂 Project Structure

```
retroforge-engine/          # Go runtime
retroforge-web/             # Next.js web app  
retroforge-docs/            # Documentation
retroforge-examples/        # Example carts
```

---

## 🎨 Cart Structure

```
game.rfe (ZIP archive)
├── manifest.json
├── main.lua
├── code/
├── levels/
│   ├── level_1.json
│   └── level_2.json
├── sprites.png
├── sprites.json
└── audio.json
```

---

## 📚 Core API Categories

1. **Graphics** - Primitives, sprites, camera
2. **Input** - Keyboard/gamepad
3. **Audio** - SFX, music
4. **Map** - Tilemaps, flags, parallax
5. **Memory** - Persistent storage
6. **Math** - Helpers, random
7. **Utility** - Time, strings, tables

---

## 🗓️ Development Timeline

- **Weeks 1-4:** Engine MVP
- **Weeks 5-7:** Audio + Advanced Nodes
- **Weeks 8-10:** Web App
- **Weeks 11-12:** Polish + Examples
- **Weeks 13-14:** Android + Cross-Platform

---

## 📄 Documentation Files

1. **[RETROFORGE_DESIGN.md](computer:///mnt/user-data/outputs/RETROFORGE_DESIGN.md)**
   - Complete technical specification
   - API reference
   - Architecture details

2. **[RETROFORGE_BRANDING.md](computer:///mnt/user-data/outputs/RETROFORGE_BRANDING.md)**
   - Brand identity
   - Logo concepts
   - Marketing messages
   - Community strategy

3. **[RETROFORGE_KICKOFF.md](computer:///mnt/user-data/outputs/RETROFORGE_KICKOFF.md)**
   - Immediate next steps
   - Development phases
   - Setup instructions
   - Success metrics

4. **[DESIGN_UPDATES_v0.2.md](computer:///mnt/user-data/outputs/DESIGN_UPDATES_v0.2.md)**
   - Version 0.2 changes
   - Design decisions
   - Architecture decision records

---

## ✅ Next Actions

### This Week
1. Purchase domains (retroforge.dev, retroforge.com)
2. Create GitHub organization
3. Set up repositories
4. Define color palette hex values
5. Reserve social media handles

### Next Week
1. Begin engine development (Week 1)
2. Set up SDL2 + gopher-lua
3. Implement basic rendering loop
4. Create first Lua callbacks

---

## 🎮 Example Cart (Hello World)

```lua
function _init()
  x = 240
  y = 135
end

function _update()
  if btn(2) then x = x - 2 end
  if btn(3) then x = x + 2 end
  if btn(0) then y = y - 2 end
  if btn(1) then y = y + 2 end
end

function _draw()
  cls(0)
  circfill(x, y, 10, 8)
  print("Use arrows!", 180, 20, 7)
end
```

---

## 🔗 Useful Links

- Go: https://go.dev
- gopher-lua: https://github.com/yuin/gopher-lua
- SDL2: https://www.libsdl.org
- Next.js: https://nextjs.org
- TinyGo: https://tinygo.org
- PICO-8: https://www.lexaloffle.com/pico-8.php

---

## 📞 Contact & Community

**GitHub:** github.com/retroforge (to be created)  
**Discord:** RetroForge Community (to be created)  
**Twitter/X:** @retroforge (to be created)

---

*Ready to forge some retro magic! 🔨✨*
