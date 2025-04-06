package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitDevLogger 初始化开发环境日志配置
// 入参：projectName 项目名称，如 user_srv
// 特点：
// 1. 输出带颜色的日志级别
// 2. 使用人性化的时间格式
// 3. 输出调用位置信息
// 4. 开发环境默认记录 Debug 及以上级别的日志
func InitDevLogger(projectName string) {
	// 使用 zap 的开发配置，默认记录 Debug 及以上级别
	config := zap.NewDevelopmentConfig()

	// 添加固定前缀字段
	config.InitialFields = map[string]interface{}{
		"project": projectName, // 应用名称
	}

	// 修改 EncoderConfig，使日志更易读
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                      // 使用大写带颜色的日志级别
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // 东八区时间格式 毫秒级
		enc.AppendString(t.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02T15:04:05.000+08:00"))
	}

	// 构建 logger
	logger, _ := config.Build()

	// 设置 caller skip，跳过第一层调用位置
	// 由于我们封装了日志包，需要跳过一层才能显示真实的调用位置
	// 【重要】必须使用此包封装的方法打日志，如果直接使用 zap.S() 或 zap.L()，由于跳过了一层，将无法显示正确的调用位置
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	// 替换全局的 logger，这样可以直接使用 zap.S() 或 zap.L() 打日志
	zap.ReplaceGlobals(logger)
}

// InitPrdLogger 初始化生产环境日志配置
// 入参：projectName 项目名称，如 user_srv
// 特点：
// 1. 使用 JSON 格式输出，便于日志收集和解析
// 2. 日志文件自动轮转，避免单个文件过大
// 3. 错误日志单独收集
// 4. Info 及以上级别会记录到 app.log
// 5. Error 及以上级别会记录到 error.log
// 6. 同时在控制台输出 Info 及以上级别的日志
func InitPrdLogger(projectName string, config ...PrdLoggerConfig) {
	// 确保日志路径存在
	logDir := fmt.Sprintf("/usr/local/yeying/projects/%s/logs", projectName) // 临时路径老有问题，没深究
	allProjectLogDir := "/usr/local/yeying/unilogs"
	if len(config) > 0 {
		// 如果传入了配置，则使用传入的配置
		if config[0].logDir != "" {
			logDir = config[0].logDir
		}
		if config[0].allProjectLogDir != "" {
			allProjectLogDir = config[0].allProjectLogDir
		}
	}
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	// 配置 JSON 编码器，定义日志格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",     // 时间戳字段名
		LevelKey:      "level",         // 日志级别字段名
		NameKey:       "logger",        // logger名字字段名
		CallerKey:     "caller",        // 调用者字段名
		FunctionKey:   zapcore.OmitKey, // 调用函数名字段名，这里选择省略
		MessageKey:    "msg",           // 消息字段名
		StacktraceKey: "stacktrace",    // 堆栈跟踪字段名
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // 东八区时间格式 毫秒级
			enc.AppendString(t.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02T15:04:05.000+08:00"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder, // 持续时间使用秒作为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
	}

	// 配置普通日志文件的轮转规则 -- 记录在当前项目下
	appLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, fmt.Sprintf("%s.log", projectName)), // 日志文件路径
		MaxSize:    100,                                                       // 单个文件最大尺寸，单位 MB
		MaxBackups: 60,                                                        // 保留旧文件的最大个数
		MaxAge:     30,                                                        // 保留旧文件的最大天数
		Compress:   true,                                                      // 是否压缩/归档旧文件
	}

	// 配置错误日志文件的轮转规则 -- 记录在当前项目下
	errorLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, fmt.Sprintf("err_%s.log", projectName)), // 错误日志文件路径
		MaxSize:    100,                                                           // 单个文件最大尺寸，单位 MB
		MaxBackups: 120,                                                           // 保留旧文件的最大个数
		MaxAge:     60,                                                            // 保留旧文件的最大天数
		Compress:   true,                                                          // 是否压缩/归档旧文件
	}

	// 配置所有日志文件的轮转规则 -- 所有项目的所有日志都写入到 all.log，方便临时分析问题
	allLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(allProjectLogDir, "all.log"), // 所有日志文件路径
		MaxSize:    300,                                        // 单个文件最大尺寸，单位 MB
		MaxBackups: 10,                                         // 保留旧文件的最大个数
		MaxAge:     3,                                          // 保留旧文件的最大天数
		Compress:   true,                                       // 是否压缩/归档旧文件
	}

	// 配置日志级别过滤器
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel // Error 及以上级别
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel // Info 到 Warn 级别
	})

	// 配置多个输出核心
	// 使用 NewTee 将日志输出到多个位置
	core := zapcore.NewTee(
		// 1. Info 到 Warn 级别的日志写入 app.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(appLogWriter),
			lowPriority,
		),
		// 2. Error 及以上级别的日志写入 error.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(errorLogWriter),
			highPriority,
		),
		// 3. Info 及以上级别的日志同时输出到控制台
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			lowPriority,
		),
		// 4. 所有级别的日志写入 all.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(allLogWriter),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return true }), // 记录所有级别的日志
		),
	)

	// 构建最终的 logger
	logger := zap.New(core,
		zap.AddCaller(),                   // 添加调用者信息
		zap.AddCallerSkip(1),              // 跳过一层调用栈，显示实际的调用位置
		zap.AddStacktrace(zap.ErrorLevel), // Error 及以上级别显示堆栈信息
		// 添加固定前缀字段
		zap.Fields(
			zap.String("project", projectName), // 应用名称
		),
	)

	// 替换全局的 logger
	zap.ReplaceGlobals(logger)
}

type PrdLoggerConfig struct {
	logDir           string
	allProjectLogDir string
}
