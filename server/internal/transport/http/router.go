package http

import (
	"net/http"

	"cms/server/internal/config"
	"cms/server/internal/repository"
	"cms/server/internal/transport/http/handler"
	"cms/server/internal/transport/http/middleware"
	"cms/server/pkg/minio"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger UI (opsional)
	// r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public API
	api := r.Group("/api")

	// Auth endpoints
	auth := handler.NewAuthHandler(cfg, db)
	api.POST("/auth/login", auth.Login)
	api.POST("/auth/register", auth.Register) // demo

	// Protected routes (require JWT)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))

	auditRepo := repository.NewAuditRepository(db)

	{
		// ContentType handler (editor & admin)
		ctRepo := repository.NewContentTypeRepository(db)
		ct := handler.NewContentTypeHandler(ctRepo)
		ctGroup := protected.Group("/content-types")
		ctGroup.Use(middleware.RequireRole("Editor", "Admin"))
		ctGroup.POST("", ct.Create)
		ctGroup.GET("", ct.List)
		ctGroup.GET("/:id", ct.Detail)
		ctGroup.PUT("/:id", ct.Update)
		ctGroup.DELETE("/:id", ct.Delete)
		ctGroup.POST("/:id/fields", ct.AddField)

		// Entry handler
		entryRepo := repository.NewEntryRepository(db, auditRepo)
		entry := handler.NewEntryHandler(entryRepo)
		entryGroup := protected.Group("/entries/:slug")
		entryGroup.Use(middleware.RequireRole("Editor", "Admin"))
		entryGroup.POST("", entry.Create)
		entryGroup.GET("", entry.List)
		entryGroup.GET("/:id", entry.Detail)
		entryGroup.PUT("/:id", entry.Update)
		entryGroup.DELETE("/:id", entry.Delete)
		entryGroup.POST("/:id/publish", entry.Publish)
		entryGroup.POST("/:id/rollback/:version", entry.Rollback)

		// Media handler
		minioClient := minio.New(
			cfg.MinIOEndpoint,
			cfg.MinIOAccessKey,
			cfg.MinIOSecretKey,
			cfg.MinIOBucket,
			cfg.MinIOUseSSL,
		)
		mediaRepo := repository.NewMediaRepository(db)
		media := handler.NewMediaHandler(minioClient, mediaRepo)

		mediaGroup := protected.Group("/media")
		mediaGroup.Use(middleware.RequireRole("Editor", "Admin"))
		mediaGroup.POST("", media.Upload)
		mediaGroup.GET("/preview/:id", media.Preview)
		mediaGroup.GET("", media.List)
		mediaGroup.DELETE("/:id", media.Delete)

		// // Admin-only endpoints
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("Admin"))

		roleRepo := repository.NewRoleRepository(db)
		role := handler.NewRoleHandler(roleRepo) // atau repo khusus
		admin.GET("/roles", role.List)
		admin.POST("/roles", role.Create)

		userRepo := repository.NewUserRepository(db)
		user := handler.NewUserHandler(userRepo)
		admin.GET("/users", user.List)
		admin.GET("/users/:id/roles", user.GetRoles)
		admin.POST("/users/:id/roles", user.SetRoles)
	}

	entryRepo := repository.NewEntryRepository(db, auditRepo)
	publicHandler := handler.NewPublicHandler(entryRepo)
	r.GET("/api/public/:slug", publicHandler.ListPublished)
	r.GET("/api/public/:slug/:id", publicHandler.GetPublished)

	return r
}
