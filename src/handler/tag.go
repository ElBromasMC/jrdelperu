package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// HandleTagsIndex muestra la lista de tags
func (h *Handler) HandleTagsIndex(c echo.Context) error {
	sessionData := c.Get("session").(*service.SessionData)

	tags, err := h.queries.ListCategoryTags(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener tags")
	}

	return Render(c, http.StatusOK, view.AdminTagsIndex(sessionData.Username, tags))
}

// HandleTagCreate procesa la creación de un tag
func (h *Handler) HandleTagCreate(c echo.Context) error {
	tagName := c.FormValue("tag_name")
	positionStr := c.FormValue("position_num")

	if tagName == "" || positionStr == "" {
		return c.String(http.StatusBadRequest, "Nombre y posición son requeridos")
	}

	position, err := strconv.ParseInt(positionStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Posición inválida")
	}

	tag, err := h.queries.CreateCategoryTag(c.Request().Context(), repository.CreateCategoryTagParams{
		TagName:     tagName,
		PositionNum: int32(position),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al crear tag")
	}

	// Retornar la nueva fila de la tabla
	return Render(c, http.StatusOK, view.TagRow(tag))
}

// HandleTagUpdate procesa la actualización de un tag
func (h *Handler) HandleTagUpdate(c echo.Context) error {
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de tag inválido")
	}

	tagName := c.FormValue("tag_name")
	positionStr := c.FormValue("position_num")

	if tagName == "" || positionStr == "" {
		return c.String(http.StatusBadRequest, "Nombre y posición son requeridos")
	}

	position, err := strconv.ParseInt(positionStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Posición inválida")
	}

	err = h.queries.UpdateCategoryTag(c.Request().Context(), repository.UpdateCategoryTagParams{
		TagID:       int32(tagID),
		TagName:     tagName,
		PositionNum: int32(position),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al actualizar tag")
	}

	// Obtener el tag actualizado
	tag, err := h.queries.GetCategoryTag(c.Request().Context(), int32(tagID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener tag")
	}

	// Retornar la fila actualizada
	return Render(c, http.StatusOK, view.TagRow(tag))
}

// HandleTagDelete elimina un tag
func (h *Handler) HandleTagDelete(c echo.Context) error {
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de tag inválido")
	}

	err = h.queries.DeleteCategoryTag(c.Request().Context(), int32(tagID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar tag")
	}

	return c.NoContent(http.StatusOK)
}
