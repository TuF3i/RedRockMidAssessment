package kafka

import (
	"RedRockMidAssessment-Consumer/core"
	"errors"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
)

func ConnectToKafka() (sarama.ConsumerGroup, error) {
	// 生成配置
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	//cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	//cfg.Version = sarama.V4_0_0_0

	// 连接Broker
	dst := fmt.Sprintf("%v:%v", core.Config.Mq.Kafka.Addr, core.Config.Mq.Kafka.Port) // 生成连接地址
	if strings.TrimSpace(core.Config.Mq.Kafka.GroupID) == "" {
		return nil, errors.New("kafka group id must not be empty")
	} // 空ID检查
	group, err := sarama.NewConsumerGroup([]string{dst}, core.Config.Mq.Kafka.GroupID, cfg)
	if err != nil {
		return nil, err
	}

	//core.Group = group
	return group, nil
}
