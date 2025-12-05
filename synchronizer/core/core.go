package core

import (
	"RedRockMidAssessment-Synchronizer/core/models"
	"RedRockMidAssessment-Synchronizer/core/utils/snowflake"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Config            *models.Config           // 同步节点配置
	Logger            *zap.Logger              // 同步节点日志核心
	SnowFlake         snowflake.Snowflake      // SnowFlakeID生成器
	RedisConn         *redis.Client            // Redis连接
	Producer          sarama.AsyncProducer     // kafka生产者
	PartitionConsumer sarama.PartitionConsumer // kafka消息接收器
	TimerStop         chan struct{}            // 计时器停止信号
	TaskQ             chan struct{}            // 全局任务队列

	READ_TOPIC        string = "Sync"
	WRITE_TOPIC       string = "Consumer"
	DEFAULT_PARTITION int32  = 0

	TimerStatus = 0 // 计时器状态
)
