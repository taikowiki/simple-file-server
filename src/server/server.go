package server

import (
	"file-taiko-wiki/src/server/router"
	"file-taiko-wiki/src/server/serverUtil"

	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	server := gin.Default()
	server.Use(serverUtil.GetCORSMiddleware())
	server.GET("/img/:fileName", router.GetImgViewRouter())
	server.GET("/fumen/:songNo/:difficulty", router.GetFumenViewRouter())
	server.GET("/cookie-test", router.GetCookieTestRouter())
	server.POST("/upload/img", router.GetImgUploadRouter())
	server.POST("/upload/link", router.GetLinkUploadRouter())

	return server
}
