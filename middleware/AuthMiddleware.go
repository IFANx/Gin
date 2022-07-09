package middleware

import (
	"Gin/common"
	"Gin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	获取authentication header
		tokenString := ctx.GetHeader("Authorization")
		//	验证token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 400, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 400, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//	获取通过验证的token种的Userid
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//	用户不存在,比如用户已经被是删除
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 400, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//	若用户存在，将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
