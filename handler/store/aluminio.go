package store

import (
	"alc/handler/util"
	"alc/view/store/aluminio"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleAluminioIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, aluminio.Index())
}
