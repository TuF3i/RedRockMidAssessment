package fruit

import (
	"RedRockMidAssessment-Consumer/core/dao"
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/dao/mysql"
	viper "RedRockMidAssessment/core/utils/config"
	zap "RedRockMidAssessment/core/utils/log"
	"RedRockMidAssessment/core/utils/snowflake"
	"os"

	"gitee.com/liumou_site/logger"
)

const (
	CONFIG_PATH = "./data/config/config.yaml"
	LOG_PATH    = "./data/log"
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
}
