package storage

type FileService interface {
	UploadFile(data []byte, folder, filename string) (string, error)
	DeleteFile(filename string) (bool, error)
}
