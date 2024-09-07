package storage

import (
	"os"
	"path"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
	}
}

func (storage *LocalStorage) UploadFile(data []byte, filename string) (string, error) {
	permissions := os.FileMode(0644)
	filePath := path.Join(storage.BasePath, filename)
	err := os.WriteFile(filePath, data, permissions)

	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (storage *LocalStorage) DeleteFile(filename string) (bool, error) {
	filePath := path.Join(storage.BasePath, filename)

	_, err := os.Stat(filePath)

	if err != nil {
		return false, os.ErrNotExist
	}

	err = os.Remove(filePath)

	if err != nil {
		return false, err
	}

	return true, nil
}
