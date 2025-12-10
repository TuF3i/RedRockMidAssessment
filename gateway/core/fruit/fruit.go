package fruit

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/api"
	"RedRockMidAssessment/core/dao/mysql"
	"RedRockMidAssessment/core/dao/redis"
	"RedRockMidAssessment/core/kafka"
	"RedRockMidAssessment/core/models"
	viper "RedRockMidAssessment/core/utils/config"
	zap "RedRockMidAssessment/core/utils/log"
	"RedRockMidAssessment/core/utils/md5"
	"RedRockMidAssessment/core/utils/snowflake"
	"fmt"
	"os"

	"gitee.com/liumou_site/logger"
)

const (
	CONFIG_PATH = "./data/config/config.yaml"
	LOG_PATH    = "./data/logs"
)

func GenesisFruit() {
	// 初始化日志
	logs := logger.NewLogger(1)
	logs.Modular = "GenesisFruit"

	// 初始化配置文件
	logs.Debug("Started to init mod <viper>")
	conf, err := viper.InitConfig(CONFIG_PATH)
	if err != nil {
		logs.Warn("Init mod <viper> error: %v", err.Error())
		os.Exit(1)
	}
	core.Config = conf
	logs.Info("Successfully loaded mod <viper>")

	// 初始化日志
	logs.Debug("Started to init mod <zap>")
	zapCore := zap.InitZap(LOG_PATH)
	core.Logger = zapCore
	logs.Info("Successfully loaded mod <zap>")

	// 初始化SnowFlake
	logs.Debug("Started to init mod <snowflake>")
	generator, err := snowflake.NewSnowflake(core.Config.SnowFlake.MachineID)
	if err != nil {
		logs.Warn("Init mod <snowflake> error: %v", err.Error())
		os.Exit(1)
	}
	core.SnowFlake = generator
	logs.Info("Successfully loaded mod <snowflake>")

	// 初始化MySQL连接
	logs.Debug("Started to init mod <MySQL>")
	mysqlConn, err := mysql.ConnectToMySQL(true, false)
	if err != nil {
		logs.Warn("Init mod <MySQL> error: %v", err.Error())
		os.Exit(1)
	}
	core.MysqlConn = mysqlConn
	logs.Info("Successfully loaded mod <MySQL>")

	// 初始化Redis连接
	logs.Debug("Started to init mod <redis>")
	redisConn, err := redis.ConnectToRedis()
	if err != nil {
		logs.Warn("Init mod <redis> error: %v", err.Error())
		os.Exit(1)
	}
	core.RedisConn = redisConn
	logs.Info("Successfully loaded mod <redis>")

	// 初始化kafka生产者
	logs.Debug("Started to init mod <kafka-producer>")
	producer, err := kafka.NewProducer()
	if err != nil {
		logs.Warn("Init mod <kafka-producer> error: %v", err.Error())
		os.Exit(1)
	}
	core.Producer = producer
	logs.Info("Successfully loaded mod <kafka-producer>")

	// 添加Root用户
	logs.Debug("Started to add user <admin>")
	err = CreateAdminUser()
	if err != nil {
		logs.Warn("Init mod <kafka-producer> error: %v", err.Error())
		os.Exit(1)
	}
	logs.Info("Successfully added user <admin>")

	// 启动HertzAPI
	api.HertzApi()
	// 测试日志
	core.Logger.Info(fmt.Sprintf("DEBUG"))

}

func WorldEndingFruit() {
	// 初始化日志
	logs := logger.NewLogger(1)
	logs.Modular = "WorldEndingFruit"

	// 关闭kafka生产者
	logs.Debug("Started to clean mod <kafka-producer>")
	if err := core.Producer.Close(); err != nil {
		logs.Warn("Cleaning mod <kafka-producer> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <kafka-producer>")

	// 关闭Redis连接
	logs.Debug("Started to clean mod <redis>")
	if err := core.RedisConn.Close(); err != nil {
		logs.Warn("Cleaning mod <redis> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <redis>")

	// 关闭MySQL连接
	//（其实并不需要，gorm.db不持有socket）

	// 清除日志缓存
	logs.Debug("Started to clean mod <zap>")
	core.Logger.Sync()
	logs.Info("Successfully cleaned mod <zap>")

}

func CreateAdminUser() error {
	// 用户信息
	userForm := models.Student{
		Role:         false,
		Name:         "admin",
		StudentID:    "1234567890",
		StudentClass: "adminGroup",
		Password:     md5.GenMD5("P@ssw0rd=Ping12345"),
		Sex:          0,
		Grade:        4,
		Age:          19,
	}
	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	// 检测用户是否存在
	result := tx.Where("student_id = ?", userForm.StudentID).Find(&models.Student{})
	if err := result.Error; err != nil {
		tx.Rollback()
		return err
	}
	// 判断用户是否存在
	if result.RowsAffected != 0 {
		tx.Rollback()
		return nil
	}
	// 添加用户
	if err := tx.Create(&userForm).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	tx.Commit()
	return nil
}
