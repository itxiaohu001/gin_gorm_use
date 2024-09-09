package dto

import "io"

type FileInfo struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	Path        string `json:"path"`
	StorageType string `json:"storage_type"`
}

type SaveFileResponse struct {
	File    FileInfo `json:"file"`
	Message string   `json:"message"`
}

type GetFileResponse struct {
	File     FileInfo `json:"file"`
	Reader   io.Reader
	Message  string `json:"message"`
}

type DeleteFileResponse struct {
	Message string `json:"message"`
}

type ListFilesResponse struct {
	Files   []FileInfo `json:"files"`
	Message string     `json:"message"`
}