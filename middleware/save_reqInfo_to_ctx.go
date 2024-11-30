package middleware

// import (
// 	"context"
// 	"encoding/base64"
// 	"encoding/json"
// 	"strconv"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"google.golang.org/grpc/metadata"
// )

// func SaveReqInfo2Ctx() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// 1. 存requestId
// 		requestId := uuid.New().String()
// 		ctx.Set("request_id", requestId)

// 		// 2. 存unsafeUID （token中解析的, 为确保速度，未验证签名, 故unsafe）
// 		unsafeUID := getUnsafeUID(ctx)
// 		ctx.Set("unsafe_uid", unsafeUID)

// 		// 3. requestId，unsafeUID存入cctx，用于传给grpc服务，告知请求信息
// 		cctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
// 			"request_id", requestId, // metadata中，key会被转为小写，所以统一用蛇形
// 			"unsafe_uid", unsafeUID,
// 		))
// 		ctx.Set("cctx", cctx)
// 	}
// }

// // 工具函数
// func getUnsafeUID(ctx *gin.Context) string {
// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		return "no-token"
// 	}
// 	authHeaderSli := strings.Split(authHeader, " ")
// 	if len(authHeaderSli) != 2 {
// 		return "invalid-auth-format"
// 	}
// 	tokenSli := strings.Split(authHeaderSli[1], ".")
// 	if len(tokenSli) != 3 {
// 		return "invalid-token-format"
// 	}
// 	uidB64 := tokenSli[1]

// 	uidBytes, err := base64.StdEncoding.DecodeString(uidB64)
// 	if err != nil {
// 		return "invalid-uidB64"
// 	}
// 	claims := myCustomClaims{}
// 	err = json.Unmarshal(uidBytes, &claims)
// 	if err != nil {
// 		return "invalid-claims"
// 	}
// 	return strconv.Itoa(int(claims.UID))
// }
