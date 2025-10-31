# RetroForge - Project Naming & Branding

**Final Decision:** October 29, 2025  
**Version:** 1.0

---

## Project Name: RetroForge

### Why "RetroForge"?

**Retro**
- Captures the nostalgic, classic gaming aesthetic
- References 8-bit and 16-bit era constraints
- Appeals to indie developers and retro enthusiasts
- Immediately communicates the project's vintage inspiration

**Forge**
- Emphasizes the creative, crafting aspect of game development
- Suggests strength, durability, and quality
- Conveys active creation ("forging" games)
- Implies a tool/workshop environment

**Combined Benefits:**
- ✅ Memorable and easy to pronounce
- ✅ Descriptive without being too literal
- ✅ Appeals to both creators and players
- ✅ Professional yet approachable
- ✅ Good for branding (logo possibilities)
- ✅ Domain availability likely (retroforge.dev, retroforge.com)

---

## File Extensions

### `.rfs` - RetroForge Source
**Usage:** Development and editing
- **RFS** = **R**etro**F**orge **S**ource
- Uncompressed ZIP archive
- Easy to inspect and debug
- Version control friendly
- Editable in the web app
- Can be shared for collaboration

**Example:**
```
my-game.rfs
└── (uncompressed cart contents)
```

### `.rfe` - RetroForge Executable
**Usage:** Distribution and sharing
- **RFE** = **R**etro**F**orge **E**xecutable
- Compressed ZIP archive
- Digitally signed
- Smaller file size
- Ready for distribution
- Can be run directly in browser or desktop app

**Example:**
```
my-game.rfe
└── (compressed, signed cart contents)
```

---

## Brand Identity

### Logo Concepts
**Primary Elements:**
- Anvil or hammer (forging metaphor)
- Pixel art aesthetic
- Retro color palette
- Simple, iconic design

**Color Scheme Ideas:**
- Warm metals: Bronze, copper, gold tones
- Classic retro: Cyan, magenta, yellow, black
- Forge fire: Orange, red, yellow gradients with dark background

### Typography
- **Primary Font:** Pixel/bitmap font for retro feel
- **Secondary Font:** Modern sans-serif for readability
- Consider: Press Start 2P, VT323, or custom pixel font

### Tagline Options
- "Forge Your Retro Dreams"
- "Where Pixels Meet Passion"
- "Craft Classic Games"
- "Forge. Play. Share."
- "Modern Tools. Retro Soul."

---

## Domain Strategy

### Primary Domains (Recommended)
- **retroforge.dev** - Main site (developers)
- **retroforge.com** - Alternative/marketing site
- **play.retroforge.dev** - Web player
- **docs.retroforge.dev** - Documentation
- **api.retroforge.dev** - API if needed

### Social Media Handles
- Twitter/X: @retroforge or @retroforge_dev
- GitHub: github.com/retroforge
- Discord: RetroForge Community
- Reddit: r/retroforge

---

## Naming Conventions

### Repository Names
```
retroforge-engine       # Go engine (WASM + native)
retroforge-web          # Next.js web application
retroforge-docs         # Documentation site
retroforge-examples     # Example carts
retroforge-community    # Community resources
```

### Package Names
**Go:**
```go
github.com/retroforge/engine
github.com/retroforge/engine/lua
github.com/retroforge/engine/graphics
github.com/retroforge/engine/audio
```

**NPM:**
```
@retroforge/web
@retroforge/types
@retroforge/sdk
```

### URL Structure
```
retroforge.dev                    # Landing page
retroforge.dev/create             # Web editor
retroforge.dev/play               # Browse carts
retroforge.dev/play/[cart-id]     # Play specific cart
retroforge.dev/learn              # Tutorials
docs.retroforge.dev               # Documentation
docs.retroforge.dev/api           # API reference
docs.retroforge.dev/guides        # Guides
```

---

## File Type Associations

### MIME Types
```
application/x-retroforge-source      # .rfs
application/x-retroforge-executable  # .rfe
```

### Desktop Integration
**macOS:**
- Application: RetroForge.app
- Associated with .rfs and .rfe files
- Custom icons for both file types

**Windows:**
- Application: RetroForge.exe
- File associations in registry
- Custom file icons

**Linux:**
- Application: retroforge
- Desktop entry with MIME types
- Custom file icons via desktop theme

---

## Marketing Messages

### For Developers
- "Create retro games with modern tools"
- "No complex engine—just pure creativity"
- "Finish your game in a weekend"
- "Learn game development the fun way"

### For Players
- "Discover new retro games daily"
- "Play in your browser—no downloads"
- "Support indie creators"
- "Nostalgia meets innovation"

### For Educators
- "Teach programming through game development"
- "Perfect for coding workshops"
- "Safe, educational environment"
- "Real projects, instant results"

---

## Competitive Positioning

### Compared to PICO-8
**RetroForge advantages:**
- ✅ Free and open source
- ✅ Modern web-based editor
- ✅ Higher resolution (320×240 vs 128×128)
- ✅ More colors (50 vs 16)
- ✅ More code tokens (16k vs 8k)
- ✅ 8 layers with parallax
- ✅ Built for web from the start

### Compared to TIC-80
**RetroForge advantages:**
- ✅ Better web integration
- ✅ Modern TypeScript web app
- ✅ Cloud cart repository built-in
- ✅ More curated experience
- ✅ Better for non-programmers (visual editors)

### Unique Selling Points
1. **Perfect balance** between constraints and capability
2. **Web-first** design (WASM from day one)
3. **Beautiful color system** (50 color palettes based on color theory)
4. **Non-programmer friendly** (visual piano roll, sprite editor)
5. **Modern tech stack** (Go, WASM, Next.js, TypeScript)
6. **Open source philosophy** while protecting creators

---

## Community Building

### Launch Strategy
1. **Week 1-2:** Build in public on Twitter/X
2. **Week 3-4:** Share early prototypes
3. **Month 2:** Invite alpha testers
4. **Month 3:** Launch beta with example carts
5. **Month 4:** Public launch with game jam

### Community Channels
- **Discord:** Main community hub
  - #general
  - #showcase (finished games)
  - #help (support)
  - #development (engine development)
  - #feedback
- **GitHub Discussions:** Technical discussions
- **Twitter/X:** Updates and showcases
- **YouTube:** Tutorial videos
- **Itch.io:** Alternative cart hosting

### Content Strategy
- Weekly dev logs
- Tutorial series (beginner to advanced)
- Featured cart of the week
- Developer spotlights
- Monthly game jams with themes

---

## Legal Considerations

### Trademark
- Consider trademark registration for "RetroForge"
- Check existing trademarks in gaming/software space
- File in relevant jurisdictions (US, EU)

### Open Source License
**Recommended:** MIT License
- Permissive
- Allows commercial use
- Simple and widely understood
- Good for community growth

**Alternative:** Apache 2.0
- Patent protection
- More explicit rights granted

### Cart Licensing
- Default: Carts are open source (code visible)
- Optional: Creator can choose license per cart
- Encourage: Creative Commons for assets

---

## Version Naming

### Engine Versions
- **Format:** Major.Minor.Patch (Semantic Versioning)
- **Example:** RetroForge Engine v1.0.0

### API Versions
- **Format:** v1, v2, etc.
- **Stability:** v1 stable, v2 beta, etc.

### Cart Format Versions
- **Format:** RFS/RFE format version in manifest
- **Example:** "format_version": "1.0"

---

## Success Metrics

### Launch Goals (3 months)
- 100+ registered users
- 50+ published carts
- 10+ active community members
- 1000+ cart plays

### Year 1 Goals
- 1,000+ registered users
- 500+ published carts
- Active Discord community (100+ members)
- 10,000+ cart plays
- First game jam completed

### Long-term Vision
- Premier platform for retro game development
- Thriving creator community
- Educational adoption (schools, workshops)
- Sustainable open-source project
- Potential commercial hosting options

---

## Next Steps

### Immediate Actions
1. ✅ **Finalize name** → RetroForge (DONE)
2. ✅ **Choose file extensions** → .rfs and .rfe (DONE)
3. 🔲 **Check domain availability**
   - retroforge.dev
   - retroforge.com
4. 🔲 **Create social media accounts**
   - Twitter/X
   - GitHub organization
5. 🔲 **Design initial logo**
   - Sketch concepts
   - Get community feedback
6. 🔲 **Set up repositories**
   - retroforge-engine
   - retroforge-web
7. 🔲 **Build landing page**
   - Explain vision
   - Collect early interest emails

### Development Phase
1. Engine MVP (Go + gopher-lua)
2. Web app MVP (Next.js + TypeScript)
3. WASM compilation pipeline
4. First example cart
5. Alpha testing with 5-10 users

### Marketing Phase
1. Launch on Product Hunt
2. Post on HackerNews
3. Share on Reddit (r/gamedev, r/indiedev)
4. Reach out to retro gaming YouTubers
5. Host launch game jam

---

## Conclusion

**RetroForge** represents the perfect blend of nostalgia and innovation. The name captures the creative, crafting nature of game development while honoring the retro aesthetic that inspires us. With thoughtful branding, clear file conventions, and a strong community focus, RetroForge is positioned to become a beloved platform for retro game creators worldwide.

**Mission:** Empower creators to forge amazing retro games with modern tools and an amazing community.

**Vision:** A thriving ecosystem where anyone can create, share, and enjoy retro games.

**Values:**
- 🎨 Creativity over complexity
- 🤝 Community over competition  
- 📚 Learning by doing
- 🔓 Open by default
- ⚡ Finish your projects
