package log

import "go.uber.org/zap"

// 用于打印不含上下文信息的，纯日志
// 因为统一跳过1层调用堆栈信息，所以也必须都包装一层
type Pure struct{}

// 示例 log.Pure{}.Debug("打日志了", "order_id", order_id)
func (p Pure) Debug(msg string, kv ...interface{}) {
	zap.S().Debugw(msg, kv...)
}

// 示例 log.Pure{}.Info("打日志了", "order_id", order_id)
func (p Pure) Info(msg string, kv ...interface{}) {
	zap.S().Infow(msg, kv...)
}

// 示例 log.Pure{}.Warn("打日志了", "order_id", order_id)
func (p Pure) Warn(msg string, kv ...interface{}) {
	zap.S().Warnw(msg, kv...)
}

// 示例 log.Pure{}.Error("打日志了", "order_id", order_id)
func (p Pure) Error(msg string, kv ...interface{}) {
	zap.S().Errorw(msg, kv...)
}

// 示例 log.Pure{}.Fatal("打日志了", "order_id", order_id)
func (p Pure) Fatal(msg string, kv ...interface{}) {
	zap.S().Fatalw(msg, kv...)
}
