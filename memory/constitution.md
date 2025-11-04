# RetroForge Engine - Project Constitution

**Version:** 1.0  
**Date:** October 30, 2025  
**Project:** RetroForge Engine (Go Runtime)

---

## 🎯 Project Mission

**"Create the most accessible and powerful fantasy console for retro game development"**

RetroForge Engine provides a modern, cross-platform runtime for creating retro-style games with the creative constraints that encourage finished projects while offering more capability than traditional 8-bit consoles.

---

## 🏗️ Core Principles

### 1. **Creative Constraints Drive Innovation**
- **Limited resources** (16,384 tokens, 256 sprites, 50 colors) spark creativity
- **Fixed resolution** (480×270) ensures consistent experience
- **Simple APIs** reduce complexity and encourage experimentation
- **Performance limits** (60 FPS, 64MB WASM) force optimization

### 2. **Modern Tools, Retro Soul**
- **Go runtime** for cross-platform compatibility
- **Node system** for modern game architecture
- **Box2D physics** for realistic interactions
- **Lua scripting** for accessible programming
- **Retro aesthetic** with pixel-perfect rendering

### 3. **Developer Experience First**
- **Dual API** (high-level nodes + low-level direct)
- **Comprehensive documentation** with examples
- **Clear error messages** and debugging tools
- **Fast iteration** with hot reloading
- **Cross-platform** development and deployment

### 4. **Community and Learning**
- **Open source** by default
- **Learn by example** through readable cart code
- **Educational focus** for teaching programming
- **Community sharing** of carts and knowledge
- **No barriers** to entry or participation

---

## 🎮 Technical Standards

### Performance Requirements
- **60 FPS** stable on all target platforms
- **<16.67ms** frame budget (60 FPS)
- **<64MB** WASM memory usage
- **<2 seconds** cart loading time
- **<20ms** audio latency

### Code Quality Standards
- **Go 1.23+** with modern idioms
- **Comprehensive testing** (unit + integration)
- **Performance benchmarks** for critical paths
- **Memory profiling** and optimization
- **Cross-platform** compatibility testing

### API Design Principles
- **Consistent naming** across all APIs
- **Lua-friendly** function signatures
- **Clear documentation** with examples
- **Backward compatibility** for stable APIs
- **Progressive enhancement** for new features

### State Machine System
- **All example carts MUST use module-based state system** (`rf.import()`)
- **State modules require these functions:**
  - `_INIT()` - Called once when state is first created
  - `_UPDATE(dt)` - Called every frame (dt is delta time in seconds)
  - `_DRAW()` - Called every frame
  - `_HANDLE_INPUT()` - Called every frame (before update)
  - `_DONE()` - Called once when state is destroyed
- **Optional functions:**
  - `_ENTER()` - Called every time state becomes active
  - `_EXIT()` - Called every time state becomes inactive
- **State file naming:** `play_state.lua` → state name is `"play"`
- **Engine startup flow (AUTOMATIC, handled by engine):**
  1. **First state:** Engine splash screen (2 seconds or any input to skip)
  2. **Second state:** After engine splash:
     - If `splash_state.lua` exists → transition to `"splash"` state
     - Else → transition to `initialState` (defaults to `"menu"`)
  3. **NO manual state transitions needed** - engine handles everything
- **IMPORTANT:** 
  - If you only have one state, name it `menu_state.lua` so it registers as `"menu"` and matches the engine's default
  - Import states directly in `main.lua` (not in `_INIT()` function)
  - Engine automatically transitions after splash - do NOT call `game.changeState()` manually at startup
- **Button API:** `rf.btn(index)` and `rf.btnp(index)` use **numeric indices 0-10**, NOT string names
  - Button indices: 0=SELECT, 1=START, 2=UP, 3=DOWN, 4=LEFT, 5=RIGHT, 6=A, 7=B, 8=X, 9=Y, 10=TURBO
- **Movement/continuous actions** go in `_UPDATE(dt)` using `dt` parameter
- **Button presses** (`rf.btnp()`) go in `_HANDLE_INPUT()`
- **Button holds** (`rf.btn()`) can go in either `_HANDLE_INPUT()` or `_UPDATE(dt)`
- **Example structure:**
  ```lua
  -- main.lua
  local play_state = rf.import("play_state.lua")
  
  -- play_state.lua
  function _INIT() end
  function _ENTER() end
  function _HANDLE_INPUT()
    if rf.btnp("a") then -- Instant actions
  end
  function _UPDATE(dt)
    if rf.btn("up") then -- Continuous actions using dt
  end
  function _DRAW() end
  function _EXIT() end
  function _DONE() end
  ```

---

## 🚀 Development Philosophy

### Spec-Driven Development
- **Constitution first** - this document guides all decisions
- **Specifications** before implementation
- **User stories** drive feature development
- **Testing** as part of the specification process
- **Iterative refinement** based on feedback

### Quality Over Speed
- **Thorough testing** before release
- **Performance optimization** as a first-class concern
- **Documentation** updated with every change
- **Code review** required for all changes
- **User feedback** incorporated regularly

### Open and Transparent
- **Public development** process
- **Open source** licensing (MIT)
- **Community input** on major decisions
- **Regular updates** on progress
- **Honest communication** about challenges

---

## 🎯 Success Metrics

### Technical Goals
- ✅ **60 FPS** on all target platforms
- ✅ **<64MB** memory usage for typical carts
- ✅ **<2 seconds** cart loading time
- ✅ **Zero crashes** in stable releases
- ✅ **100% API** test coverage

### User Experience Goals
- ✅ **<5 minutes** to first running cart
- ✅ **Clear documentation** for all features
- ✅ **Helpful error messages** for common issues
- ✅ **Smooth performance** across all platforms
- ✅ **Intuitive APIs** for common tasks

### Community Goals
- ✅ **Active community** of developers
- ✅ **High-quality** example carts
- ✅ **Educational** content and tutorials
- ✅ **Regular** community events and jams
- ✅ **Positive** developer feedback

---

## 🚫 What We Don't Do

### Anti-Patterns to Avoid
- **Over-engineering** simple features
- **Breaking changes** without migration path
- **Platform lock-in** or vendor dependencies
- **Complex configuration** or setup
- **Performance regressions** in stable releases

### Features We Won't Add
- **True 3D rendering** (2.5D raytracing only)
- **Complex shader systems** (simple materials only)
- **Real-time networking** (local multiplayer only)
- **Cloud storage** (local-first approach)
- **Proprietary formats** (open standards only)

---

## 🔄 Decision Making Process

### Major Decisions
1. **Constitution alignment** - Does this fit our principles?
2. **Community input** - What do users think?
3. **Technical feasibility** - Can we implement it well?
4. **Performance impact** - Does it meet our standards?
5. **Documentation** - Can we explain it clearly?

### Implementation Process
1. **Specification** - Detailed technical spec
2. **Prototype** - Proof of concept
3. **Testing** - Comprehensive test suite
4. **Documentation** - Complete API docs
5. **Release** - Stable, well-tested feature

---

## 📚 Learning and Growth

### Continuous Improvement
- **Regular retrospectives** on development process
- **User feedback** analysis and incorporation
- **Performance monitoring** and optimization
- **Community engagement** and support
- **Technical debt** management and reduction

### Knowledge Sharing
- **Documentation** as a first-class deliverable
- **Example code** for all features
- **Tutorial content** for learning paths
- **Community contributions** recognition
- **Open source** knowledge sharing

---

## 🎉 Project Values

### Core Values
- **🎨 Creativity** - Constraints spark innovation
- **🤝 Community** - Together we build better
- **📚 Learning** - Education through creation
- **🔓 Openness** - Transparent and accessible
- **⚡ Performance** - Fast and efficient
- **🎮 Fun** - Making games should be enjoyable

### Quality Standards
- **Reliability** - Stable and predictable
- **Usability** - Easy to learn and use
- **Performance** - Fast and efficient
- **Maintainability** - Clean and well-documented
- **Extensibility** - Easy to add new features

---

## 🚀 Future Vision

### Short Term (6 months)
- **Core engine** complete and stable
- **Node system** fully implemented
- **Physics integration** working smoothly
- **Audio system** with all features
- **Cross-platform** builds working

### Medium Term (1 year)
- **Voxel raytracing** system (Phase 2)
- **Performance optimizations** complete
- **Comprehensive documentation** finished
- **Community** actively contributing
- **Example games** showcasing capabilities

### Long Term (2+ years)
- **Industry standard** for fantasy consoles
- **Educational adoption** in schools
- **Thriving ecosystem** of tools and games
- **Cross-platform** mobile support
- **Advanced features** based on community needs

---

**This constitution serves as the foundation for all RetroForge Engine development decisions. When in doubt, refer to these principles and values to guide the way forward.**

---

*"Forge Your Retro Dreams" - RetroForge Engine Constitution* 🔨✨
