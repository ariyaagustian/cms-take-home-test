package handler

import (
	"cms/server/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(r repository.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) GetRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	roles, err := h.repo.GetRoles(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *UserHandler) SetRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var body struct {
		Roles []int `json:"roles"` // disesuaikan dengan payload frontend
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload tidak valid"})
		return
	}

	if len(body.Roles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roles tidak boleh kosong"})
		return
	}

	if err := h.repo.SetRoles(c.Request.Context(), id, body.Roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengatur role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil mengatur role"})
}
