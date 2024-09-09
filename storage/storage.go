package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"test/config"
	"test/logger"
	"test/models"
)

// FileInfo 包含文件的元数据
type FileInfo struct {
	ID          string
	Filename    string
	Size        int64
	Path        string
	StorageType models.StorageType
}

// ObjectStorage 定义了对象存储的接口
type ObjectStorage interface {
	// SaveFile 保存文件到对象存储，返回文件信息
	SaveFile(file *multipart.FileHeader) (FileInfo, error)
	// GetFile 从对象存储获取文件
	GetFile(fileID string) (io.Reader, string, error)
	// DeleteFile 从对象存储删除文件
	DeleteFile(fileID string) error
	// ListFiles 获取文件列表
	ListFiles() ([]FileInfo, error)
	// GetStorageType 返回当前存储类型
	GetStorageType() models.StorageType
}

// CurrentStorage 全局变量，用于存储当前使用的对象存储实现
var CurrentStorage ObjectStorage

// InitStorage 初始化存储
func InitStorage() error {
	storageType := config.AppConfig.Storage.Type
	switch models.StorageType(storageType) {
	case models.StorageTypeLocal:
		basePath := config.AppConfig.Storage.LocalPath
		if basePath == "" {
			return fmt.Errorf("本地存储路径未配置")
		}
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("创建本地存储目录失败: %v", err)
		}
		var err error
		CurrentStorage, err = NewLocalStorage(basePath)
		if err != nil {
			return err
		}
		logger.Sugar.Info("本地存储初始化成功")
	// 可以在这里添加其他存储类型的初始化，如 S3、阿里云 OSS 等
	default:
		return fmt.Errorf("不支持的存储类型: %s", storageType)
	}
	return nil
}
