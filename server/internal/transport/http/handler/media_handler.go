package handler

import (
	"cms/server/internal/model"
	"cms/server/internal/repository"
	"cms/server/pkg/minio"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaHandler struct {
	MinioClient *minio.Client
	Repository  repository.MediaRepository
}

func NewMediaHandler(m *minio.Client, repo repository.MediaRepository) *MediaHandler {
	return &MediaHandler{
		MinioClient: m,
		Repository:  repo,
	}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file tidak ditemukan"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuka file"})
		return
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membaca file"})
		return
	}

	objectName := uuid.New().String() + filepath.Ext(fileHeader.Filename)

	err = h.MinioClient.Upload(c.Request.Context(), objectName, fileHeader.Header.Get("Content-Type"), buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload ke minio gagal"})
		return
	}

	meta := map[string]any{
		"source":        "upload",
		"uploaded_from": c.ClientIP(),
	}
	metaJSON, _ := json.Marshal(meta)

	asset := &model.MediaAsset{
		Filename:  fileHeader.Filename,
		Mime:      fileHeader.Header.Get("Content-Type"),
		SizeBytes: fileHeader.Size,
		URL:       objectName,
		Meta:      metaJSON,
		CreatedAt: time.Now(),
	}

	if err := h.Repository.Save(c.Request.Context(), asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       asset.ID,
		"filename": asset.Filename,
		"url":      asset.URL,
		"meta":     meta,
	})
}

func (h *MediaHandler) Preview(c *gin.Context) {
	id := c.Param("id")

	asset, err := h.Repository.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file tidak ditemukan"})
		return
	}

	url, err := h.MinioClient.GenerateSignedURL(c.Request.Context(), asset.URL, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal generate signed url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preview_url": url})
}

func (h *MediaHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	items, total, err := h.Repository.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *MediaHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Hapus dari MinIO (opsional)
	asset, err := h.Repository.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media tidak ditemukan"})
		return
	}
	_ = h.MinioClient.Delete(c.Request.Context(), asset.URL) // ignore error

	// Hapus dari DB
	if err := h.Repository.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal hapus media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "media berhasil dihapus"})
}
