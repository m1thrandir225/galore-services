package shared

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

// BytesFromFile returns a byte array from a given FileHeader
func BytesFromFile(header *multipart.FileHeader) ([]byte, error) {
	opened, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := opened.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	fileData, err := io.ReadAll(opened)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

// UUIDFileRename renames a given filename to a random uuid whilst preserving the extension
func UUIDFileRename(filename string) (string, error) {
	extension := filename[strings.LastIndex(filename, ".")+1:]

	if len(extension) < 1 {
		return "", errors.New("there was a problem parsing the extension")
	}

	newName := uuid.New().String()

	return fmt.Sprintf("%s.%s", newName, extension), nil
}

// URLFromFilePath returns the corrent URL given the server address and filePath
func URLFromFilePath(serverAddress, filePath string) string {
	return fmt.Sprintf("%s/%s", serverAddress, filePath)
}
