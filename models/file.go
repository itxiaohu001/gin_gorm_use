package models

import "gorm.io/gorm"

type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeMinio StorageType = "minio"
)

type File struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Filename    string
	Size        int64
	Path        string
	StorageType StorageType
}
