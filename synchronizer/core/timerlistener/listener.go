package timerlistener

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/worker"
	"time"
)

func Listener() {
	for {
		select {
		case <-core.TaskQ:
			worker.RunWorker()
			time.Sleep(5 * time.Second)
		}
	}
}
