package main

import (
	"alc/assets"
	"alc/config"
	"alc/handler"
	"alc/repository"
	"alc/service"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	if os.Getenv("ENV") == "development" {
		e.Debug = true
	}

	// Connect to database
	dbURL := os.Getenv("POSTGRESQL_URL")
	if dbURL == "" {
		log.Fatal("POSTGRESQL_URL environment variable is required")
	}

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Verify database connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	// Initialize repository
	queries := repository.New(dbpool)

	// Initialize auth service
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable is required")
	}
	authService := service.NewSessionAuthService(queries, sessionSecret)

	// Initialize file service
	fileService := service.NewFileService(queries)

	// Initialize handlers
	h := handler.New(authService, fileService, queries)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	// Static files
	e.StaticFS("/static", echo.MustSubFS(assets.Assets, "static"))

	// Uploaded files directory
	e.Static(config.UPLOADS_PATH, config.UPLOADS_SAVEDIR)

	// Page routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/nosotros", h.HandleNosotrosShow)
	e.GET("/descargas", h.HandleDescargasShow)
	e.GET("/galeria", h.HandleGaleriaShow)
	e.GET("/contacto", h.HandleContactoShow)

	// Projects routes
	e.GET("/proyectos", h.HandleProjectsShow)
	e.GET("/proyectos/:slug", h.HandleProjectDetailShow)

	// Store routes
	e.GET("/servicio", h.HandleStoreIndexShow)

	e.GET("/servicio/vidrios", h.HandleVidrioIndexShow)
	e.GET("/servicio/vidrios/:catSlug", h.HandleVidrioCategoryShow)

	e.GET("/servicio/aluminios", h.HandleAluminioIndexShow)
	e.GET("/servicio/aluminios/:catSlug", h.HandleAluminioCategoryShow)
	e.GET("/servicio/aluminios/:catSlug/:itemSlug", h.HandleAluminioItemShow)

	e.GET("/servicio/upvc", h.HandleUPVCIndexShow)
	e.GET("/servicio/upvc/:catSlug", h.HandleUPVCCategoryShow)
	e.GET("/servicio/upvc/:catSlug/:itemSlug", h.HandleUPVCItemShow)

	// Admin routes (public)
	e.GET("/admin/login", h.HandleAdminLoginShow)
	e.POST("/admin/login", h.HandleAdminLoginSubmit)
	e.GET("/admin/logout", h.HandleAdminLogout)

	// Admin routes (protected)
	admin := e.Group("/admin", h.AdminAuthMiddleware)
	admin.GET("/dashboard", h.HandleAdminDashboard)

	// File management
	admin.GET("/files", h.HandleFilesIndex)
	admin.POST("/files/upload", h.HandleFileUpload)
	admin.GET("/files/:id", h.HandleFileGet)
	admin.GET("/files/:id/edit", h.HandleFileEdit)
	admin.PUT("/files/:id/display-name", h.HandleFileUpdateDisplayName)
	admin.DELETE("/files/:id", h.HandleFileDelete)

	// Tags management
	admin.GET("/tags", h.HandleTagsIndex)
	admin.POST("/tags", h.HandleTagCreate)
	admin.PUT("/tags/:id", h.HandleTagUpdate)
	admin.DELETE("/tags/:id", h.HandleTagDelete)

	// Categories management
	admin.GET("/categories", h.HandleCategoriesIndex)
	admin.GET("/categories/new", h.HandleCategoryNewForm)
	admin.POST("/categories", h.HandleCategoryCreate)
	admin.GET("/categories/:id/edit", h.HandleCategoryEditForm)
	admin.PUT("/categories/:id", h.HandleCategoryUpdate)
	admin.DELETE("/categories/:id", h.HandleCategoryDelete)

	// Category Features management
	admin.GET("/categories/:id/features", h.HandleCategoryFeaturesIndex)
	admin.POST("/categories/:id/features", h.HandleCategoryFeatureCreate)
	admin.GET("/categories/:id/features/:featureId", h.HandleCategoryFeatureGet)
	admin.GET("/categories/:id/features/:featureId/edit", h.HandleCategoryFeatureEdit)
	admin.PUT("/categories/:id/features/:featureId", h.HandleCategoryFeatureUpdate)
	admin.DELETE("/categories/:id/features/:featureId", h.HandleCategoryFeatureDelete)

	// Items management
	admin.GET("/categories/:id/items", h.HandleItemsIndex)
	admin.GET("/categories/:id/items/new", h.HandleItemNewForm)
	admin.POST("/categories/:id/items", h.HandleItemCreate)
	admin.GET("/categories/:id/items/:itemId/edit", h.HandleItemEditForm)
	admin.PUT("/categories/:id/items/:itemId", h.HandleItemUpdate)
	admin.DELETE("/categories/:id/items/:itemId", h.HandleItemDelete)

	// Projects management
	admin.GET("/projects", h.HandleAdminProjectsIndex)
	admin.GET("/projects/new", h.HandleAdminProjectNewForm)
	admin.POST("/projects", h.HandleAdminProjectCreate)
	admin.GET("/projects/:id", h.HandleAdminProjectDetail)
	admin.GET("/projects/:id/edit", h.HandleAdminProjectEditForm)
	admin.PUT("/projects/:id", h.HandleAdminProjectUpdate)
	admin.DELETE("/projects/:id", h.HandleAdminProjectDelete)
	admin.POST("/projects/:id/images", h.HandleAdminProjectImageAdd)
	admin.POST("/projects/:id/images/upload", h.HandleAdminProjectImageUpload)
	admin.PUT("/projects/:id/images/:imageId/order", h.HandleProjectImageUpdateOrder)
	admin.PUT("/projects/:id/images/:imageId/featured", h.HandleProjectImageUpdateFeatured)
	admin.DELETE("/projects/:id/images/:imageId", h.HandleAdminProjectImageDelete)

	// Start server
	log.Println("Starting server on :8080")
	log.Fatalln(e.Start(":8080"))
}
