package store

import (
	"alc/config"
	"alc/handler/util"
	"alc/view/store/vidrio"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleVidrioIndexShow(c echo.Context) error {
	cats := config.VidrioCategories
	return util.Render(c, http.StatusOK, vidrio.Index(cats))
}

func (h *Handler) HandleVidrioCategoryShow(c echo.Context) error {
	cat := config.VidrioCategories[0]
	is := config.MonoliticoItems
	fs := config.MonoliticoFeatures
	return util.Render(c, http.StatusOK, vidrio.Category(cat, is, fs))
}
