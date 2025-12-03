package zap

import (
	"RedRockMidAssessment-Consumer/core"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logLevel = zap.DebugLevel           //debug级别
var timeFmt = "2006-01-02 15:04:05.000" //格式化时间

func InitZap(path string) {
	// 轮转文件
	logFile := &lumberjack.Logger{
		Filename:   path,  // 日志路径
		MaxSize:    10,    // 单个文件最大10M
		MaxBackups: 3,     // 最多3个备份文件夹
		MaxAge:     7,     // 日志最多保留7天
		Compress:   false, // 是否压缩日志
	}

	//JSON编码器
	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",                                 //时间戳的字段名
		LevelKey:       "level",                              //日志级别字段名
		NameKey:        "logger",                             //产生日志模块的字段名（Logger.Named("xxx") 时才有）eg."logger":"http.server"
		CallerKey:      "caller",                             //代码位置字段名
		MessageKey:     "msg",                                //日志正文字段名
		StacktraceKey:  "stacktrace",                         //堆栈字段名，触发Panic才有
		LineEnding:     zapcore.DefaultLineEnding,            //每条日志结尾换行符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,        //级别全小写编码
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFmt), //时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,       //time.Duration采用小数编码 eg."1.234"（秒，浮点）
		EncodeCaller:   zapcore.ShortCallerEncoder,           //代码文件路径
	})

	//JSON日志核心
	jsonCore := zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(logFile),
		logLevel,
	)

	//控制台日志编码器
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "T", //同上
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFmt),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	//控制台日志核心
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel)

	//合并两个日志核心
	core.Logger = zap.New(zapcore.NewTee(jsonCore, consoleCore), zap.AddCaller())
}
