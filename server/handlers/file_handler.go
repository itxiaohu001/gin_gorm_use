package handlers

import (
	"fmt"
	"io"
	"net/http"
	"test/dto"
	"test/services"

	"github.com/gin-gonic/gin"
)

// UploadFile godoc
// @Summary 上传文件
// @Description 上传一个文件到服务器
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 200 {object} dto.SaveFileResponse "文件上传成功"
// @Failure 400 {object} dto.ErrorResponse "无效的文件"
// @Failure 500 {object} dto.ErrorResponse "文件上传失败"
// @Security ApiKeyAuth
// @Router /files/upload [post]
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "无效的文件", Message: err.Error()})
		return
	}

	resp, err := services.SaveUploadedFile(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "文件上传失败", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DownloadFile godoc
// @Summary 下载文件
// @Description 从服务器下载指定的文件
// @Tags files
// @Produce octet-stream
// @Param fileID path string true "文件ID"
// @Success 200 {file} binary "文件内容"
// @Failure 404 {object} dto.ErrorResponse "文件不存在"
// @Failure 500 {object} dto.ErrorResponse "下载文件失败"
// @Security ApiKeyAuth
// @Router /files/download/{fileID} [get]
func DownloadFile(c *gin.Context) {
	fileID := c.Param("fileID")

	resp, err := services.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "文件不存在", Message: err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.File.Filename))
	c.Header("Content-Type", "application/octet-stream")

	_, err = io.Copy(c.Writer, resp.Reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "下载文件失败", Message: err.Error()})
		return
	}
}

// DeleteFile godoc
// @Summary 删除文件
// @Description 从服务器删除指定的文件
// @Tags files
// @Produce json
// @Param fileID path string true "文件ID"
// @Success 200 {object} dto.DeleteFileResponse "文件删除成功"
// @Failure 404 {object} dto.ErrorResponse "文件不存在"
// @Failure 500 {object} dto.ErrorResponse "删除文件失败"
// @Security ApiKeyAuth
// @Router /files/{fileID} [delete]
func DeleteFile(c *gin.Context) {
	fileID := c.Param("fileID")

	resp, err := services.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "删除文件失败", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListFiles godoc
// @Summary 获取文件列表
// @Description 获取当前用户的所有文件列表
// @Tags files
// @Produce json
// @Success 200 {object} dto.ListFilesResponse "文件列表"
// @Failure 500 {object} dto.ErrorResponse "获取文件列表失败"
// @Security ApiKeyAuth
// @Router /files [get]
func ListFiles(c *gin.Context) {
	resp, err := services.ListFiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "获取文件列表失败", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
