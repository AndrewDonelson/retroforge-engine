//go:build js && wasm

package input

import (
	"syscall/js"
)

// logResetStepFlag logs when ResetStepFlag() is called (WASM only)
func logResetStepFlag() {
	// No-op now - ResetStepFlag doesn't clear cur anymore
}

// logStep logs when Step() is called (WASM only)
func logStep() {
	console := js.Global().Get("console")
	if console.Truthy() {
		if js.Global().Get("rf_step_log_count").IsUndefined() {
			js.Global().Set("rf_step_log_count", 0)
		}
		logCount := js.Global().Get("rf_step_log_count").Int()
		js.Global().Set("rf_step_log_count", logCount+1)
		if logCount < 10 || logCount%60 == 0 {
			console.Call("log", "[Input] Step() called - prev = cur, clearing cur only if prev=false")
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

// logBtnp logs a Btnp() call when it returns true (WASM only)
func logBtnp(i int, cur bool, prev bool, stepped bool) {
	console := js.Global().Get("console")
	if console.Truthy() {
		console.Call("log", "[Input] Btnp(", i, ") = TRUE! cur:", cur, "prev:", prev, "stepped:", stepped)
	}
}

