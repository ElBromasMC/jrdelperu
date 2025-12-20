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

// HandleCategoriesIndex muestra la lista de categorías
func (h *Handler) HandleCategoriesIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Get filter by material type
	materialTypeFilter := c.QueryParam("type")

	var categories []repository.Category
	var err error

	// Filter by material type if specified
	if materialTypeFilter != "" {
		var materialType repository.MaterialType
		switch materialTypeFilter {
		case "vidrio":
			materialType = repository.MaterialTypeVidrio
		case "aluminio":
			materialType = repository.MaterialTypeAluminio
		case "upvc":
			materialType = repository.MaterialTypeUpvc
		default:
			return c.String(http.StatusBadRequest, "Tipo de material inválido")
		}
		categories, err = h.queries.ListCategoriesByMaterialType(ctx, materialType)
	} else {
		categories, err = h.queries.ListAllCategories(ctx)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar categorías")
	}

	// Load images for categories
	imagesMap := make(map[int32]repository.StaticFile)
	for _, cat := range categories {
		if cat.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, cat.ImageID.Int32)
			if err == nil {
				imagesMap[cat.ImageID.Int32] = img
			}
		}
	}

	return Render(c, http.StatusOK, view.AdminCategoriesPage(sessionData.Username, categories, imagesMap))
}

// HandleCategoryNewForm muestra el formulario para crear una nueva categoría
func (h *Handler) HandleCategoryNewForm(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Get images and tags for dropdowns
	images, _ := h.queries.ListImages(ctx)
	tags, _ := h.queries.ListCategoryTags(ctx)
	pdfs, _ := h.queries.ListPDFs(ctx)

	return Render(c, http.StatusOK, view.AdminCategoryFormPage(sessionData.Username, nil, images, tags, pdfs, false))
}

// HandleCategoryCreate crea una nueva categoría
func (h *Handler) HandleCategoryCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse form data
	materialTypeStr := c.FormValue("material_type")
	slug := c.FormValue("slug")
	name := c.FormValue("name")
	description := c.FormValue("description")
	longDescription := c.FormValue("long_description")
	imageIDStr := c.FormValue("image_id")
	secondaryImageIDStr := c.FormValue("secondary_image_id")
	tagIDStr := c.FormValue("tag_id")
	pdfIDStr := c.FormValue("pdf_id")

	// Convert material type
	var materialType repository.MaterialType
	switch materialTypeStr {
	case "vidrio":
		materialType = repository.MaterialTypeVidrio
	case "aluminio":
		materialType = repository.MaterialTypeAluminio
	case "upvc":
		materialType = repository.MaterialTypeUpvc
	default:
		return c.String(http.StatusBadRequest, "Tipo de material inválido")
	}

	// Convert optional IDs
	var imageID pgtype.Int4
	if imageIDStr != "" {
		id, err := strconv.Atoi(imageIDStr)
		if err == nil {
			imageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	var secondaryImageID pgtype.Int4
	if secondaryImageIDStr != "" {
		id, err := strconv.Atoi(secondaryImageIDStr)
		if err == nil {
			secondaryImageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	var tagID pgtype.Int4
	if tagIDStr != "" {
		id, err := strconv.Atoi(tagIDStr)
		if err == nil {
			tagID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	var pdfID pgtype.Int4
	if pdfIDStr != "" {
		id, err := strconv.Atoi(pdfIDStr)
		if err == nil {
			pdfID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Create category
	_, err := h.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		MaterialType:      materialType,
		Slug:              slug,
		Name:              name,
		Description:       description,
		LongDescription:   longDescription,
		ImageID:           imageID,
		SecondaryImageID:  secondaryImageID,
		TagID:             tagID,
		PdfID:             pdfID,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al crear categoría: %v", err))
	}

	// Redirect to categories list
	c.Response().Header().Set("HX-Redirect", "/admin/categories")
	return c.NoContent(http.StatusOK)
}

// HandleCategoryEditForm muestra el formulario para editar una categoría
func (h *Handler) HandleCategoryEditForm(c echo.Context) error {
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

	// Get images, tags, and pdfs for dropdowns
	images, _ := h.queries.ListImages(ctx)
	tags, _ := h.queries.ListCategoryTags(ctx)
	pdfs, _ := h.queries.ListPDFs(ctx)

	return Render(c, http.StatusOK, view.AdminCategoryFormPage(sessionData.Username, &category, images, tags, pdfs, true))
}

// HandleCategoryUpdate actualiza una categoría existente
func (h *Handler) HandleCategoryUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Parse form data
	slug := c.FormValue("slug")
	name := c.FormValue("name")
	description := c.FormValue("description")
	longDescription := c.FormValue("long_description")
	imageIDStr := c.FormValue("image_id")
	secondaryImageIDStr := c.FormValue("secondary_image_id")
	tagIDStr := c.FormValue("tag_id")
	pdfIDStr := c.FormValue("pdf_id")

	// Convert optional IDs
	var imageID pgtype.Int4
	if imageIDStr != "" {
		id, err := strconv.Atoi(imageIDStr)
		if err == nil {
			imageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	var secondaryImageID pgtype.Int4
	if secondaryImageIDStr != "" {
		id, err := strconv.Atoi(secondaryImageIDStr)
		if err == nil {
			secondaryImageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	var tagID pgtype.Int4
	if tagIDStr != "" {
		id, err := strconv.Atoi(tagIDStr)
		if err == nil {
			tagID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	var pdfID pgtype.Int4
	if pdfIDStr != "" {
		id, err := strconv.Atoi(pdfIDStr)
		if err == nil {
			pdfID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Update category
	err = h.queries.UpdateCategory(ctx, repository.UpdateCategoryParams{
		CategoryID:       int32(categoryID),
		Slug:             slug,
		Name:             name,
		Description:      description,
		LongDescription:  longDescription,
		ImageID:          imageID,
		SecondaryImageID: secondaryImageID,
		TagID:            tagID,
		PdfID:            pdfID,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al actualizar categoría: %v", err))
	}

	// Redirect to categories list
	c.Response().Header().Set("HX-Redirect", "/admin/categories")
	return c.NoContent(http.StatusOK)
}

// HandleCategoryDelete elimina una categoría
func (h *Handler) HandleCategoryDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Get category ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de categoría inválido")
	}

	// Delete category (cascade will delete items and features)
	err = h.queries.DeleteCategory(ctx, int32(categoryID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al eliminar categoría: %v", err))
	}

	// Return empty response (HTMX will remove the row)
	return c.NoContent(http.StatusOK)
}
