package handler

import (
	"alc/service"
	"alc/view"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleAdminLoginShow muestra la página de login
func (h *Handler) HandleAdminLoginShow(c echo.Context) error {
	// Si ya está autenticado, redirigir al dashboard
	session, _ := h.authService.GetSessionStore().Get(c.Request(), service.SessionName)
	if _, err := service.GetSessionData(session); err == nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	return Render(c, http.StatusOK, view.AdminLoginPage(""))
}

// AdminLoginRequest representa la solicitud de login
type AdminLoginRequest struct {
	Identifier string `form:"identifier" validate:"required"`
	Password   string `form:"password" validate:"required"`
}

// HandleAdminLoginSubmit procesa el formulario de login
func (h *Handler) HandleAdminLoginSubmit(c echo.Context) error {
	var req AdminLoginRequest
	if err := c.Bind(&req); err != nil {
		return Render(c, http.StatusOK, view.AdminLoginFormWithValues("", req.Identifier, "Por favor complete todos los campos"))
	}

	// Autenticar usuario
	sessionData, err := h.authService.Login(c.Request().Context(), req.Identifier, req.Password)
	if err != nil {
		var errorMsg string
		if errors.Is(err, service.ErrInvalidCredentials) {
			errorMsg = "Usuario o contraseña incorrectos"
		} else if errors.Is(err, service.ErrUserNotActive) {
			errorMsg = "Usuario inactivo. Contacte al administrador"
		} else {
			errorMsg = "Error al iniciar sesión. Intente nuevamente"
		}
		// Return 200 OK so HTMX swaps the content with error message, preserve identifier
		return Render(c, http.StatusOK, view.AdminLoginFormWithValues("", req.Identifier, errorMsg))
	}

	// Obtener o crear sesión
	session, err := h.authService.GetSessionStore().Get(c.Request(), service.SessionName)
	if err != nil {
		return Render(c, http.StatusOK, view.AdminLoginFormWithValues("", req.Identifier, "Error al crear sesión: "+err.Error()))
	}

	// Guardar datos en la sesión
	service.SetSessionData(session, sessionData)

	// Guardar sesión FIRST (before setting redirect header)
	if err := session.Save(c.Request(), c.Response()); err != nil {
		// Return 200 OK so user can see the error
		return Render(c, http.StatusOK, view.AdminLoginFormWithValues("", req.Identifier, "Error al guardar sesión: "+err.Error()))
	}

	// Redirigir al dashboard (AFTER session is saved)
	c.Response().Header().Set("HX-Redirect", "/admin/dashboard")
	return c.NoContent(http.StatusOK)
}

// HandleAdminLogout cierra la sesión del administrador
func (h *Handler) HandleAdminLogout(c echo.Context) error {
	session, err := h.authService.GetSessionStore().Get(c.Request(), service.SessionName)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	// Invalidar sesión
	session.Options.MaxAge = -1
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Error al cerrar sesión")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/login")
}

// HandleAdminDashboard muestra el dashboard del administrador
func (h *Handler) HandleAdminDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	// Obtener datos de sesión del contexto (añadido por el middleware)
	sessionData := c.Get("session").(*service.SessionData)

	// Obtener conteos
	categoriesCount, err := h.queries.CountCategories(ctx)
	if err != nil {
		categoriesCount = 0
	}

	itemsCount, err := h.queries.CountAllItems(ctx)
	if err != nil {
		itemsCount = 0
	}

	messagesCount, err := h.queries.CountUnreadContactSubmissions(ctx)
	if err != nil {
		messagesCount = 0
	}

	// Contar archivos (imágenes + PDFs)
	imagesCount, err := h.queries.CountImages(ctx)
	if err != nil {
		imagesCount = 0
	}

	pdfsCount, err := h.queries.CountPDFs(ctx)
	if err != nil {
		pdfsCount = 0
	}
	filesCount := imagesCount + pdfsCount

	return Render(c, http.StatusOK, view.AdminDashboard(
		sessionData.Username,
		int(categoriesCount),
		int(itemsCount),
		int(messagesCount),
		int(filesCount),
	))
}
