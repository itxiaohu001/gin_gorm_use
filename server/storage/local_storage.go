package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"test/models"

	"github.com/google/uuid"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) SaveFile(file *multipart.FileHeader) (FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return FileInfo{}, err
	}
	defer src.Close()

	fileID := uuid.New().String()
	fileName := fmt.Sprintf("%s%s", fileID, filepath.Ext(file.Filename))
	filePath := filepath.Join(s.basePath, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return FileInfo{}, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		ID:          fileID,
		Filename:    file.Filename,
		Size:        file.Size,
		StorageType: models.StorageTypeLocal,
	}, nil
}

func (s *LocalStorage) GetFile(fileID string) (io.Reader, string, error) {
	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, "", err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), fileID) {
			filePath := filepath.Join(s.basePath, file.Name())
			f, err := os.Open(filePath)
			if err != nil {
				return nil, "", err
			}
			return f, file.Name(), nil
		}
	}

	return nil, "", fmt.Errorf("file not found")
}

func (s *LocalStorage) DeleteFile(fileID string) error {
	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), fileID) {
			return os.Remove(filepath.Join(s.basePath, file.Name()))
		}
	}

	return fmt.Errorf("file not found")
}

func (s *LocalStorage) ListFiles() ([]FileInfo, error) {
	var fileInfos []FileInfo

	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}
			fileID := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			fileInfos = append(fileInfos, FileInfo{
				ID:          fileID,
				Filename:    file.Name(),
				Size:        info.Size(),
				StorageType: models.StorageTypeLocal,
			})
		}
	}

	return fileInfos, nil
}

func (s *LocalStorage) GetStorageType() models.StorageType {
	return models.StorageTypeLocal
}
