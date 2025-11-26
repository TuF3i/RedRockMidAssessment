package core

import (
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/snowflake"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger *zap.Logger
	Config *models.Config

	MysqlConn *gorm.DB
	RedisConn *redis.Client
	SnowFlake *snowflake.Snowflake
)
