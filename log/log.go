package log

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type LogBasicInfo struct {
	RequestId string
	Uid       string
}

func parseInfoFromContext(ctx context.Context) *LogBasicInfo {
	lbi := &LogBasicInfo{}
	switch ctx.(type) {
	case *gin.Context:
		gctx := ctx.(*gin.Context)
		if v, exist := gctx.Get("request_id"); exist {
			lbi.RequestId = v.(string)
		}
		if v, exist := gctx.Get("uid"); exist {
			lbi.Uid = v.(string)
		}
	case context.Context:
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			lbi.RequestId = md["request_id"][0]
			lbi.Uid = md["uid"][0]
		}
	}
	return lbi
}

// 示例 log.Debug(ctx, "调试一下", "order_id", order_id)
func Debug(ctx context.Context, msg string, kv ...interface{}) {
	lbi := parseInfoFromContext(ctx)
	args := append([]interface{}{
		"request_id", lbi.RequestId,
		"uid", lbi.Uid,
	}, kv...)
	zap.S().Debugw(msg, args...)
}

// 示例 log.Info(ctx, "操作成功啦", "order_id", order_id)
func Info(ctx context.Context, msg string, kv ...interface{}) {
	lbi := parseInfoFromContext(ctx)
	args := append([]interface{}{
		"request_id", lbi.RequestId,
		"uid", lbi.Uid,
	}, kv...)
	zap.S().Infow(msg, args...)
}

// 示例 log.Warn(ctx, "警告", "order_id", order_id)
func Warn(ctx context.Context, msg string, kv ...interface{}) {
	lbi := parseInfoFromContext(ctx)
	args := append([]interface{}{
		"request_id", lbi.RequestId,
		"uid", lbi.Uid,
	}, kv...)
	zap.S().Warnw(msg, args...)
}

// 示例 log.Error(ctx, "出错啦", "order_id", order_id)
func Error(ctx context.Context, msg string, kv ...interface{}) {
	lbi := parseInfoFromContext(ctx)
	args := append([]interface{}{
		"request_id", lbi.RequestId,
		"uid", lbi.Uid,
	}, kv...)
	zap.S().Errorw(msg, args...)
}

// 示例 log.Fatal(ctx, "严重错误", "order_id", order_id)
func Fatal(ctx context.Context, msg string, kv ...interface{}) {
	lbi := parseInfoFromContext(ctx)
	args := append([]interface{}{
		"request_id", lbi.RequestId,
		"uid", lbi.Uid,
	}, kv...)
	zap.S().Fatalw(msg, args...)
}
