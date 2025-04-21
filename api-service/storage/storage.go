package storage

// Description:
// Interface for a file service.
type FileService interface {
	UploadFile(data []byte, folder, fileName string) (string, error)
	DeleteFile(filePath string) error
	ReplaceFile(filePath string, data []byte) (string, error)
}
