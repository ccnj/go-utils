package middleware

import (
	"errors"
	"net/http"

	"strings"
	"time"

	"github.com/ccnj/go-utils/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type myCustomClaims struct {
	UID  string `json:"uid"`
	Role int32  `json:"role"`
	jwt.RegisteredClaims
}

// claims.UID 用户id
// claims.RegisteredClaims.ExpiresAt 过期时间
func parseToken(tokenStr string, jwtSigningKey string) (claims *myCustomClaims, e error) {
	token, err := jwt.ParseWithClaims(tokenStr, &myCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSigningKey), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		log.Pure{}.Info("parse token err", "err", err)
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			// 签名错误，说明有人在攻击了，记录一下
			log.Pure{}.Error("签名错误，疑似受到攻击",
				"tokenStr", tokenStr,
				"err", err,
			)
		}
		return nil, err
	} else if claims, ok := token.Claims.(*myCustomClaims); ok {
		// fmt.Println(claims.UID, claims.RegisteredClaims.ExpiresAt)
		return claims, nil
	} else {
		return nil, errors.New("未知的claims类型, 无法继续")
	}
}

func canSkip(ctx *gin.Context, skipPathsPrefix []string) bool {
	for _, path := range skipPathsPrefix {
		if strings.HasPrefix(ctx.Request.URL.Path, path) {
			return true
		}
	}
	return false
}

func ValidateToken(jwtSigningKey string, skipPathsPrefix []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 跳过不需要验证的路由
		if canSkip(ctx, skipPathsPrefix) {
			// ctx.Next() // 会自动执行，可以不显示写出来
			ctx.Set("uid", "")
			return
		}

		// 获取请求头中的Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"errCode": 401,
				"errMsg":  "尚未登录，请先登录～",
			})
			ctx.Abort() // 必须显式地中止。因为gin中，即使没有ctx.Next()，也会在中间件结束时自动执行下一个
			return
		}

		// 从Authorization头中提取Token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"errCode": 401,
				"errMsg":  "授权头格式错误",
			})
			ctx.Abort()
			return
		}

		// 验证token有效性
		claims, err := parseToken(token, jwtSigningKey)
		if err != nil {
			// 验证失败
			var errMsg string
			if errors.Is(err, jwt.ErrTokenExpired) {
				errMsg = "身份认证已过期，请重新登录"
			} else {
				errMsg = "身份认证失败，请先登录"
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"errCode": 401,
				"errMsg":  errMsg,
			})
			ctx.Abort()
			return
		}

		// 保存uid至ctx中
		ctx.Set("uid", claims.UID)
		ctx.Set("role", int64(claims.Role)) // ctx无getInt32方法，所以存int64，取的时候也必须ctx.GetInt64("role") GetInt取不到
		// 执行后续中间件
		// ctx.Next()

	}
}
