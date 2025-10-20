package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// HandleCategoryFeaturesIndex muestra las características de una categoría
func (h *Handler) HandleCategoryFeaturesIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Get category
	category, err := h.queries.GetCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Get features for this category
	features, err := h.queries.ListCategoryFeatures(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar características")
	}

	return Render(c, http.StatusOK, view.AdminCategoryFeaturesPage(sessionData.Username, category, features))
}

// HandleCategoryFeatureGet devuelve una fila de característica en modo visualización
func (h *Handler) HandleCategoryFeatureGet(c echo.Context) error {
	ctx := c.Request().Context()

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Get feature ID
	featureIDStr := c.Param("featureId")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de característica inválido")
	}

	// Get feature
	feature, err := h.queries.GetCategoryFeature(ctx, int32(featureID))
	if err != nil {
		return c.String(http.StatusNotFound, "Característica no encontrada")
	}

	return Render(c, http.StatusOK, view.AdminCategoryFeatureRow(int32(categoryID), feature))
}

// HandleCategoryFeatureEdit devuelve una fila de característica en modo edición
func (h *Handler) HandleCategoryFeatureEdit(c echo.Context) error {
	ctx := c.Request().Context()

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Get feature ID
	featureIDStr := c.Param("featureId")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de característica inválido")
	}

	// Get feature
	feature, err := h.queries.GetCategoryFeature(ctx, int32(featureID))
	if err != nil {
		return c.String(http.StatusNotFound, "Característica no encontrada")
	}

	return Render(c, http.StatusOK, view.AdminCategoryFeatureRowEdit(int32(categoryID), feature))
}

// HandleCategoryFeatureCreate crea una nueva característica
func (h *Handler) HandleCategoryFeatureCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Parse form data
	name := c.FormValue("name")
	description := c.FormValue("description")

	// Create feature
	feature, err := h.queries.CreateCategoryFeature(ctx, repository.CreateCategoryFeatureParams{
		CategoryID:  int32(categoryID),
		Name:        name,
		Description: description,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al crear característica: %v", err))
	}

	// Return the new row (with OOB to remove empty state if it exists)
	return Render(c, http.StatusOK, view.AdminCategoryFeatureRowWithEmptyRemoval(int32(categoryID), feature))
}

// HandleCategoryFeatureUpdate actualiza una característica
func (h *Handler) HandleCategoryFeatureUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Get feature ID
	featureIDStr := c.Param("featureId")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de característica inválido")
	}

	// Parse form data
	name := c.FormValue("name")
	description := c.FormValue("description")

	// Update feature
	err = h.queries.UpdateCategoryFeature(ctx, repository.UpdateCategoryFeatureParams{
		FeatureID:   int32(featureID),
		Name:        name,
		Description: description,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al actualizar característica: %v", err))
	}

	// Get updated feature
	feature, err := h.queries.GetCategoryFeature(ctx, int32(featureID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener característica actualizada")
	}

	// Return the updated row
	return Render(c, http.StatusOK, view.AdminCategoryFeatureRow(int32(categoryID), feature))
}

// HandleCategoryFeatureDelete elimina una característica
func (h *Handler) HandleCategoryFeatureDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Get feature ID
	featureIDStr := c.Param("featureId")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de característica inválido")
	}

	// Delete feature
	err = h.queries.DeleteCategoryFeature(ctx, int32(featureID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al eliminar característica: %v", err))
	}

	// Return empty response (HTMX will remove the row)
	return c.NoContent(http.StatusOK)
}
