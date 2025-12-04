package worker

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/flitter"
	"RedRockMidAssessment-Consumer/core/models"
	"context"
	"encoding/json"
	"strconv"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Worker struct {
	blanket chan struct{} // 令牌桶（用于并发控制）
}

func NewConsumer(max int) *Worker {
	return &Worker{blanket: make(chan struct{}, max)}
}

func (c *Worker) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Worker) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Worker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 过滤topic
		if msg.Topic != "Consumer" {
			return nil
		}
		// 拿令牌
		c.blanket <- struct{}{}
		// 异步处理
		go func(m *sarama.ConsumerMessage) {
			var commander models.Commander
			// 释放令牌
			defer func() { <-c.blanket }()
			// 生成snowflakeID
			traceID := core.SnowFlake.TraceID()
			ctx := context.WithValue(context.Background(), "TraceID", traceID)
			// 解析JSON
			if err := json.Unmarshal(m.Value, &commander); err != nil {
				core.Logger.Error(
					"Unmarshal JSON Byte Error",
					zap.String("snowflake", traceID),
					zap.String("detail", err.Error()),
					zap.String("raw_message", string(m.Value)),
				)
			}
			// 调用flitter生成业务方法
			biz := flitter.GetRelatedHandleFunc(commander)
			// 调用业务方法
			biz.Do(ctx, commander)
			// 输出处理日志
			core.Logger.Debug(
				"Handle Message With Offset",
				zap.String("snowflake", traceID),
				zap.String("offset", strconv.FormatInt(msg.Offset, 10)),
				zap.String("raw_message", string(m.Value)),
			)
			// ACK
			sess.MarkMessage(m, "")
		}(msg)
	}
	return nil
}
