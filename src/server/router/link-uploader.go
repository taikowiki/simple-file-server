package router

import (
	"file-taiko-wiki/src/server/db"
	"file-taiko-wiki/src/server/serverUtil"
	"file-taiko-wiki/src/util"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 요청 바디 구조체
type LinkUploadRequest struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func GetLinkUploadRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// req link url
		var reqData LinkUploadRequest
		if err := ctx.ShouldBindJSON(&reqData); err != nil {
			log.Printf("[LinkUpload] Request error: Failed to bind JSON from IP: %s", ctx.ClientIP())
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if reqData.URL == "" {
			log.Printf("[LinkUpload] Request error: URL is empty from IP: %s", ctx.ClientIP())
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 인증 로직
		var userUUID string = "server"
		apiKey := os.Getenv("API_KEY")

		// 1. API Key 확인
		if reqData.Key != "" {
			if reqData.Key != apiKey {
				log.Printf("[LinkUpload] Auth failed: Invalid API Key from IP: %s", ctx.ClientIP())
				ctx.JSON(http.StatusForbidden, gin.H{"error": "API Key가 올바르지 않습니다."})
				return
			}
			userUUID = "server"
			log.Printf("[LinkUpload] Auth success: API Key used from IP: %s", ctx.ClientIP())
		} else {
			// 2. 쿠키 확인
			userToken, err := ctx.Cookie("auth-user")
			if err != nil {
				log.Printf("[LinkUpload] Auth failed: No auth-user cookie from IP: %s", ctx.ClientIP())
				ctx.JSON(http.StatusForbidden, gin.H{"error": "인증 정보가 없습니다."})
				return
			}

			authKey := os.Getenv("AUTH_KEY")
			authData, err := serverUtil.Decipher(userToken, authKey)
			if err != nil {
				log.Printf("[LinkUpload] Auth failed: Decipher error: %v from IP: %s", err, ctx.ClientIP())
				ctx.JSON(http.StatusForbidden, gin.H{"error": "인증 정보가 올바르지 않습니다."})
				return
			}

			userDataResults, err := db.GetUserDataByProvider(authData.Provider, authData.ProviderId)
			if err != nil || len(userDataResults) == 0 {
				log.Printf("[LinkUpload] Auth failed: User not found for %s/%s from IP: %s", authData.Provider, authData.ProviderId, ctx.ClientIP())
				ctx.JSON(http.StatusForbidden, gin.H{"error": "유저 정보를 찾을 수 없습니다."})
				return
			}

			userData := userDataResults[0]
			gradeRaw, ok := userData["grade"]
			if !ok {
				log.Printf("[LinkUpload] Auth failed: Grade info missing for user %v from IP: %s", userData["UUID"], ctx.ClientIP())
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "권한 정보를 확인할 수 없습니다."})
				return
			}

			var grade float64
			switch v := gradeRaw.(type) {
			case float64:
				grade = v
			case int:
				grade = float64(v)
			default:
				log.Printf("[LinkUpload] Auth failed: Invalid grade type for user %v from IP: %s", userData["UUID"], ctx.ClientIP())
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "권한 형식이 올바르지 않습니다."})
				return
			}

			if grade < 9 {
				log.Printf("[LinkUpload] Auth failed: Insufficient grade (%v) for user %v from IP: %s", grade, userData["UUID"], ctx.ClientIP())
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "권한이 부족합니다."})
				return
			}

			userUUID, _ = userData["UUID"].(string)
			log.Printf("[LinkUpload] Auth success: User %v authenticated via cookie from IP: %s", userUUID, ctx.ClientIP())
		}

		log.Printf("[LinkUpload] Starting download: %s for user %v", reqData.URL, userUUID)

		// fetch
		resp, err := http.Get(reqData.URL)
		if err != nil || resp.StatusCode != http.StatusOK {
			status := "unknown"
			if resp != nil {
				status = resp.Status
			}
			log.Printf("[LinkUpload] Fetch error: %v (Status: %s) for URL: %s", err, status, reqData.URL)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// ext
		contentType := resp.Header.Get("Content-Type")
		exts, _ := mime.ExtensionsByType(contentType)
		var ext string
		if len(exts) > 0 {
			ext = exts[0]
		}
		log.Printf("[LinkUpload] Detected content type: %s, extension: %s", contentType, ext)

		// file name
		var fileName string
		uploadDir := filepath.Join(util.FileDir(), "img")
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}
		for {
			fileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
			filePath := filepath.Join(uploadDir, fileName)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				break
			}
		}

		// save
		savePath := filepath.Join(uploadDir, fileName)
		file, err := os.Create(savePath)
		if err != nil {
			log.Printf("[LinkUpload] File creation error: %v (path: %s)", err, savePath)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		written, err := io.Copy(file, resp.Body)
		if err != nil {
			log.Printf("[LinkUpload] File write error: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// DB 로그 기록
		db.NewFileLog(userUUID, reqData.URL, fileName)

		log.Printf("[LinkUpload] Successfully saved: %s -> %s (%d bytes) for user %v", reqData.URL, fileName, written, userUUID)

		ctx.JSON(http.StatusOK, gin.H{
			"fileName": fileName,
		})
	}
}
