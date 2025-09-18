package handler

import (
	"net/http"

	"cms/server/internal/repository"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	repo repository.RoleRepository
}

func NewRoleHandler(repo repository.RoleRepository) *RoleHandler {
	return &RoleHandler{repo: repo}
}

// GET /admin/roles
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data role"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// POST /admin/roles
func (h *RoleHandler) Create(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama role wajib diisi"})
		return
	}

	if err := h.repo.Create(c.Request.Context(), body.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat role"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "berhasil membuat role"})
}
