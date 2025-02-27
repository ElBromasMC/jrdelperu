package main

import (
	"alc/handler/page"
	"alc/handler/store"
	"alc/handler/util"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	if os.Getenv("ENV") == "development" {
		e.Debug = true
	}

	// Initialize handlers
	// TODO: Create Initializers
	ph := page.Handler{}
	sh := store.Handler{}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	// Static files
	static(e)

	// Page routes
	e.GET("/", ph.HandleIndexShow)
	e.GET("/nosotros", ph.HandleNosotrosShow)
	e.GET("/descargas", ph.HandleDescargasShow)
	e.GET("/galeria", ph.HandleGaleriaShow)
	e.GET("/contacto", ph.HandleContactoShow)

	// Store routes
	e.GET("/servicio", sh.HandleIndexShow)

	e.GET("/servicio/vidrios", sh.HandleVidrioIndexShow)
	e.GET("/servicio/vidrios/:catSlug", sh.HandleVidrioCategoryShow)

	e.GET("/servicio/aluminios", sh.HandleAluminioIndexShow)
	e.GET("/servicio/aluminios/:catSlug", sh.HandleAluminioCategoryShow)
	e.GET("/servicio/aluminios/:catSlug/:itemSlug", sh.HandleAluminioItemShow)

	e.GET("/servicio/upvc", sh.HandleUPVCIndexShow)
	e.GET("/servicio/upvc/:catSlug", sh.HandleUPVCCategoryShow)
	e.GET("/servicio/upvc/:catSlug/:itemSlug", sh.HandleUPVCItemShow)

	// Error handler
	e.HTTPErrorHandler = util.HTTPErrorHandler

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(e.Start(":" + port))
}
