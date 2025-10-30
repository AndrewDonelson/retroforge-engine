package app

import "sync/atomic"

var quitFlag int32

func RequestQuit() { atomic.StoreInt32(&quitFlag, 1) }
func QuitRequested() bool { return atomic.LoadInt32(&quitFlag) == 1 }
func Reset() { atomic.StoreInt32(&quitFlag, 0) }


