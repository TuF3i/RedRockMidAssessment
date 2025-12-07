package kafka

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/utils/response"
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func Writer(ctx context.Context, jsonBuffer []byte) response.Response {
	msg := &sarama.ProducerMessage{
		Topic:     core.WRITE_TOPIC,
		Partition: core.DEFAULT_PARTITION,
		Value:     sarama.ByteEncoder(jsonBuffer),
	}

	_, _, err := core.Producer.SendMessage(msg)
	if err != nil {
		core.Logger.Error(
			"Send Kafka Message Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}

	return response.OperationSuccess
}
