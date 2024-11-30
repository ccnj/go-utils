package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenRequestId2Ctx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 存requestId
		requestId := uuid.New().String()
		ctx.Set("request_id", requestId)
	}
}
