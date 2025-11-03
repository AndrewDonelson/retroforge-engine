//go:build js && wasm

package input

import (
	"syscall/js"
)

//export rf_set_shift
func rf_set_shift(down bool) {
	// TURBO button is now mapped to shift key
	Set(BtnTURBO, down)
}

//export rf_set_button
func rf_set_button(name string, down bool) {
	SetByName(name, down)
}

//export rf_set_btn
func rf_set_btn(idx int, down bool) {
	oldCur := cur[idx]
	oldPrev := prev[idx]
	Set(idx, down)
	newCur := cur[idx]
	
	// Log button state changes for debugging (always log to see timing)
	console := js.Global().Get("console")
	if console.Truthy() {
		console.Call("log", "[Input.Set] Button", idx, ":", oldCur, "->", newCur, "prev:", oldPrev)
	}
}

// logResetStepFlag logs when ResetStepFlag() is called (WASM only)
func logResetStepFlag() {
	// No-op now - ResetStepFlag doesn't clear cur anymore
}

// logStep logs when Step() is called (WASM only)
func logStep() {
	console := js.Global().Get("console")
	if console.Truthy() {
		// Always log first few times, then periodically
		if js.Global().Get("rf_step_log_count").IsUndefined() {
			js.Global().Set("rf_step_log_count", 0)
		}
		logCount := js.Global().Get("rf_step_log_count").Int()
		js.Global().Set("rf_step_log_count", logCount+1)
		if logCount < 10 || logCount%60 == 0 {
			// Log button states before Step() to see what's being copied
			// Convert Go slice to JavaScript array
			jsArray := js.Global().Get("Array").New()
			for i := 0; i < 11; i++ {
				if cur[i] {
					jsArray.Call("push", i)
				}
			}
			console.Call("log", "[Input] Step() called - copying cur to prev. Active buttons before Step:", jsArray)
		}
	}
}

// logSet logs a Set() call (WASM only)
func logSet(i int, down bool, frameInput bool, oldCur bool, newCur bool) {
	if js.Global().Get("rf_set_log_count").IsUndefined() {
		js.Global().Set("rf_set_log_count", 0)
	}
	logCount := js.Global().Get("rf_set_log_count").Int()
	js.Global().Set("rf_set_log_count", logCount+1)
	if logCount < 10 || (down && logCount%30 == 0) {
		console := js.Global().Get("console")
		if console.Truthy() {
			console.Call("log", "[Input] Set(", i, ",", down, ") frameInput:", frameInput, "cur:", oldCur, "->", newCur)
		}
	}
}

// logBtn logs a Btn() call (WASM only)
func logBtn(i int, result bool, cur bool, prev bool) {
	if js.Global().Get("rf_btn_log_count").IsUndefined() {
		js.Global().Set("rf_btn_log_count", 0)
	}
	logCount := js.Global().Get("rf_btn_log_count").Int()
	js.Global().Set("rf_btn_log_count", logCount+1)
	if logCount < 10 || logCount%60 == 0 {
		console := js.Global().Get("console")
		if console.Truthy() {
			console.Call("log", "[Input] Btn(", i, ") =", result, "cur:", cur, "prev:", prev)
		}
	}
}

// logBtnp logs Btnp() calls only when a button press is detected (WASM only)
func logBtnp(i int, cur bool, prev bool, stepped bool) {
	result := cur && !prev
	// Only log when button press is actually detected (result = true)
	// This prevents spam when buttons are checked but not pressed
	if result && i <= 5 {
		console := js.Global().Get("console")
		if console.Truthy() {
			console.Call("log", "[Input.Btnp] Button", i, "PRESSED ✓ cur:", cur, "prev:", prev)
		}
	}
}

// logBtnpFalse logs when Btnp() returns false but cur is true (potential timing issue)
func logBtnpFalse(i int, cur bool, prev bool) {
	if i <= 5 && cur {
		// Button is currently pressed but btnp returned false (prev must also be true)
		console := js.Global().Get("console")
		if console.Truthy() {
			console.Call("warn", "[Input] Btnp(", i, ") = FALSE but cur=true, prev=", prev, " - edge missed!")
		}
	}
}
