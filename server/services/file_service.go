package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"test/database"
	"test/dto"
	"test/logger"
	"test/models"
	"test/storage"
	"test/utils"
)

func SaveUploadedFile(ctx context.Context, file *multipart.FileHeader) (dto.SaveFileResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "SaveUploadedFile",
	)

	if storage.CurrentStorage == nil {
		err := errors.New("对象存储未初始化")
		logger.Sugar.Errorw("保存文件失败",
			"error", err,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.SaveFileResponse{}, err
	}

	fileInfo, err := storage.CurrentStorage.SaveFile(file)
	if err != nil {
		logger.Sugar.Errorw("保存文件失败",
			"error", err,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.SaveFileResponse{}, fmt.Errorf("保存文件失败: %w", err)
	}

	// 将文件信息保存到数据库
	dbFile := models.File{
		ID:          fileInfo.ID,
		Filename:    fileInfo.Filename,
		Size:        fileInfo.Size,
		Path:        fileInfo.Path,
		StorageType: models.StorageType(fileInfo.StorageType),
	}

	if err := database.DB.Create(&dbFile).Error; err != nil {
		logger.Sugar.Errorw("保存文件信息到数据库失败",
			"error", err,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.SaveFileResponse{}, fmt.Errorf("保存文件信息到数据库失败: %w", err)
	}

	return dto.SaveFileResponse{
		File: dto.FileInfo{
			ID:          dbFile.ID,
			Filename:    dbFile.Filename,
			Size:        dbFile.Size,
			Path:        dbFile.Path,
			StorageType: string(dbFile.StorageType),
		},
		Message: "文件上传成功",
	}, nil
}

func GetFile(ctx context.Context, fileID string) (dto.GetFileResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "GetFile",
	)

	// 从数据库中获取文件信息
	var dbFile models.File
	if err := database.DB.Where("id = ?", fileID).First(&dbFile).Error; err != nil {
		logger.Sugar.Errorw("从数据库获取文件信息失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.GetFileResponse{}, fmt.Errorf("获取文件信息失败: %w", err)
	}

	if storage.CurrentStorage == nil {
		err := errors.New("对象存储未初始化")
		logger.Sugar.Errorw("获取文件失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.GetFileResponse{}, err
	}

	reader, _, err := storage.CurrentStorage.GetFile(dbFile.Path)
	if err != nil {
		logger.Sugar.Errorw("获取文件失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.GetFileResponse{}, fmt.Errorf("获取文件失败: %w", err)
	}

	return dto.GetFileResponse{
		File: dto.FileInfo{
			ID:          dbFile.ID,
			Filename:    dbFile.Filename,
			Size:        dbFile.Size,
			Path:        dbFile.Path,
			StorageType: string(dbFile.StorageType),
		},
		Reader:  reader,
		Message: "文件获取成功",
	}, nil
}

func DeleteFile(ctx context.Context, fileID string) (dto.DeleteFileResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "DeleteFile",
	)

	// 从数据库中获取文件信息
	var dbFile models.File
	if err := database.DB.Where("id = ?", fileID).First(&dbFile).Error; err != nil {
		logger.Sugar.Errorw("从数据库获取文件信息失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.DeleteFileResponse{}, fmt.Errorf("获取文件信息失败: %w", err)
	}

	if storage.CurrentStorage == nil {
		err := errors.New("对象存储未初始化")
		logger.Sugar.Errorw("删除文件失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.DeleteFileResponse{}, err
	}

	err := storage.CurrentStorage.DeleteFile(dbFile.Path)
	if err != nil {
		logger.Sugar.Errorw("删除文件失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.DeleteFileResponse{}, fmt.Errorf("删除文件失败: %w", err)
	}

	// 从数据库中删除文件记录
	if err := database.DB.Delete(&dbFile).Error; err != nil {
		logger.Sugar.Errorw("从数据库删除文件记录失败",
			"error", err,
			"fileID", fileID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.DeleteFileResponse{}, fmt.Errorf("删除文件记录失败: %w", err)
	}

	return dto.DeleteFileResponse{
		Message: "文件删除成功",
	}, nil
}

func ListFiles(ctx context.Context) (dto.ListFilesResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "ListFiles",
	)

	var dbFiles []models.File
	if err := database.DB.Find(&dbFiles).Error; err != nil {
		logger.Sugar.Errorw("从数据库获取文件列表失败",
			"error", err,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.ListFilesResponse{}, fmt.Errorf("获取文件列表失败: %w", err)
	}

	files := make([]dto.FileInfo, len(dbFiles))
	for i, dbFile := range dbFiles {
		files[i] = dto.FileInfo{
			ID:          dbFile.ID,
			Filename:    dbFile.Filename,
			Size:        dbFile.Size,
			Path:        dbFile.Path,
			StorageType: string(dbFile.StorageType),
		}
	}

	return dto.ListFilesResponse{
		Files:   files,
		Message: "获取文件列表成功",
	}, nil
}
