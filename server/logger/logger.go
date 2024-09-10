package logger

import (
	"fmt"
	"os"
	"test/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func Init() error {
	// 从配置中获取日志级别
	level, err := zapcore.ParseLevel(config.AppConfig.Logger.Level)
	if err != nil {
		return fmt.Errorf("解析日志级别失败: %v", err)
	}

	// 配置 lumberjack 日志切割
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.AppConfig.Logger.Filename,
		MaxSize:    config.AppConfig.Logger.MaxSize, // 单位：MB
		MaxBackups: config.AppConfig.Logger.MaxBackups,
		MaxAge:     config.AppConfig.Logger.MaxAge, // 单位：天
		Compress:   config.AppConfig.Logger.Compress,
	})

	// 创建 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建文件输出的 core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		level,
	)

	// 创建控制台输出的 core
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, level)

	// 使用 NewTee 将两个 core 组合
	core := zapcore.NewTee(fileCore, consoleCore)

	// 创建 logger
	Log = zap.New(core, zap.AddCaller())
	Sugar = Log.Sugar()

	return nil
}
