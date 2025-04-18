package storage

import (
	"os"
	"path"
	"path/filepath"

	"github.com/m1thrandir225/galore-services/util"
)

// Description:
// LocalStorage implementation of the FileService interface
//
// Parameters:
// BasePath: string
type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
	}
}

func (storage *LocalStorage) UploadFile(data []byte, folder, fileName string) (string, error) {
	permissions := os.FileMode(0644)

	folderPath := path.Join(storage.BasePath, folder)

	err := os.MkdirAll(folderPath, os.FileMode(0700))
	if err != nil {
		return "", err
	}

	uniqueFilename, err := util.UuidFileRename(fileName)
	if err != nil {
		return "", err
	}

	filePath := path.Join(folderPath, uniqueFilename)

	err = os.WriteFile(filePath, data, permissions)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (storage *LocalStorage) DeleteFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return os.ErrNotExist
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (storage *LocalStorage) ReplaceFile(filePath string, data []byte) (string, error) {
	dirs := filepath.Dir(filePath)

	userIdFolder := filepath.Base(dirs)

	fileName := filepath.Base(filePath)

	err := storage.DeleteFile(filePath)
	if err != nil {
		return "", err
	}

	newFilePath, err := storage.UploadFile(data, userIdFolder, fileName)
	if err != nil {
		return "", err
	}
	return newFilePath, nil
}
