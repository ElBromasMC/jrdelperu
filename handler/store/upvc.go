package store

import (
	"alc/config"
	"alc/handler/util"
	"alc/view/store/upvc"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleUPVCIndexShow(c echo.Context) error {
	cats := config.UPVCCategories
	return util.Render(c, http.StatusOK, upvc.Index(cats))
}

func (h *Handler) HandleUPVCCategoryShow(c echo.Context) error {
	cat := config.UPVCCategories[0]
	is := config.LuminaItems
	pdf := ""
	return util.Render(c, http.StatusOK, upvc.Category(cat, is, pdf))
}

func (h *Handler) HandleUPVCItemShow(c echo.Context) error {
	i := config.LuminaItems[0]
	imgs := config.Imgs
	return util.Render(c, http.StatusOK, upvc.Item(i, imgs))
}
