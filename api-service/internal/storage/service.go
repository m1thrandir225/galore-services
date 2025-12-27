// Package storage provides a way to manage and store files.
package storage

// Service is an interface for uploading and managing files
type Service interface {
	UploadFile(data []byte, folder, fileName string) (string, error)
	DeleteFile(filePath string) error
	ReplaceFile(filePath string, data []byte) (string, error)
}
