// Package tile2iso provides tools for converting 2D top-down sprites into
// isometric (2.5D) tiles by combining three textures: a top face, left side face,
// and right side face. The converter supports multiple lighting modes to create
// depth and dimensionality.
//
// Core Functions:
//   - CreateIsometricTile: Convert three sprites (top, left, right) into an isometric tile
//   - Apply lighting modes: normal, basic, full, gradient
//   - Transform top face to isometric projection
//
// All functions work with RetroForge sprites (from sprites.json) and support
// static, frames, and animation sprite types.
package tile2iso
