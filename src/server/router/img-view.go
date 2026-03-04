package router

import (
	"file-taiko-wiki/src/util"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetImgViewRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("fileName")
		if fileName == "" {
			ctx.AbortWithStatus(404)
			return
		}

		fileName = filepath.Clean(fileName)
		if strings.Contains(fileName, "..") {
			ctx.AbortWithStatus(400)
			return
		}

		filePath := filepath.Join(util.MainDir(), "img", fileName)
		stat, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(filePath, "not exists.")
				ctx.AbortWithStatus(404)
				return
			}
			ctx.AbortWithStatus(500)
			return
		}
		if stat.IsDir() {
			log.Println(filePath, "is a directory.")
			ctx.AbortWithStatus(404)
			return
		}

		ctx.File(filePath)
	}
}
