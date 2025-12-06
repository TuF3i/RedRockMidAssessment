package fruit

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/dao"
	"RedRockMidAssessment-Synchronizer/core/kafka"
	"RedRockMidAssessment-Synchronizer/core/timerlistener"
	viper "RedRockMidAssessment-Synchronizer/core/utils/config"
	zap "RedRockMidAssessment-Synchronizer/core/utils/log"
	"RedRockMidAssessment-Synchronizer/core/utils/snowflake"
	"fmt"
	"os"
	"time"

	"gitee.com/liumou_site/logger"
)

const (
	CONFIG_PATH = "./data/config/config.yaml"
	LOG_PATH    = "./data/log"
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
	generator, err := snowflake.NewSnowflake(core.Config.Snowflake.MachineID)
	if err != nil {
		logs.Warn("Init mod <snowflake> error: %v", err.Error())
		os.Exit(1)
	}
	core.SnowFlake = generator
	logs.Info("Successfully loaded mod <snowflake>")

	// 初始化Redis连接
	logs.Debug("Started to init mod <redis>")
	redisConn, err := dao.ConnectToRedis()
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

	// 初始化kafka消费者
	pc, err := kafka.NewConsumer()
	logs.Debug("Started to init mod <kafka-consumer>")
	if err != nil {
		logs.Warn("Init mod <kafka-consumer> error: %v", err.Error())
		os.Exit(1)
	}
	core.PartitionConsumer = pc
	logs.Info("Successfully loaded mod <kafka-consumer>")

	// 初始化通信管道
	core.TimerStop = make(chan struct{}) // 计时器停止指令
	core.TaskQ = make(chan struct{}, 5)  // 全局任务队列

	//启动一些监视携程
	go kafka.KLogger()          // 启动kafka日志携程
	go kafka.MsgReader()        // 启动信息监听携程
	go timerlistener.Listener() // 启动计时器信号监听携程

	fmt.Println()
	logs.Alert("Server Started Successfully, present logs: ")
}

func WorldEndingFruit() {
	// 初始化日志
	logs := logger.NewLogger(1)
	logs.Modular = "WorldEndingFruit"

	// 关闭消费者,等待消息处理完毕
	logs.Debug("Started to clean mod <kafka-consumer>")
	err := core.PartitionConsumer.Close()
	if err != nil {
		logs.Warn("Cleaning mod <kafka-consumer> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <kafka-consumer>")

	// 进行关闭前的一次同步
	go func() { // 防止死锁
		if core.TimerStatus != 0 {
			core.TimerStatus = 0
			core.TimerStop <- struct{}{}
			core.TaskQ <- struct{}{}
		}
	}()
	time.Sleep(500 * time.Millisecond)

	// 清除日志缓存
	logs.Debug("Started to clean mod <zap>")
	if err := core.Logger.Sync(); err != nil {
		logs.Warn("Cleaning mod <zap> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <zap>")

	// 关闭生产者
	logs.Debug("Started to clean mod <kafka-producer>")
	core.GlobalWg.Wait() // 等待结束
	if err := core.Producer.Close(); err != nil {
		logs.Warn("Cleaning mod <kafka-producer> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <kafka-producer>")

	// 关闭Redis缓存
	logs.Debug("Started to clean mod <redis>")
	if err := core.RedisConn.Close(); err != nil {
		logs.Warn("Cleaning mod <redis> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <redis>")
}
