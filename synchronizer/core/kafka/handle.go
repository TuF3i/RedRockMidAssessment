package kafka

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/models"
	"RedRockMidAssessment-Synchronizer/core/timer"
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func Writer(jsonBuffer []byte) {
	msg := &sarama.ProducerMessage{
		Topic:     core.WRITE_TOPIC,
		Partition: core.DEFAULT_PARTITION,
		Value:     sarama.ByteEncoder(jsonBuffer),
	}
	core.Producer.Input() <- msg
}

func KLogger(ctx context.Context) {
	for {
		select {
		case suc := <-core.Producer.Successes():
			core.Logger.Info(
				">>> Async Sent Successful",
				// zap.String("snowflake", traceID),
				zap.Int32("partition", suc.Partition),
				zap.Int64("offset", suc.Offset),
				zap.Any("key", suc.Key),
				zap.Any("value", suc.Value),
			)
		case fail := <-core.Producer.Errors():
			core.Logger.Info(
				">>> Async Sent Failed",
				// zap.String("snowflake", traceID),
				zap.String("Error", fail.Error()),
			)
		case <-ctx.Done():
			return
		}
	}
}

func MsgReader(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			core.Logger.Info("MsgReader received context canceled, exiting")
			return

		case msg, ok := <-core.PartitionConsumer.Messages():
			if !ok {
				// 消费者通道被关闭，直接退出
				core.Logger.Info("PartitionConsumer.Messages channel closed, exiting")
				return
			}

			var commander models.Commander
			// 解析 JSON
			if err := json.Unmarshal(msg.Value, &commander); err != nil {
				core.Logger.Error("Unmarshal JSON Byte Error",
					zap.String("detail", err.Error()),
					zap.String("raw_message", string(msg.Value)),
				)
				return
			}

			// 判断操作类型
			if commander.Role != "CourseSelectionEvent" {
				return
			}

			// 类型断言
			status, ok := commander.Msg.(string)
			if !ok {
				core.Logger.Error("Type Assertion Error",
					zap.Any("detail", commander.Msg),
					zap.String("raw_message", string(msg.Value)),
				)
				return
			}

			// 判断操作
			switch status {
			case "open":
				if core.TimerStatus != 1 {
					core.TimerStatus = 1
					go timer.Timer()
				}
			case "close":
				if core.TimerStatus != 0 {
					core.TimerStatus = 0
					core.TimerStop <- struct{}{}
					core.TaskQ <- struct{}{}
				}
			default:
				core.Logger.Error("Unknown Status", zap.Any("detail", status))
			}
		}
	}
}
