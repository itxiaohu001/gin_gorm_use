package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"test/database"
	"test/models"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (m *MinioStorage) SaveFile(file *multipart.FileHeader) (FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return FileInfo{}, err
	}
	defer src.Close()

	fileID := uuid.New().String()
	objectName := fmt.Sprintf("%s%s", fileID, filepath.Ext(file.Filename))

	_, err = m.client.PutObject(context.Background(), m.bucketName, objectName, src, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	if err != nil {
		return FileInfo{}, err
	}

	// 保存文件元数据到数据库
	fileModel := models.File{
		Filename:    file.Filename,
		Size:        file.Size,
		Path:        objectName,
		StorageType: models.StorageTypeMinio, // 设置存储类型为 MinIO
	}
	if err := database.DB.Create(&fileModel).Error; err != nil {
		return FileInfo{}, fmt.Errorf("保存文件元数据失败: %w", err)
	}

	return FileInfo{
		ID:          fileID,
		Filename:    file.Filename,
		Size:        file.Size,
		StorageType: models.StorageTypeMinio,
	}, nil
}

func (m *MinioStorage) GetFile(fileID string) (io.Reader, string, error) {
	var fileModel models.File
	if err := database.DB.Where("id = ?", fileID).First(&fileModel).Error; err != nil {
		return nil, "", &FileNotFoundError{FileID: fileID}
	}

	if fileModel.StorageType != models.StorageTypeMinio {
		return nil, "", &StorageTypeError{
			Expected: models.StorageTypeMinio,
			Actual:   fileModel.StorageType,
		}
	}

	object, err := m.client.GetObject(context.Background(), m.bucketName, fileModel.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("获取文件失败: %w", err)
	}

	return object, fileModel.Filename, nil
}

func (m *MinioStorage) DeleteFile(fileID string) error {
	var fileModel models.File
	if err := database.DB.Where("id = ?", fileID).First(&fileModel).Error; err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	err := m.client.RemoveObject(context.Background(), m.bucketName, fileModel.Path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	if err := database.DB.Delete(&fileModel).Error; err != nil {
		return fmt.Errorf("删除文件元数据失败: %w", err)
	}

	return nil
}

func (m *MinioStorage) ListFiles() ([]FileInfo, error) {
	var files []models.File
	if err := database.DB.Find(&files).Error; err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %w", err)
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Filename: file.Filename,
			Size:     file.Size,
		})
	}

	return fileInfos, nil
}

func (m *MinioStorage) GetStorageType() models.StorageType {
	return models.StorageTypeMinio
}
