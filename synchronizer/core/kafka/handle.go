package kafka

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/models"
	"RedRockMidAssessment-Synchronizer/core/timer"
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

func KLogger() {
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
		}
	}
}

func MsgReader() {
	for msg := range core.PartitionConsumer.Messages() {
		var commander models.Commander
		// 解析JSON
		if err := json.Unmarshal(msg.Value, &commander); err != nil {
			core.Logger.Error(
				"Unmarshal JSON Byte Error",
				zap.String("detail", err.Error()),
				zap.String("raw_message", string(msg.Value)),
			)
			return
		}
		// 判断操作类型
		if commander.Role != "CourseSelectionEvent" {
			return
		}
		//类型断言
		status, ok := commander.Msg.(string)
		if !ok {
			core.Logger.Error(
				"Type Assertion Error",
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
			core.Logger.Error(
				"Unknow Status",
				zap.Any("detail", status),
			)
		}
	}
}
