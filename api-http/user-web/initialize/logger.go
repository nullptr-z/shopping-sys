package initialize

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() error {
	mode := viper.Get("mode")
	logFilePath := fmt.Sprint(viper.Get("log.log_file"))
	var config zap.Config
	if mode == gin.DebugMode {
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
		return err
	}
	// 创建SugaredLogger
	// sugar := logger.Sugar()

	defer logger.Sync() // 退出时缓冲区的日志都刷到磁盘里

	// 替换全局的 Logger; 其他地方直接 zap.L就能访问这个 logger 了
	zap.ReplaceGlobals(logger)
	fmt.Println("Logger initialized.......")
	return nil
}

// 自定义日志输出格式
func LoggerFormateOutput(g *gin.Context) {
	// 请求前
	startTime := time.Now()

	// 复制请求体，以便日志记录后仍可读取
	var requestBody bytes.Buffer
	if g.Request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(g.Request.Body)
		requestBody.Write(bodyBytes)
		// 重新设置请求体，以供后续使用
		g.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	queryParams := g.Request.URL.Query()
	// fmt.Println("request arguments:", requestBody.String())
	// fmt.Println("query arguments:", queryParams)

	// 处理请求
	g.Next()

	// 请求后
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	statusCode := g.Writer.Status()
	clientIP := g.ClientIP()

	// body := prettyPrintJSON(requestBody)
	fmt.Println("body:", requestBody.String())
	fmt.Println("queryParams:", queryParams)

	// 使用方括号[]格式化日志内容
	zap.L().Info("requestDetails",
		zap.String("method", g.Request.Method),
		zap.String("uri", g.Request.RequestURI),
		zap.Int("status", statusCode),
		zap.String("latency", fmt.Sprintf("[%s]", latencyTime)),
		zap.String("clientIP", fmt.Sprintf("[%s]", clientIP)),
		// zap.String("RequestArguments", fmt.Sprint("body:", requestBody.String())),
		// zap.String("QueryParams", fmt.Sprintf("[%s]", queryParams)),
		// zap.String("formData", fmt.Sprintf("[%s]", formData)),
	)

}
