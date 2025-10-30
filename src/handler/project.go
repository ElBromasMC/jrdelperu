package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// HandleAdminProjectsIndex muestra la lista de proyectos
func (h *Handler) HandleAdminProjectsIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Obtener todos los proyectos
	projects, err := h.queries.ListAllProjects(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener proyectos")
	}

	return Render(c, http.StatusOK, view.AdminProjectsPage(sessionData.Username, projects))
}

// HandleAdminProjectDetail muestra el detalle de un proyecto con sus imágenes
func (h *Handler) HandleAdminProjectDetail(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	// Obtener el proyecto
	project, err := h.queries.GetProject(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusNotFound, "Proyecto no encontrado")
	}

	// Obtener imágenes del proyecto
	images, err := h.queries.ListProjectImages(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	// Obtener todas las imágenes disponibles (para el selector)
	allImages, err := h.queries.ListImages(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes disponibles")
	}

	return Render(c, http.StatusOK, view.AdminProjectDetailPage(sessionData.Username, project, images, allImages))
}

// HandleAdminProjectNewForm muestra el formulario para crear un proyecto
func (h *Handler) HandleAdminProjectNewForm(c echo.Context) error {
	sessionData := c.Get("session").(*service.SessionData)
	return Render(c, http.StatusOK, view.AdminProjectFormPage(sessionData.Username, repository.Project{}, false, "Nuevo Proyecto"))
}

// HandleAdminProjectEditForm muestra el formulario para editar un proyecto
func (h *Handler) HandleAdminProjectEditForm(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	// Obtener el proyecto
	project, err := h.queries.GetProject(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusNotFound, "Proyecto no encontrado")
	}

	return Render(c, http.StatusOK, view.AdminProjectFormPage(sessionData.Username, project, true, "Editar Proyecto"))
}

// HandleAdminProjectCreate crea un nuevo proyecto
func (h *Handler) HandleAdminProjectCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse form data
	description := c.FormValue("description")
	slug := c.FormValue("slug")
	location := c.FormValue("location")
	period := c.FormValue("period")
	areaM2Str := c.FormValue("area_m2")
	service := c.FormValue("service")
	displayOrderStr := c.FormValue("display_order")
	isVisibleStr := c.FormValue("is_visible")

	// Validar campos requeridos
	if description == "" || slug == "" || location == "" || period == "" || service == "" {
		return c.String(http.StatusBadRequest, "Todos los campos requeridos deben ser completados")
	}

	// Parse display order
	displayOrder := int32(0)
	if displayOrderStr != "" {
		order, err := strconv.ParseInt(displayOrderStr, 10, 32)
		if err == nil {
			displayOrder = int32(order)
		}
	}

	// Parse area_m2
	var areaM2 pgtype.Numeric
	if areaM2Str != "" {
		if err := areaM2.Scan(areaM2Str); err == nil {
			areaM2.Valid = true
		}
	}

	// Parse is_visible (checkbox)
	isVisible := pgtype.Bool{
		Bool:  isVisibleStr == "on",
		Valid: true,
	}

	// Crear el proyecto
	_, err := h.queries.CreateProject(ctx, repository.CreateProjectParams{
		Slug:         slug,
		Description:  description,
		Location:     location,
		Period:       period,
		AreaM2:       areaM2,
		Service:      service,
		DisplayOrder: displayOrder,
		IsVisible:    isVisible,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.String(http.StatusBadRequest, "Ya existe un proyecto con ese slug")
		}
		return c.String(http.StatusInternalServerError, "Error al crear proyecto")
	}

	// Return HTMX redirect
	c.Response().Header().Set("HX-Redirect", "/admin/projects")
	return c.NoContent(http.StatusOK)
}

// HandleAdminProjectUpdate actualiza un proyecto
func (h *Handler) HandleAdminProjectUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	// Parse form data
	description := c.FormValue("description")
	slug := c.FormValue("slug")
	location := c.FormValue("location")
	period := c.FormValue("period")
	areaM2Str := c.FormValue("area_m2")
	service := c.FormValue("service")
	displayOrderStr := c.FormValue("display_order")
	isVisibleStr := c.FormValue("is_visible")

	// Validar campos requeridos
	if description == "" || slug == "" || location == "" || period == "" || service == "" {
		return c.String(http.StatusBadRequest, "Todos los campos requeridos deben ser completados")
	}

	// Parse display order
	displayOrder := int32(0)
	if displayOrderStr != "" {
		order, err := strconv.ParseInt(displayOrderStr, 10, 32)
		if err == nil {
			displayOrder = int32(order)
		}
	}

	// Parse area_m2
	var areaM2 pgtype.Numeric
	if areaM2Str != "" {
		if err := areaM2.Scan(areaM2Str); err == nil {
			areaM2.Valid = true
		}
	}

	// Parse is_visible (checkbox)
	isVisible := pgtype.Bool{
		Bool:  isVisibleStr == "on",
		Valid: true,
	}

	// Actualizar el proyecto
	err = h.queries.UpdateProject(ctx, repository.UpdateProjectParams{
		ProjectID:    int32(projectID),
		Slug:         slug,
		Description:  description,
		Location:     location,
		Period:       period,
		AreaM2:       areaM2,
		Service:      service,
		DisplayOrder: displayOrder,
		IsVisible:    isVisible,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.String(http.StatusBadRequest, "Ya existe un proyecto con ese slug")
		}
		return c.String(http.StatusInternalServerError, "Error al actualizar proyecto")
	}

	// Return HTMX redirect
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/admin/projects/%d", projectID))
	return c.NoContent(http.StatusOK)
}

// HandleAdminProjectDelete elimina un proyecto
func (h *Handler) HandleAdminProjectDelete(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	// Eliminar el proyecto (las imágenes asociadas se eliminan en cascada)
	err = h.queries.DeleteProject(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar proyecto")
	}

	// Return HTMX redirect
	c.Response().Header().Set("HX-Redirect", "/admin/projects")
	return c.NoContent(http.StatusOK)
}

// HandleAdminProjectImageAdd agrega una imagen a un proyecto
func (h *Handler) HandleAdminProjectImageAdd(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	// Parse form data
	imageIDStr := c.FormValue("image_id")
	displayOrderStr := c.FormValue("display_order")
	isFeaturedStr := c.FormValue("is_featured")

	imageID, err := strconv.ParseInt(imageIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de imagen inválido")
	}

	displayOrder := int32(0)
	if displayOrderStr != "" {
		order, err := strconv.ParseInt(displayOrderStr, 10, 32)
		if err == nil {
			displayOrder = int32(order)
		}
	}

	isFeatured := pgtype.Bool{
		Bool:  isFeaturedStr == "on",
		Valid: true,
	}

	// Si se marca como destacada, primero desmarcar todas las demás
	if isFeatured.Bool {
		_ = h.queries.UnsetAllFeaturedProjectImages(ctx, int32(projectID))
	}

	// Agregar la imagen al proyecto
	err = h.queries.AddProjectImage(ctx, repository.AddProjectImageParams{
		ProjectID:    int32(projectID),
		ImageID:      int32(imageID),
		DisplayOrder: displayOrder,
		IsFeatured:   isFeatured,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.String(http.StatusBadRequest, "Esta imagen ya está asociada al proyecto")
		}
		return c.String(http.StatusInternalServerError, "Error al agregar imagen")
	}

	// Obtener todas las imágenes del proyecto para actualizar el grid
	images, err := h.queries.ListProjectImages(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	return Render(c, http.StatusOK, view.AdminProjectImagesGrid(int32(projectID), images))
}

// HandleAdminProjectImageDelete elimina una imagen de un proyecto
func (h *Handler) HandleAdminProjectImageDelete(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	imageIDStr := c.Param("imageId")

	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	imageID, err := strconv.ParseInt(imageIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de imagen inválido")
	}

	// Eliminar la asociación
	err = h.queries.RemoveProjectImage(ctx, repository.RemoveProjectImageParams{
		ProjectID: int32(projectID),
		ImageID:   int32(imageID),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar imagen")
	}

	// Obtener todas las imágenes del proyecto para actualizar el grid
	images, err := h.queries.ListProjectImages(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	return Render(c, http.StatusOK, view.AdminProjectImagesGrid(int32(projectID), images))
}

// HandleAdminProjectImageUpload sube una imagen y la asocia al proyecto directamente
func (h *Handler) HandleAdminProjectImageUpload(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	// Obtener archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "No se proporcionó archivo")
	}

	// Obtener otros parámetros del formulario
	displayName := c.FormValue("display_name")
	displayOrderStr := c.FormValue("display_order")
	isFeaturedStr := c.FormValue("is_featured")

	// Subir archivo usando el servicio de archivos
	result, err := h.fileService.UploadFile(ctx, file, displayName)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error al subir archivo: %v", err))
	}

	// Parse display order
	displayOrder := int32(0)
	if displayOrderStr != "" {
		order, err := strconv.ParseInt(displayOrderStr, 10, 32)
		if err == nil {
			displayOrder = int32(order)
		}
	}

	// Parse is_featured
	isFeatured := pgtype.Bool{
		Bool:  isFeaturedStr == "on",
		Valid: true,
	}

	// Si se marca como destacada, primero desmarcar todas las demás
	if isFeatured.Bool {
		_ = h.queries.UnsetAllFeaturedProjectImages(ctx, int32(projectID))
	}

	// Asociar la imagen al proyecto
	err = h.queries.AddProjectImage(ctx, repository.AddProjectImageParams{
		ProjectID:    int32(projectID),
		ImageID:      result.FileID,
		DisplayOrder: displayOrder,
		IsFeatured:   isFeatured,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al asociar imagen al proyecto")
	}

	// Obtener todas las imágenes del proyecto para actualizar el grid
	images, err := h.queries.ListProjectImages(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	return Render(c, http.StatusOK, view.AdminProjectImagesGrid(int32(projectID), images))
}

// HandleProjectImageUpdateOrder actualiza el orden de una imagen en un proyecto
func (h *Handler) HandleProjectImageUpdateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse IDs
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	imageID, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de imagen inválido")
	}

	// Parse display_order from form
	displayOrderStr := c.FormValue("display_order")
	displayOrder, err := strconv.Atoi(displayOrderStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Orden inválido")
	}

	// Update order in database
	err = h.queries.UpdateProjectImageOrder(ctx, repository.UpdateProjectImageOrderParams{
		ProjectID:    int32(projectID),
		ImageID:      int32(imageID),
		DisplayOrder: int32(displayOrder),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al actualizar orden")
	}

	// Get updated images
	images, err := h.queries.ListProjectImages(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	return Render(c, http.StatusOK, view.AdminProjectImagesGrid(int32(projectID), images))
}

// HandleProjectImageUpdateFeatured actualiza el estado de imagen destacada
func (h *Handler) HandleProjectImageUpdateFeatured(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse IDs
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de proyecto inválido")
	}

	imageID, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de imagen inválido")
	}

	// Parse is_featured from form
	isFeaturedStr := c.FormValue("is_featured")
	isFeatured := isFeaturedStr == "true"

	// If setting as featured, first unset all featured images for this project
	if isFeatured {
		err = h.queries.UnsetAllFeaturedProjectImages(ctx, int32(projectID))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error al actualizar destacadas")
		}
	}

	// Update featured status
	err = h.queries.UpdateProjectImageFeatured(ctx, repository.UpdateProjectImageFeaturedParams{
		ProjectID:  int32(projectID),
		ImageID:    int32(imageID),
		IsFeatured: pgtype.Bool{Bool: isFeatured, Valid: true},
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al actualizar estado")
	}

	// Get updated images
	images, err := h.queries.ListProjectImages(ctx, int32(projectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	return Render(c, http.StatusOK, view.AdminProjectImagesGrid(int32(projectID), images))
}
