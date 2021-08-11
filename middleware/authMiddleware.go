package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"zhiyudong.cn/gin-test/common"
	"zhiyudong.cn/gin-test/model"
	"zhiyudong.cn/gin-test/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			response.Response(ctx, http.StatusUnauthorized, 40001, nil, "权限不足")
			// ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort() // 抛弃请求
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { // 解析失败 || token 失效
			response.Response(ctx, http.StatusUnauthorized, 40001, nil, "权限不足")
			ctx.Abort() // 抛弃请求
			return
		}

		userId := claims.UserId
		DB := common.GetDB()
		var user model.User

		DB.First(&user, userId)

		// 用户
		if user.ID == 0 {
			response.Response(ctx, http.StatusUnauthorized, 40001, nil, "权限不足")
			ctx.Abort() // 抛弃请求
			return
		}

		// 用户存在  将 user 的信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
