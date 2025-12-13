package fruit

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/dao"
	"RedRockMidAssessment-Synchronizer/core/kafka"
	"RedRockMidAssessment-Synchronizer/core/timer"
	"RedRockMidAssessment-Synchronizer/core/timerlistener"
	viper "RedRockMidAssessment-Synchronizer/core/utils/config"
	zap "RedRockMidAssessment-Synchronizer/core/utils/log"
	"RedRockMidAssessment-Synchronizer/core/utils/snowflake"
	"context"
	"fmt"
	"os"
	"time"

	"gitee.com/liumou_site/logger"
)

const (
	CONFIG_PATH = "./data/config/config.yaml"
	LOG_PATH    = "./data/logs"
)

var (
	cancel context.CancelFunc
	ctx    context.Context
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
	logs.Debug("Started to init mod <kafka-consumer>")
	pc, err := kafka.NewConsumer()
	if err != nil {
		logs.Warn("Init mod <kafka-consumer> error: %v", err.Error())
		os.Exit(1)
	}
	core.PartitionConsumer = pc
	logs.Info("Successfully loaded mod <kafka-consumer>")

	// 初始化通信管道
	core.TimerStop = make(chan struct{}) // 计时器停止指令
	core.TaskQ = make(chan struct{}, 5)  // 全局任务队列

	// 初始化Context
	ctx, cancel = context.WithCancel(context.Background())

	//启动一些监视携程
	go kafka.KLogger(ctx)          // 启动kafka日志携程
	go kafka.MsgReader(ctx)        // 启动信息监听携程
	go timerlistener.Listener(ctx) // 启动计时器信号监听携程

	// 检测并恢复计时器状态
	logs.Debug("Started to recover <timer>")
	if err := RecoverTimerStatus(); err != nil {
		logs.Warn("Recover <timer> error: %v", err.Error())
		os.Exit(1)
	}
	logs.Info("Successfully recovered <timer>")

	fmt.Println()
	logs.Alert("Server Started Successfully, present logs: ")
}

func WorldEndingFruit() {
	// 初始化日志
	logs := logger.NewLogger(1)
	logs.Modular = "WorldEndingFruit"

	// 发送Cancel消息
	cancel()

	// 关闭消费者,等待消息处理完毕
	logs.Debug("Started to clean mod <kafka-consumer>")
	err := core.PartitionConsumer.Close()
	if err != nil {
		logs.Warn("Cleaning mod <kafka-consumer> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <kafka-consumer>")

	// 进行关闭前的一次同步
	if core.TimerStatus != 0 {
		core.TimerStatus = 0
		core.TimerStop <- struct{}{}
		core.TaskQ <- struct{}{}
	}
	//go func() { // 防止死锁
	//	if core.TimerStatus != 0 {
	//		core.TimerStatus = 0
	//		core.TimerStop <- struct{}{}
	//		core.TaskQ <- struct{}{}
	//	}
	//}()
	time.Sleep(500 * time.Millisecond)

	// 等待结束
	core.GlobalWg.Wait()

	// 关闭生产者
	logs.Debug("Started to clean mod <kafka-producer>")
	if err := core.Producer.Close(); err != nil {
		logs.Warn("Cleaning mod <kafka-producer> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <kafka-producer>")

	// 关闭管道
	logs.Debug("Started to close some <pipeline>")
	close(core.TaskQ)
	close(core.TimerStop)
	logs.Info("Successfully closed some <pipeline>")

	time.Sleep(2 * time.Second)

	// 关闭Redis缓存
	logs.Debug("Started to clean mod <redis>")
	if err := core.RedisConn.Close(); err != nil {
		logs.Warn("Cleaning mod <redis> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <redis>")

	// 清除日志缓存
	logs.Debug("Started to clean mod <zap>")
	if err := core.Logger.Sync(); err != nil {
		logs.Warn("Cleaning mod <zap> error: %v", err.Error())
	}
	logs.Info("Successfully cleaned mod <zap>")
}

func RecoverTimerStatus() error {
	// 恢复计时器状态
	var KEY = "courseSelection:status"
	val, err := core.RedisConn.Get(context.Background(), KEY).Result()
	if err != nil {
		return err
	}
	if val == "1" && core.TimerStatus == 0 {
		core.TimerStatus = 1
		go timer.Timer()
	}
	return nil
}
