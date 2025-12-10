package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// HandleItemsIndex muestra la lista de items de una categoría
func (h *Handler) HandleItemsIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Obtener la categoría
	category, err := h.queries.GetCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Obtener items de la categoría
	items, err := h.queries.ListItemsByCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener items")
	}

	// Cargar imágenes de los items
	imagesMap := make(map[int32]repository.StaticFile)
	for _, item := range items {
		if item.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, item.ImageID.Int32)
			if err == nil {
				imagesMap[item.ImageID.Int32] = img
			}
		}
	}

	return Render(c, http.StatusOK, view.AdminItemsPage(sessionData.Username, category, items, imagesMap))
}

// HandleItemNewForm muestra el formulario para crear un item
func (h *Handler) HandleItemNewForm(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Obtener la categoría
	category, err := h.queries.GetCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Obtener imágenes disponibles
	images, err := h.queries.ListImages(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	// Obtener PDFs disponibles
	pdfs, err := h.queries.ListPDFs(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener PDFs")
	}

	return Render(c, http.StatusOK, view.AdminItemFormPage(sessionData.Username, category, repository.Item{}, images, pdfs, false))
}

// HandleItemEditForm muestra el formulario para editar un item
func (h *Handler) HandleItemEditForm(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de item inválido")
	}

	// Obtener la categoría
	category, err := h.queries.GetCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Obtener el item
	item, err := h.queries.GetItem(ctx, int32(itemID))
	if err != nil {
		return c.String(http.StatusNotFound, "Item no encontrado")
	}

	// Verificar que el item pertenece a la categoría
	if item.CategoryID != int32(categoryID) {
		return c.String(http.StatusBadRequest, "El item no pertenece a esta categoría")
	}

	// Obtener imágenes disponibles
	images, err := h.queries.ListImages(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	// Obtener PDFs disponibles
	pdfs, err := h.queries.ListPDFs(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener PDFs")
	}

	return Render(c, http.StatusOK, view.AdminItemFormPage(sessionData.Username, category, item, images, pdfs, true))
}

// HandleItemCreate crea un nuevo item
func (h *Handler) HandleItemCreate(c echo.Context) error {
	ctx := c.Request().Context()

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Verificar que la categoría existe
	_, err = h.queries.GetCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Parsear formulario
	slug := c.FormValue("slug")
	name := c.FormValue("name")
	description := c.FormValue("description")
	longDescription := c.FormValue("long_description")

	// Validación básica
	if name == "" || slug == "" {
		return c.String(http.StatusBadRequest, "Nombre y slug son requeridos")
	}

	// Parsear image_id opcional
	var imageID pgtype.Int4
	if imageIDStr := c.FormValue("image_id"); imageIDStr != "" {
		if id, err := strconv.ParseInt(imageIDStr, 10, 32); err == nil {
			imageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Parsear secondary_image_id opcional
	var secondaryImageID pgtype.Int4
	if secondaryImageIDStr := c.FormValue("secondary_image_id"); secondaryImageIDStr != "" {
		if id, err := strconv.ParseInt(secondaryImageIDStr, 10, 32); err == nil {
			secondaryImageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Parsear pdf_id opcional
	var pdfID pgtype.Int4
	if pdfIDStr := c.FormValue("pdf_id"); pdfIDStr != "" {
		if id, err := strconv.ParseInt(pdfIDStr, 10, 32); err == nil {
			pdfID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Crear item
	_, err = h.queries.CreateItem(ctx, repository.CreateItemParams{
		CategoryID:       int32(categoryID),
		Slug:             slug,
		Name:             name,
		Description:      description,
		LongDescription:  longDescription,
		ImageID:          imageID,
		SecondaryImageID: secondaryImageID,
		PdfID:            pdfID,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al crear item: %v", err))
	}

	// Redirigir a la lista de items (HTMX redirect)
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/admin/categories/%d/items", categoryID))
	return c.NoContent(http.StatusOK)
}

// HandleItemUpdate actualiza un item existente
func (h *Handler) HandleItemUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de item inválido")
	}

	// Verificar que el item existe y pertenece a la categoría
	item, err := h.queries.GetItem(ctx, int32(itemID))
	if err != nil {
		return c.String(http.StatusNotFound, "Item no encontrado")
	}

	if item.CategoryID != int32(categoryID) {
		return c.String(http.StatusBadRequest, "El item no pertenece a esta categoría")
	}

	// Parsear formulario
	slug := c.FormValue("slug")
	name := c.FormValue("name")
	description := c.FormValue("description")
	longDescription := c.FormValue("long_description")

	// Validación básica
	if name == "" || slug == "" {
		return c.String(http.StatusBadRequest, "Nombre y slug son requeridos")
	}

	// Parsear image_id opcional
	var imageID pgtype.Int4
	if imageIDStr := c.FormValue("image_id"); imageIDStr != "" {
		if id, err := strconv.ParseInt(imageIDStr, 10, 32); err == nil {
			imageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Parsear secondary_image_id opcional
	var secondaryImageID pgtype.Int4
	if secondaryImageIDStr := c.FormValue("secondary_image_id"); secondaryImageIDStr != "" {
		if id, err := strconv.ParseInt(secondaryImageIDStr, 10, 32); err == nil {
			secondaryImageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Parsear pdf_id opcional
	var pdfID pgtype.Int4
	if pdfIDStr := c.FormValue("pdf_id"); pdfIDStr != "" {
		if id, err := strconv.ParseInt(pdfIDStr, 10, 32); err == nil {
			pdfID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Actualizar item
	err = h.queries.UpdateItem(ctx, repository.UpdateItemParams{
		ItemID:           int32(itemID),
		Slug:             slug,
		Name:             name,
		Description:      description,
		LongDescription:  longDescription,
		ImageID:          imageID,
		SecondaryImageID: secondaryImageID,
		PdfID:            pdfID,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al actualizar item: %v", err))
	}

	// Redirigir a la lista de items (HTMX redirect)
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/admin/categories/%d/items", categoryID))
	return c.NoContent(http.StatusOK)
}

// HandleItemDelete elimina un item
func (h *Handler) HandleItemDelete(c echo.Context) error {
	ctx := c.Request().Context()

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de item inválido")
	}

	// Verificar que el item existe y pertenece a la categoría
	item, err := h.queries.GetItem(ctx, int32(itemID))
	if err != nil {
		return c.String(http.StatusNotFound, "Item no encontrado")
	}

	if item.CategoryID != int32(categoryID) {
		return c.String(http.StatusBadRequest, "El item no pertenece a esta categoría")
	}

	// Eliminar item (las imágenes de la galería se eliminan en cascada)
	if err := h.queries.DeleteItem(ctx, int32(itemID)); err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar item")
	}

	// Retornar vacío para que HTMX elimine la fila
	return c.NoContent(http.StatusOK)
}
