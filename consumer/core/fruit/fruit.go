package fruit

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/dao"
	"RedRockMidAssessment-Consumer/core/kafka"
	viper "RedRockMidAssessment-Consumer/core/utils/config"
	zap "RedRockMidAssessment-Consumer/core/utils/log"
	"RedRockMidAssessment-Consumer/core/utils/snowflake"
	"RedRockMidAssessment-Consumer/core/worker"
	"fmt"
	"os"
	"time"

	"gitee.com/liumou_site/logger"
)

const (
	CONFIG_PATH = "./data/config/config.yaml"
	LOG_PATH    = "./data/logs"
)

func GenesisFruit() {
	//初始化日志
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
	generator, err := snowflake.NewSnowflake(core.Config.SFID.MachineID)
	if err != nil {
		logs.Warn("Init mod <snowflake> error: %v", err.Error())
		os.Exit(1)
	}
	core.SnowFlake = generator
	logs.Info("Successfully loaded mod <snowflake>")

	// 初始化MySQL连接
	logs.Debug("Started to init mod <MySQL>")
	mysqlConn, err := dao.ConnectToMySQL(true, true)
	if err != nil {
		logs.Warn("Init mod <MySQL> error: %v", err.Error())
		os.Exit(1)
	}
	core.MysqlConn = mysqlConn
	logs.Info("Successfully loaded mod <MySQL>")

	// 初始化全局WaitGroup
	worker.InitGlobalWg()

	// 初始化kafka消费者组
	logs.Debug("Started to init mod <kafka-consumer-group>")
	cg, err := kafka.ConnectToKafka()
	if err != nil {
		logs.Warn("Init mod <kafka-consumer-group> error: %v", err.Error())
		os.Exit(1)
	}
	core.Group = cg
	logs.Info("Successfully loaded mod <kafka-consumer-group>")

	// 启动消费者组
	logs.Debug("Started to Setup <kafka-consumer-group-thread>")
	cancel := kafka.StartWorkingThread()
	core.Cancel = cancel
	logs.Info("Successfully Setup <kafka-consumer-group-thread>")

	// 完成提示
	fmt.Println()
	logs.Alert("Server Started Successfully, present logs: ")
	core.Logger.Debug("Debug")
}

func WorldEndingFruit() {
	//初始化日志
	logs := logger.NewLogger(1)
	logs.Modular = "WorldEndingFruit"

	// 发送取消消息（暂时没用，先留着）
	logs.Info("Sending Cancel() Message...")
	core.Cancel()
	time.Sleep(time.Second)

	// 等待全部ACK
	logs.Info("Waiting for consumer...")
	if core.GlobalWg != nil {
		worker.WaitGlobalWg()
	}

	// 关闭底层连接
	logs.Debug("Started to clean mod <kafka-consumer>")
	err := core.Group.Close()
	if err != nil {
		logs.Warn("Cleaning mod <kafka-consumer> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <kafka-consumer>")

	// 关闭MySQL连接
	//（其实并不需要，gorm.db不持有socket）

	// 清除日志缓存
	logs.Debug("Started to clean mod <zap>")
	if err := core.Logger.Sync(); err != nil {
		logs.Warn("Cleaning mod <zap> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <zap>")
}
