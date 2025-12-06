package kafka

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/worker"
	"context"
	"log"
)

func work(ctx context.Context, h *worker.Worker) {
	for {
		// 一旦重平衡或连接断，sarama 会自动重试
		if err := core.Group.Consume(ctx, []string{core.CONSUME_TOPIC}, h); err != nil {
			log.Printf("consume error: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func StartWorkingThread() context.CancelFunc {
	// 构造消费者
	h := worker.NewWorker(core.Config.Mq.Kafka.BlanketPeek)
	// 上下文
	ctx, cancel := context.WithCancel(context.Background())
	// 启动工作携程
	go work(ctx, h)

	return cancel
}
