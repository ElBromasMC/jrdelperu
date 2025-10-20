package handler

import (
	"alc/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AdminAuthMiddleware verifica que el usuario esté autenticado como admin
func (h *Handler) AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener sesión
		session, err := h.authService.GetSessionStore().Get(c.Request(), service.SessionName)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/login")
		}

		// Validar datos de sesión
		sessionData, err := service.GetSessionData(session)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/login")
		}

		// Almacenar datos de sesión en el contexto para uso posterior
		c.Set("session", sessionData)

		return next(c)
	}
}
