package storage

import (
	"fmt"
	"test/models"
)

type StorageTypeError struct {
	Expected models.StorageType
	Actual   models.StorageType
}

func (e *StorageTypeError) Error() string {
	return fmt.Sprintf("存储类型不匹配: 期望 %s 存储，实际为 %s", e.Expected, e.Actual)
}

type FileNotFoundError struct {
	FileID string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("文件不存在: %s", e.FileID)
}
