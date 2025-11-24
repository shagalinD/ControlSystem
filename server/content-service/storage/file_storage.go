package storage

import (
	"io"
	"os"
	"path/filepath"
)

type FileStorage struct {
    UploadPath string
}

func NewFileStorage(uploadPath string) *FileStorage {
    // Создаем директорию если не существует
    os.MkdirAll(uploadPath, 0755)
    return &FileStorage{
        UploadPath: uploadPath,
    }
}

func (fs *FileStorage) SaveFile(filename string, file io.Reader) (string, error) {
    filepath := filepath.Join(fs.UploadPath, filename)
    
    out, err := os.Create(filepath)
    if err != nil {
        return "", err
    }
    defer out.Close()
    
    _, err = io.Copy(out, file)
    if err != nil {
        return "", err
    }
    
    return filepath, nil
}

func (fs *FileStorage) GetFile(filepath string) (*os.File, error) {
    return os.Open(filepath)
}

func (fs *FileStorage) DeleteFile(filepath string) error {
    return os.Remove(filepath)
}

func (fs *FileStorage) FileExists(filepath string) bool {
    _, err := os.Stat(filepath)
    return !os.IsNotExist(err)
}