// internal/transport/http/handler/content_type_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"cms/server/internal/model"
	"cms/server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContentTypeHandler struct {
	repo repository.ContentTypeRepository
}

func NewContentTypeHandler(repo repository.ContentTypeRepository) *ContentTypeHandler {
	return &ContentTypeHandler{repo: repo}
}

func (h *ContentTypeHandler) Create(c *gin.Context) {
	var in struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct := model.ContentType{
		Name: in.Name,
		Slug: in.Slug,
	}
	if err := h.repo.Create(c.Request.Context(), &ct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ct})
}

func (h *ContentTypeHandler) List(c *gin.Context) {
	list, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *ContentTypeHandler) Detail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	ct, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ct})
}

func (h *ContentTypeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	var in struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, in.Name, in.Slug); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *ContentTypeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ContentTypeHandler) AddField(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	var in struct {
		Name    string         `json:"name" binding:"required"`
		Kind    string         `json:"kind" binding:"required"`
		Options map[string]any `json:"options"` // <- object JSON, bukan string
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi kind (boleh tambahkan image & wysiwyg)
	allowedKinds := map[string]bool{
		"text": true, "string": true, "number": true, "bool": true, "date": true,
		"json": true, "select": true, "image": true, "wysiwyg": true,
	}
	if !allowedKinds[in.Kind] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid field kind"})
		return
	}

	// marshal options ke JSON (default: {} bila nil)
	var optsBytes []byte
	if in.Options == nil {
		optsBytes = []byte(`{}`)
	} else {
		optsBytes, err = json.Marshal(in.Options)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid options json"})
			return
		}
	}

	field := model.ContentField{
		ContentTypeID: id,
		Name:          in.Name,
		Kind:          in.Kind,
		Options:       optsBytes, // json.RawMessage
	}

	if err := h.repo.AddField(c.Request.Context(), &field); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": field})
}
