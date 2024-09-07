package storage

type FileService interface {
	UploadFile(data []byte, filename string) (string, error)
	DeleteFile(filename string) (bool, error)
}
