//go:build !(js && wasm)

package input

// logResetStepFlag is a no-op for non-WASM builds
func logResetStepFlag() {}

// logStep is a no-op for non-WASM builds
func logStep() {}

// logSet is a no-op for non-WASM builds
func logSet(i int, down bool, frameInput bool, oldCur bool, newCur bool) {}

// logBtn is a no-op for non-WASM builds
func logBtn(i int, result bool, cur bool, prev bool) {}

// logBtnp is a no-op for non-WASM builds
func logBtnp(i int, cur bool, prev bool, stepped bool) {}

