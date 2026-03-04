package router

import "github.com/gin-gonic/gin"

func GetCookieTestRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("auth-user")
		if err != nil {
			ctx.AbortWithStatus(403)
			return
		}
		ctx.String(200, cookie)
	}
}
