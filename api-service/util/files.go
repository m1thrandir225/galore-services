package util

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"strings"
)

func BytesFromFile(header *multipart.FileHeader) ([]byte, error) {
	opened, err := header.Open()

	if err != nil {
		return nil, err
	}
	defer opened.Close()

	fileData, err := io.ReadAll(opened)

	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func UuidFileRename(filename string) (string, error) {
	extension := filename[strings.LastIndex(filename, ".")+1:]

	if len(extension) < 1 {
		return "", errors.New("there was a problem parsing the extension")
	}

	newName := uuid.New().String()

	return fmt.Sprintf("%s.%s", newName, extension), nil
}

func UrlFromFilePath(serverAddress, filePath string) string {
	return fmt.Sprintf("%s/%s", serverAddress, filePath)
}
