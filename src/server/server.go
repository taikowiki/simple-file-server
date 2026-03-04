package server

import (
	"file-taiko-wiki/src/server/router"

	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	server := gin.Default()
	server.GET("/img/:fileName", router.GetImgViewRouter())
	server.GET("/fumen/:fileName", router.GetFumenViewRouter())
	server.GET("/cookie-test", router.GetCookieTestRouter())

	return server
}
