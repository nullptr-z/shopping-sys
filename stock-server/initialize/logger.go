package initialize

import (
	"fmt"
	"stock-server/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() {
	mode := global.Conf.Mode
	logFilePath := fmt.Sprint(global.Conf.Log.LogFile)
	var config zap.Config
	if mode == "debug" {
		// 开发模式下的配置
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 在开发模式下，使用彩色来区分不同级别的日志
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)                 // 设置日志级别，例如Debug级别
		config.OutputPaths = []string{logFilePath, "stdout"}                // 设置文件作为日志输出目标,和标准错误输出作为日志输出目的地
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder      // 自定义文件名和行号的输出格式
		config.EncoderConfig.StacktraceKey = "stacktrace"                   // 启用堆栈跟踪
		// config.Encoding = "console"                                         // 'json' or 'console'
	} else {
		config = zap.NewProductionConfig()         // 构建一个生产环境下推荐使用的配置
		config.OutputPaths = []string{logFilePath} // 设置文件作为日志输出目标,和标准错误输出作为日志输出目的地
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder // 修改日志时间格式为ISO8601

	// 创建Logger实例
	logger, err := config.Build()
	if err != nil {
		fmt.Println("Failed initialized Init Logger:", err.Error())
		panic("Logger")
	}
	// 创建SugaredLogger
	// sugar := logger.Sugar()

	defer logger.Sync() // 退出时缓冲区的日志都刷到磁盘里

	// 替换全局的 Logger; 其他地方直接 zap.L就能访问这个 logger 了
	zap.ReplaceGlobals(logger)
	fmt.Println("Logger initialized.......")
}
