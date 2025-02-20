package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func render_v2(c *gin.Context, data gin.H, HTMLtemplateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, HTMLtemplateName, data)
	}
}

func SaveFile(formName string, c *gin.Context) (fileNameForDB string, err error) {
	fileHeader, _ := c.FormFile(formName)
	if fileHeader == nil {
		fmt.Fprintln(c.Writer, "не был прикреплён файл")
		return "", errors.New("file doesn't send")
	}

	ext := filepath.Ext(fileHeader.Filename)
	filenameAsUNIXsecond := fmt.Sprint(time.Now().UnixNano())
	fileNameForDB = filenameAsUNIXsecond + ext
	err = c.SaveUploadedFile(fileHeader, "./v2/loaded_img/"+fileNameForDB)
	if err != nil {
		return "", err
	}

	return fileNameForDB, nil
}
