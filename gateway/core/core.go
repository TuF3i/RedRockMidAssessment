package core

import (
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/snowflake"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger *zap.Logger    // 日志组件
	Config *models.Config // 配置结构体

	MysqlConn *gorm.DB             // MySQL连接
	RedisConn *redis.Client        // Redis连接
	SnowFlake *snowflake.Snowflake // snowflake生成器
)
