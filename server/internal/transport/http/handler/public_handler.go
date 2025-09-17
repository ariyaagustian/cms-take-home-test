package handler

import (
	"net/http"
	"strconv"

	"cms/server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PublicHandler struct{ repo repository.EntryRepository }

func NewPublicHandler(repo repository.EntryRepository) *PublicHandler {
	return &PublicHandler{repo: repo}
}

// GET /api/public/:slug
func (h *PublicHandler) ListPublished(c *gin.Context) {
	slug := c.Param("slug")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sort := c.DefaultQuery("sort", "-published_at")

	items, total, err := h.repo.ListPublished(c.Request.Context(), slug, limit, offset, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GET /api/public/:slug/:id
func (h *PublicHandler) GetPublished(c *gin.Context) {
	slug := c.Param("slug")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	item, err := h.repo.GetPublished(c.Request.Context(), slug, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
