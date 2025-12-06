package kafka

import (
	"RedRockMidAssessment-Synchronizer/core"
	"fmt"

	"github.com/IBM/sarama"
)

func NewProducer() (sarama.AsyncProducer, error) {
	// 配置
	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewManualPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	// 构造生产者
	dsn := fmt.Sprintf("%v:%v", core.Config.Mq.Kafka.Addr, core.Config.Mq.Kafka.Port)
	producer, err := sarama.NewAsyncProducer([]string{dsn}, cfg)
	if err != nil {
		return nil, err
	}
	//defer producer.Close()
	return producer, nil
}

func NewConsumer() (sarama.PartitionConsumer, error) {
	// 生成配置
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	// 构造消费者
	dsn := fmt.Sprintf("%v:%v", core.Config.Mq.Kafka.Addr, core.Config.Mq.Kafka.Port)
	consumer, err := sarama.NewConsumer([]string{dsn}, cfg)
	if err != nil {
		return nil, err
	}

	// 构造消息接收器
	pc, err := consumer.ConsumePartition(core.READ_TOPIC, core.DEFAULT_PARTITION, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	//core.PartitionConsumer = pc
	return pc, nil
}
