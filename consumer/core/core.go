package core

import (
	"RedRockMidAssessment-Consumer/core/models"
	"RedRockMidAssessment-Consumer/core/utils/snowflake"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config    *models.Config       // 消费者节点配置
	Logger    *zap.Logger          // 消费者节点日志核心
	Group     sarama.ConsumerGroup // 消费者组
	MysqlConn *gorm.DB             // MySQL数据库连接
	SnowFlake snowflake.Snowflake  // SnowFlakeID生成器

	GlobalWg *sync.WaitGroup // 全局WaitGroup

	CONSUME_TOPIC = "Consumer"
)
