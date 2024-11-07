package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitDevLogger() {
	// 使用 zap 的开发配置
	config := zap.NewDevelopmentConfig()

	// 修改 EncoderConfig，使日志等级带颜色
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色日志等级
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 可选：人性化的时间格式

	// 构建 logger
	logger, _ := config.Build()

	// 设置 caller skip，跳过第一层调用位置。自己的log包封装了一层，所以这里设置跳过一层。
	// 【注意】 必须使用此包封装的方法打日志。如果直接用zap.S()，跳过一层后，就不是日志位置了
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)) // 跳过1层

	zap.ReplaceGlobals(logger) // 替换全局的logger，既zap.S() zap.L()
}

func InitPrdLogger() {

}
