package page

import (
	"alc/handler/util"
	"alc/view/page"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Index())
}

func (h *Handler) HandleNosotrosShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Nosotros())
}

func (h *Handler) HandleDescargasShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Descargas())
}

func (h *Handler) HandleGaleriaShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Galeria())
}

func (h *Handler) HandleContactoShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Contacto())
}
