package main

import (
	"alc/assets"
	"alc/handler"
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

	// Initialize handlers
	h := handler.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	// Static files
	e.StaticFS("/static", echo.MustSubFS(assets.Assets, "static"))

	// Page routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/nosotros", h.HandleNosotrosShow)
	e.GET("/descargas", h.HandleDescargasShow)
	e.GET("/galeria", h.HandleGaleriaShow)
	e.GET("/contacto", h.HandleContactoShow)

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

	// Start server
	log.Fatalln(e.Start(":8080"))
}
