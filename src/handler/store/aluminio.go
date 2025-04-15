package store

import (
	"alc/config"
	"alc/handler/util"
	"alc/view/store/aluminio"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleAluminioIndexShow(c echo.Context) error {
	cats := config.AluminioCategories
	return util.Render(c, http.StatusOK, aluminio.Index(cats))
}

func (h *Handler) HandleAluminioCategoryShow(c echo.Context) error {
	cat := config.AluminioCategories[0]
	is := config.FachadasItems
	pdf := ""
	return util.Render(c, http.StatusOK, aluminio.Category(cat, is, pdf))
}

func (h *Handler) HandleAluminioItemShow(c echo.Context) error {
	i := config.FachadasItems[0]
	imgs := config.Imgs
	return util.Render(c, http.StatusOK, aluminio.Item(i, imgs))
}
