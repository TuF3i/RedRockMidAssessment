package timerlistener

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/worker"
	"context"
	"time"
)

func Listener(ctx context.Context) {
	for {
		select {
		case <-core.TaskQ:
			worker.RunWorker()
			time.Sleep(5 * time.Second)
		case <-ctx.Done():
			return
		}
	}
}
