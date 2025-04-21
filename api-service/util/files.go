package util

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

// Description:
// Return an array of bytes from a given multipart header file. Mainly used for
// uploads that handle files for them to be stored later
//
// Parameters:
// header: *multipart.FileHeader
//
// Return:
// []byte
// error
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

// Description:
// Rename a file with a new random UUID and keep the file extension the same.
// Used when replacing files mainly
//
// Parameters:
// filename: string
//
// Return:
// string
// error
func UuidFileRename(filename string) (string, error) {
	extension := filename[strings.LastIndex(filename, ".")+1:]

	if len(extension) < 1 {
		return "", errors.New("there was a problem parsing the extension")
	}

	newName := uuid.New().String()

	return fmt.Sprintf("%s.%s", newName, extension), nil
}

// Description:
// Get a URL from a specified path. Used for getting a valid file url for a
// given publicly uploaded file
//
// Parameters:
// serverAddress: string,
// filePath: string
//
// Return:
// string
func UrlFromFilePath(serverAddress, filePath string) string {
	return fmt.Sprintf("%s/%s", serverAddress, filePath)
}
