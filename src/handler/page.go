package handler

import (
	"alc/view"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	return renderOK(c, view.Index())
}

func (h *Handler) HandleNosotrosShow(c echo.Context) error {
	return renderOK(c, view.Nosotros())
}

func (h *Handler) HandleDescargasShow(c echo.Context) error {
	return renderOK(c, view.Descargas())
}

func (h *Handler) HandleGaleriaShow(c echo.Context) error {
	return renderOK(c, view.Galeria())
}

func (h *Handler) HandleContactoShow(c echo.Context) error {
	return renderOK(c, view.Contacto())
}
