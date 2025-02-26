package store

import (
	"alc/handler/util"
	"alc/view/store/vidrio"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleVidrioIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, vidrio.Index())
}
