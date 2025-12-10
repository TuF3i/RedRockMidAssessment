package core

import (
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/snowflake"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"gorm.io/gorm"
)

var (
	Logger *hertzzap.Logger // 日志组件
	Config *models.Config   // 配置结构体

	MysqlConn *gorm.DB             // MySQL连接
	RedisConn *redis.Client        // Redis连接
	SnowFlake *snowflake.Snowflake // snowflake生成器
	Producer  sarama.SyncProducer  // kafka生产者

	WRITE_TOPIC             = "Sync" // kafka发送的默认主题
	DEFAULT_PARTITION int32 = 0      // kafka默认分区
)
