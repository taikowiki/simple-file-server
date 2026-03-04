package router

import (
	"file-taiko-wiki/src/util"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetFumenViewRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		songNo := ctx.Param("songNo")
		if songNo == "" {
			ctx.AbortWithStatus(404)
			return
		}
		difficulty := ctx.Param("difficulty")
		if difficulty == "" {
			ctx.AbortWithStatus(404)
			return
		}

		fileName := filepath.Join(filepath.Clean(songNo), filepath.Clean(difficulty)) + ".png"
		if strings.Contains(fileName, "..") {
			ctx.AbortWithStatus(400)
			return
		}

		filePath := filepath.Join(util.MainDir(), "fumen", fileName)
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
