package log

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type LogBasicInfo struct {
	RequestId string
	Uid       int64
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
			lbi.Uid = v.(int64)
		}
	case context.Context:
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			lbi.RequestId = md["request_id"][0]
			uidStr := md["uid"][0]
			if uidInt, err := strconv.Atoi(uidStr); err == nil {
				lbi.Uid = int64(uidInt)
			}
		}
	}
	return lbi
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
