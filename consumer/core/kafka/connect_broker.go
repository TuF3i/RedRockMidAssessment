package kafka

import (
	"RedRockMidAssessment-Consumer/core"
	"fmt"

	"github.com/IBM/sarama"
)

func ConnectToKafka() (sarama.ConsumerGroup, error) {
	// 生成配置
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 连接Broker
	dst := fmt.Sprintf("%v:%v", core.Config.Mq.Kafka.Addr, core.Config.Mq.Kafka.Port) // 生成连接地址
	group, err := sarama.NewConsumerGroup([]string{dst}, core.Config.Mq.Kafka.GroupID, cfg)
	if err != nil {
		return nil, err
	}

	//core.Group = group
	return group, nil
}
