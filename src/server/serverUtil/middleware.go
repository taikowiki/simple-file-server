package serverUtil

import "github.com/gin-gonic/gin"

func GetCORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "https://taiko.wiki")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Next()
	}
}
