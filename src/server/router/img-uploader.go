package router

import (
	"file-taiko-wiki/src/server/db"
	"file-taiko-wiki/src/server/serverUtil"
	"file-taiko-wiki/src/util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetImgUploadRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 인증 확인
		userToken, err := ctx.Cookie("auth-user")
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "인증 토큰이 없습니다."})
			return
		}

		authKey := os.Getenv("AUTH_KEY")
		authData, err := serverUtil.Decipher(userToken, authKey)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "인증 정보가 올바르지 않습니다."})
			return
		}

		userDataResults, err := db.GetUserDataByProvider(authData.Provider, authData.ProviderId)
		if err != nil || len(userDataResults) == 0 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "유저 정보를 찾을 수 없습니다."})
			return
		}

		userData := userDataResults[0]
		gradeRaw, ok := userData["grade"]
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "권한 정보를 확인할 수 없습니다."})
			return
		}

		// grade는 float64로 올 가능성이 큼 (JSON unmarshal 특성상)
		grade, ok := gradeRaw.(float64)
		if !ok {
			// int일 경우도 대비
			gradeInt, ok := gradeRaw.(int)
			if ok {
				grade = float64(gradeInt)
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "권한 형식이 올바르지 않습니다."})
				return
			}
		}

		if grade < 9 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "권한이 부족합니다."})
			return
		}

		// 파일 가져오기
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "파일을 찾을 수 없습니다."})
			return
		}

		// file name 생성
		var fileName string
		uploadDir := filepath.Join(util.FileDir(), "img")
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}
		for {
			ext := filepath.Ext(file.Filename)
			fileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)

			filePath := filepath.Join(uploadDir, fileName)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				break
			}
		}

		// save file
		savePath := filepath.Join(uploadDir, fileName)
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "파일 저장 중 오류가 발생했습니다."})
			return
		}

		// DB 로그 기록
		userUUID, _ := userData["UUID"].(string)
		originalName := file.Filename
		if originalName == "" {
			originalName = "undefined"
		}
		db.NewFileLog(userUUID, originalName, fileName)

		// 응답
		ctx.JSON(http.StatusOK, gin.H{
			"fileName": fileName,
		})
	}
}
