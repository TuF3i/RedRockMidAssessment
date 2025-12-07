package kafka

import (
	"RedRockMidAssessment/core"
	"fmt"

	"github.com/IBM/sarama"
)

func NewProducer() (sarama.SyncProducer, error) {
	// 配置
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewManualPartitioner
	cfg.Producer.Return.Successes = true

	// 构造生产者
	dsn := fmt.Sprintf("%v:%v", core.Config.Mq.Kafka.Addr, core.Config.Mq.Kafka.Port)
	producer, err := sarama.NewSyncProducer([]string{dsn}, cfg)
	if err != nil {
		return nil, err
	}
	//defer producer.Close()
	return producer, nil
}
