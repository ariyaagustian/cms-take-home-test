// internal/transport/http/handler/entry_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cms/server/internal/model"
	"cms/server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EntryHandler struct{ repo repository.EntryRepository }

func NewEntryHandler(repo repository.EntryRepository) *EntryHandler {
	return &EntryHandler{repo: repo}
}

// POST /api/entries/:slug
func (h *EntryHandler) Create(c *gin.Context) {
	slug := c.Param("slug")
	var in struct {
		Slug   string          `json:"slug"`
		Status string          `json:"status"`
		Data   json.RawMessage `json:"data"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var editorID *uuid.UUID
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				editorID = &id
			}
		}
	}
	e := &model.Entry{Slug: in.Slug, Status: in.Status}
	if e.Status == "" {
		e.Status = "draft"
	}
	if err := h.repo.Create(c.Request.Context(), slug, e, in.Data, editorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": e})
}

// GET /api/entries/:slug?limit=20&offset=0
func (h *EntryHandler) List(c *gin.Context) {
	slug := c.Param("slug")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.repo.List(c.Request.Context(), slug, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total, "limit": limit, "offset": offset})
}

// GET /api/entries/:slug/:id
func (h *EntryHandler) Detail(c *gin.Context) {
	slug := c.Param("slug")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	item, err := h.repo.Get(c.Request.Context(), slug, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// PUT /api/entries/:slug/:id
func (h *EntryHandler) Update(c *gin.Context) {
	slug := c.Param("slug")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	var in struct {
		Status *string         `json:"status"`
		Data   json.RawMessage `json:"data"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var editorID *uuid.UUID
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if eid, err := uuid.Parse(s); err == nil {
				editorID = &eid
			}
		}
	}
	if err := h.repo.Update(c.Request.Context(), slug, id, in.Data, in.Status, editorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DELETE /api/entries/:slug/:id
func (h *EntryHandler) Delete(c *gin.Context) {
	slug := c.Param("slug")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	var editorID *uuid.UUID
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if eid, err := uuid.Parse(s); err == nil {
				editorID = &eid
			}
		}
	}
	if err := h.repo.Delete(c.Request.Context(), slug, id, editorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /api/entries/:slug/:id/publish
func (h *EntryHandler) Publish(c *gin.Context) {
	slug := c.Param("slug")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	var editorID *uuid.UUID
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if eid, err := uuid.Parse(s); err == nil {
				editorID = &eid
			}
		}
	}
	if err := h.repo.Publish(c.Request.Context(), slug, id, time.Now(), editorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// POST /api/entries/:slug/:id/rollback/:version
func (h *EntryHandler) Rollback(c *gin.Context) {
	slug := c.Param("slug")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	ver, err := strconv.Atoi(c.Param("version"))
	if err != nil || ver <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}
	var editorID *uuid.UUID
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if eid, err := uuid.Parse(s); err == nil {
				editorID = &eid
			}
		}
	}
	if err := h.repo.Rollback(c.Request.Context(), slug, id, ver, editorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
