package worker

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/service"
	"context"
	"sync"
)

// 任务结构体
type task struct {
	courseID string
	ctx      context.Context
}

// 任务群组
type WorkerGroup struct {
	innerTask chan task
	w         sync.WaitGroup
}

func RunWorker() {
	// 生成线程池对象
	WorkerPoll := WorkerGroup{innerTask: make(chan task, 2*core.Config.Mq.Kafka.BlanketPeak), w: sync.WaitGroup{}}
	// 更新全局指针
	core.GlobalWg = &WorkerPoll.w
	// 启动工作线程
	WorkerPoll.WakeUpWorker()
	//发布工作任务
	go WorkerPoll.PublishWork() //这个放在后台吧
	//等待任务完成
	WorkerPoll.w.Wait()
}

func (w *WorkerGroup) Worker() {
	defer w.w.Done() // 任务计数器递减
	for t := range w.innerTask {
		// 进行业务操作
		service.UpdateSelectedStuList(t.ctx, t.courseID)
		service.UpdateDroppedStuList(t.ctx, t.courseID)
		service.UpdateCourseSelectedNum(t.ctx, t.courseID)
	}
}

func (w *WorkerGroup) PublishWork() {
	// 生成TraceID
	traceID := core.SnowFlake.TraceID()
	ctx := context.WithValue(context.Background(), "TraceID", traceID)
	// 获取所有课程ID
	ids := service.GetAllCourseID(ctx)
	if ids == nil { // 空ID则跳过
		return
	}
	// 类型断言
	i := ids.([]string)
	// 遍历发布任务
	for _, id := range i {
		w.innerTask <- task{courseID: id, ctx: ctx}
	}
}

func (w *WorkerGroup) WakeUpWorker() {
	for i := 0; i < core.Config.Mq.Kafka.BlanketPeak; i++ {
		go w.Worker()
		w.w.Add(1)
	}
}
