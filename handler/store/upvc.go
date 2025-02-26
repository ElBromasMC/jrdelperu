package store

import (
	"alc/handler/util"
	"alc/view/store/upvc"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleUpvcIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, upvc.Index())
}
