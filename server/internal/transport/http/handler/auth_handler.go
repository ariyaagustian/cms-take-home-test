package handler

import (
	"net/http"
	"time"

	"cms/server/internal/config"
	"cms/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg config.Config
	svc service.AuthService
}

func NewAuthHandler(cfg config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
		svc: service.NewAuthService(cfg, db),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// ✅ validasi lewat service
	user, err := h.svc.Authenticate(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ✅ generate token dengan role
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	s, _ := t.SignedString([]byte(h.cfg.JWTSecret))

	// ✅ response
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token":      s,
		"token_type": "Bearer",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var in struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), in.Name, in.Email, in.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
