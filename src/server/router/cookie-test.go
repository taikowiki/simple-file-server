package router

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetCookieTestRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("auth-user")
		if err != nil {
			log.Printf("[CookieTest] No auth-user cookie found from IP: %s", ctx.ClientIP())
			ctx.AbortWithStatus(403)
			return
		}
		log.Printf("[CookieTest] Successfully retrieved auth-user cookie from IP: %s", ctx.ClientIP())
		ctx.String(200, cookie)
	}
}
