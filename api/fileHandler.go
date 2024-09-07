package api

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func extractFileBytes(ctx *gin.Context) ([]byte, error) {
	file, _, err := ctx.Request.FormFile("file")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bytes.NewBuffer(nil)

	if _, err = io.Copy(buffer, file); err != nil {
		return nil, err
	}

	fileBytes := buffer.Bytes()

	return fileBytes, nil

}
