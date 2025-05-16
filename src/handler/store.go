package handler

import (
	"alc/config"
	"alc/view"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleStoreIndexShow(c echo.Context) error {
	return renderOK(c, view.StoreIndex())
}

func (h *Handler) HandleVidrioIndexShow(c echo.Context) error {
	cats := config.VidrioCategories
	return renderOK(c, view.StoreVidrioIndex(cats))
}

func (h *Handler) HandleVidrioCategoryShow(c echo.Context) error {
	cat := config.VidrioCategories[0]
	is := config.MonoliticoItems
	fs := config.MonoliticoFeatures
	return renderOK(c, view.StoreVidrioCategory(cat, is, fs))
}

func (h *Handler) HandleAluminioIndexShow(c echo.Context) error {
	cats := config.AluminioCategories
	return renderOK(c, view.StoreAluminioIndex(cats))
}

func (h *Handler) HandleAluminioCategoryShow(c echo.Context) error {
	cat := config.AluminioCategories[0]
	is := config.FachadasItems
	pdf := ""
	return renderOK(c, view.StoreAluminioCategory(cat, is, pdf))
}

func (h *Handler) HandleAluminioItemShow(c echo.Context) error {
	i := config.FachadasItems[0]
	imgs := config.Imgs
	return renderOK(c, view.StoreAluminioItem(i, imgs))
}

func (h *Handler) HandleUPVCIndexShow(c echo.Context) error {
	cats := config.UPVCCategories
	return renderOK(c, view.StoreUPVCIndex(cats))
}

func (h *Handler) HandleUPVCCategoryShow(c echo.Context) error {
	cat := config.UPVCCategories[0]
	is := config.LuminaItems
	pdf := ""
	return renderOK(c, view.StoreUPVCCategory(cat, is, pdf))
}

func (h *Handler) HandleUPVCItemShow(c echo.Context) error {
	i := config.LuminaItems[0]
	imgs := config.Imgs
	return renderOK(c, view.StoreUPVCItem(i, imgs))
}
