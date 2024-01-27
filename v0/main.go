package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"objectify/log"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const path = "object-storage/v0"

func init() {
	logFile := fmt.Sprintf("storage-v0-%s.log", time.Now().Format("20060102150405"))
	infoFile := fmt.Sprintf("storage-v0-info-%s.log", time.Now().Format("20060102150405"))
	errorFile := fmt.Sprintf("storage-v0-error-%s.log", time.Now().Format("20060102150405"))
	log.SetupZap([]string{infoFile, logFile}, []string{errorFile, logFile}, zap.InfoLevel)
	gin.SetMode(gin.ReleaseMode)
	// create object storage path if not exists
	err := os.MkdirAll(path, 0777)
	if err != nil {
		zap.L().Fatal("create obj storage path:", zap.Error(err))
	}
}

func main() {
	r := gin.Default()
	r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))
	r.PUT("/v0/object/:filename", put)
	r.GET("/v0/object/:filename", get)
	err := r.Run("localhost:8080")
	if err != nil {
		zap.L().Fatal("ListenAndServe", zap.Error(err))
	}
}

func put(c *gin.Context) {
	escapedFilename := c.Param("filename")
	// get unescaped filename from url
	filename, err := url.QueryUnescape(escapedFilename)
	if err != nil {
		c.String(400, "fail to get unescaped filename: %s", escapedFilename)
		return
	}
	// create file if it does not exist
	_, err = os.Open(filepath.Join(path, filename))
	if err != nil {
		file, err := os.OpenFile(filepath.Join(path, filename), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			zap.L().Error("fail to create file", zap.Error(err), zap.String("filename", filename))
			c.String(400, "fail to create file: %s", filename)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				zap.L().Error("fail to close file", zap.Error(err), zap.String("filename", filename))
			}
		}(file)
		_, err = io.Copy(file, c.Request.Body)
		if err != nil {
			zap.L().Error("fail to write file", zap.Error(err), zap.String("filename", filename))
			c.String(400, "fail to write file: %s", filename)
			return
		}
		zap.L().Info("put file", zap.String("filename", filename), zap.String("ip", c.ClientIP()))
		c.String(200, "ok")
		return
	} else {
		zap.L().Info("file already exists", zap.String("filename", filename), zap.String("ip", c.ClientIP()))
		c.String(400, "file already exists: %s", filename)
		return
	}
}

func get(c *gin.Context) {
	escapedFilename := c.Param("filename")
	// get unescaped filename from url
	filename, err := url.QueryUnescape(escapedFilename)
	if err != nil {
		c.String(400, "fail to get unescaped filename: %s", escapedFilename)
		return
	}
	file, err := os.Open(filepath.Join(path, filename))
	if err != nil {
		zap.L().Error("fail to open file", zap.Error(err), zap.String("filename", filename))
		c.String(400, "fail to open file: %s", filename)
		return
	}
	defer file.Close()
	c.Header("Content-Disposition", "attachment;filename="+filename)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		zap.L().Error("fail to response file", zap.Error(err), zap.String("filename", filename))
		c.String(400, "fail to response file: %s", filename)
		return
	}
}
