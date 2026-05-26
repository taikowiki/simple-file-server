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

		// fileName
		fileName := filepath.Join(filepath.Clean(songNo), filepath.Clean(difficulty)) + ".png"
		if strings.Contains(fileName, "..") {
			ctx.AbortWithStatus(400)
			return
		}

		// filePath
		filePath := filepath.Join(util.FileDir(), "fumen", fileName)
		stat, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("[FumenView] File not found: %s", filePath)
				ctx.AbortWithStatus(404)
				return
			}
			log.Printf("[FumenView] Error checking file: %v (path: %s)", err, filePath)
			ctx.AbortWithStatus(500)
			return
		}
		if stat.IsDir() {
			log.Printf("[FumenView] Path is a directory: %s", filePath)
			ctx.AbortWithStatus(404)
			return
		}

		// send file
		log.Printf("[FumenView] Serving file: %s", fileName)
		ctx.File(filePath)
	}
}
