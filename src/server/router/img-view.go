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

		filePath := filepath.Join(util.FileDir(), "img", fileName)
		stat, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("[ImgView] File not found: %s", filePath)
				ctx.AbortWithStatus(404)
				return
			}
			log.Printf("[ImgView] Error checking file: %v (path: %s)", err, filePath)
			ctx.AbortWithStatus(500)
			return
		}
		if stat.IsDir() {
			log.Printf("[ImgView] Path is a directory: %s", filePath)
			ctx.AbortWithStatus(404)
			return
		}

		log.Printf("[ImgView] Serving file: %s", fileName)
		ctx.File(filePath)
	}
}
