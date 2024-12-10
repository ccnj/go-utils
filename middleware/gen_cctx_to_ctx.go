// cctx是context.Context，是grpc服务层的入参之一，此处用于传递requestId，uid

package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func GenCctx2Ctx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestId string
		var uid int64
		if requestIdAny, exist := ctx.Get("request_id"); exist {
			if requestIdStr, ok := requestIdAny.(string); ok {
				requestId = requestIdStr
			}
		}
		if uidAny, exist := ctx.Get("uid"); exist {
			if uidInt64, ok := uidAny.(int64); ok {
				uid = uidInt64
			}
		}

		// requestId，unsafeUID存入cctx，用于传给grpc服务，告知请求信息
		cctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
			"request_id", requestId, // metadata中，key会被转为小写，所以统一用蛇形
			"uid", fmt.Sprintf("%d", uid),
		))
		ctx.Set("cctx", cctx)
	}
}
