package worker

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/service"
	"context"
	"sync"
)

type WorkerGroup struct {
	innerTask  chan string
	bufferSize uint
	w          sync.WaitGroup
}

func (w *WorkerGroup) Worker() {
	for courseID := range w.innerTask {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx := context.WithValue(context.Background(), "TraceID", traceID)
		// 进行业务操作
		selectedStu := service.GetSelectedStuList(ctx, courseID)

	}
}
