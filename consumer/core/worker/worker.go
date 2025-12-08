package worker

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/flitter"
	"RedRockMidAssessment-Consumer/core/models"
	"context"
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// 消费者组
type Worker struct {
	blanket chan struct{} // 并发度控制令牌桶
}

// 创建全局WaitGroup
func InitGlobalWg() {
	core.GlobalWg = &sync.WaitGroup{}
}

// 等待消费完成
func WaitGlobalWg() {
	core.GlobalWg.Wait()
}

// 新建消费者组
func NewWorker(maxConcurrency int) *Worker {
	w := Worker{blanket: make(chan struct{}, maxConcurrency)}
	return &w
}

// 消费方法
func (c *Worker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 过滤topic
		if msg.Topic != core.CONSUME_TOPIC {
			continue
		}
		// 心跳检测
		if string(msg.Value) == "ping" {
			core.Logger.Info("pong")
			sess.MarkMessage(msg, "")
			continue
		}
		// 忽略非法消息
		if len(msg.Value) < 2 || msg.Value[0] != '{' && msg.Value[0] != '[' {
			sess.MarkMessage(msg, "")
			continue
		}
		// 拿令牌
		c.blanket <- struct{}{}
		core.GlobalWg.Add(1)
		// 异步处理
		go func(data []byte) {
			var commander models.Commander
			// 释放令牌
			defer func() {
				<-c.blanket
				core.GlobalWg.Done()
			}()
			// 生成snowflakeID
			traceID := core.SnowFlake.TraceID()
			ctx := context.WithValue(context.Background(), "TraceID", traceID)
			// 解析JSON
			if err := json.Unmarshal(data, &commander); err != nil {
				core.Logger.Error(
					"Unmarshal JSON Byte Error",
					zap.String("snowflake", traceID),
					zap.String("detail", err.Error()),
					zap.String("raw_message", string(data)),
				)
				return
			}
			// 调用flitter生成业务方法
			biz := flitter.GetRelatedHandleFunc(commander)
			// 调用业务方法
			biz.Do(ctx, commander)
			// 输出处理日志
			core.Logger.Debug(
				"Handle Message With Offset",
				zap.String("snowflake", traceID),
				zap.String("raw_message", string(data)),
			)
		}(msg.Value)
		// ACK
		sess.MarkMessage(msg, "")
	}
	return nil
}

// 初始化消费者组方法（无实际意义，仅满足接口）
func (c *Worker) Setup(sarama.ConsumerGroupSession) error { return nil }

// 释放化消费者组方法（无实际意义，仅满足接口）
func (c *Worker) Cleanup(sarama.ConsumerGroupSession) error { return nil }
