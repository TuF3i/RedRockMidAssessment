package timer

import (
	"RedRockMidAssessment-Synchronizer/core"
	"time"
)

var Span = 5 * time.Minute

func Timer() {
	ticker := time.NewTicker(Span)
	defer ticker.Stop()
	for {
		select {
		case <-core.TimerStop:
			return
		case <-ticker.C:
			core.TaskQ <- struct{}{}
		}
	}
}
