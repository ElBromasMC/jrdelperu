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

// HandleArticlesIndex muestra la lista de artículos del blog
func (h *Handler) HandleArticlesIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	articles, err := h.queries.ListAllArticles(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar artículos")
	}

	// Load cover images for articles
	imagesMap := make(map[int32]repository.StaticFile)
	for _, article := range articles {
		if article.CoverImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, article.CoverImageID.Int32)
			if err == nil {
				imagesMap[article.CoverImageID.Int32] = img
			}
		}
	}

	return Render(c, http.StatusOK, view.AdminArticlesPage(sessionData.Username, articles, imagesMap))
}

// HandleArticleNewForm muestra el formulario para crear un nuevo artículo
func (h *Handler) HandleArticleNewForm(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Get images for dropdown
	images, _ := h.queries.ListImages(ctx)

	return Render(c, http.StatusOK, view.AdminArticleFormPage(sessionData.Username, nil, images, false))
}

// HandleArticleCreate crea un nuevo artículo
func (h *Handler) HandleArticleCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse form data
	slug := c.FormValue("slug")
	title := c.FormValue("title")
	summary := c.FormValue("summary")
	content := c.FormValue("content")
	author := c.FormValue("author")
	coverImageIDStr := c.FormValue("cover_image_id")

	// Validate required fields
	if slug == "" || title == "" || summary == "" || content == "" {
		return c.String(http.StatusBadRequest, "Todos los campos requeridos deben ser completados")
	}

	// Default author
	if author == "" {
		author = "JR del Perú"
	}

	// Convert optional cover image ID
	var coverImageID pgtype.Int4
	if coverImageIDStr != "" {
		id, err := strconv.Atoi(coverImageIDStr)
		if err == nil {
			coverImageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Create article
	_, err := h.queries.CreateArticle(ctx, repository.CreateArticleParams{
		Slug:         slug,
		Title:        title,
		Summary:      summary,
		Content:      content,
		CoverImageID: coverImageID,
		Author:       author,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al crear artículo: %v", err))
	}

	// Redirect to articles list
	c.Response().Header().Set("HX-Redirect", "/admin/articles")
	return c.NoContent(http.StatusOK)
}

// HandleArticleEditForm muestra el formulario para editar un artículo
func (h *Handler) HandleArticleEditForm(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Get article
	article, err := h.queries.GetArticle(ctx, int32(articleID))
	if err != nil {
		return c.String(http.StatusNotFound, "Artículo no encontrado")
	}

	// Get images for dropdown
	images, _ := h.queries.ListImages(ctx)

	return Render(c, http.StatusOK, view.AdminArticleFormPage(sessionData.Username, &article, images, true))
}

// HandleArticleUpdate actualiza un artículo existente
func (h *Handler) HandleArticleUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Parse form data
	slug := c.FormValue("slug")
	title := c.FormValue("title")
	summary := c.FormValue("summary")
	content := c.FormValue("content")
	author := c.FormValue("author")
	coverImageIDStr := c.FormValue("cover_image_id")

	// Validate required fields
	if slug == "" || title == "" || summary == "" || content == "" {
		return c.String(http.StatusBadRequest, "Todos los campos requeridos deben ser completados")
	}

	// Default author
	if author == "" {
		author = "JR del Perú"
	}

	// Convert optional cover image ID
	var coverImageID pgtype.Int4
	if coverImageIDStr != "" {
		id, err := strconv.Atoi(coverImageIDStr)
		if err == nil {
			coverImageID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Update article
	err = h.queries.UpdateArticle(ctx, repository.UpdateArticleParams{
		ArticleID:    int32(articleID),
		Slug:         slug,
		Title:        title,
		Summary:      summary,
		Content:      content,
		CoverImageID: coverImageID,
		Author:       author,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al actualizar artículo: %v", err))
	}

	// Redirect to articles list
	c.Response().Header().Set("HX-Redirect", "/admin/articles")
	return c.NoContent(http.StatusOK)
}

// HandleArticleDelete elimina un artículo
func (h *Handler) HandleArticleDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Delete article (cascade will delete FAQs)
	err = h.queries.DeleteArticle(ctx, int32(articleID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al eliminar artículo: %v", err))
	}

	// Return empty response (HTMX will remove the row)
	return c.NoContent(http.StatusOK)
}

// HandleArticlePublish publica un artículo
func (h *Handler) HandleArticlePublish(c echo.Context) error {
	ctx := c.Request().Context()

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Publish article
	err = h.queries.PublishArticle(ctx, int32(articleID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al publicar artículo: %v", err))
	}

	// Redirect to refresh the page
	c.Response().Header().Set("HX-Redirect", "/admin/articles")
	return c.NoContent(http.StatusOK)
}

// HandleArticleUnpublish despublica un artículo
func (h *Handler) HandleArticleUnpublish(c echo.Context) error {
	ctx := c.Request().Context()

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Unpublish article
	err = h.queries.UnpublishArticle(ctx, int32(articleID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al despublicar artículo: %v", err))
	}

	// Redirect to refresh the page
	c.Response().Header().Set("HX-Redirect", "/admin/articles")
	return c.NoContent(http.StatusOK)
}

// HandleArticleDetail muestra el detalle del artículo con FAQs
func (h *Handler) HandleArticleDetail(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Get article
	article, err := h.queries.GetArticle(ctx, int32(articleID))
	if err != nil {
		return c.String(http.StatusNotFound, "Artículo no encontrado")
	}

	// Get article FAQs
	faqs, err := h.queries.ListArticleFAQs(ctx, int32(articleID))
	if err != nil {
		faqs = []repository.ArticleFaq{}
	}

	// Get cover image if exists
	var coverImage *repository.StaticFile
	if article.CoverImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, article.CoverImageID.Int32)
		if err == nil {
			coverImage = &img
		}
	}

	return Render(c, http.StatusOK, view.AdminArticleDetailPage(sessionData.Username, article, faqs, coverImage))
}

// HandleArticleFAQCreate crea una FAQ para un artículo
func (h *Handler) HandleArticleFAQCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Parse form data
	question := c.FormValue("question")
	answer := c.FormValue("answer")
	displayOrderStr := c.FormValue("display_order")

	// Validate required fields
	if question == "" || answer == "" {
		return c.String(http.StatusBadRequest, "Pregunta y respuesta son requeridas")
	}

	// Parse display order
	displayOrder := int32(0)
	if displayOrderStr != "" {
		if order, err := strconv.Atoi(displayOrderStr); err == nil {
			displayOrder = int32(order)
		}
	}

	// Create FAQ
	_, err = h.queries.CreateArticleFAQ(ctx, repository.CreateArticleFAQParams{
		ArticleID:    int32(articleID),
		Question:     question,
		Answer:       answer,
		DisplayOrder: displayOrder,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al crear FAQ: %v", err))
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListArticleFAQs(ctx, int32(articleID))
	return Render(c, http.StatusOK, view.AdminArticleFAQsGrid(int32(articleID), faqs))
}

// HandleArticleFAQUpdate actualiza una FAQ de artículo
func (h *Handler) HandleArticleFAQUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get IDs
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	faqIDStr := c.Param("faqId")
	faqID, err := strconv.Atoi(faqIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de FAQ inválido")
	}

	// Parse form data
	question := c.FormValue("question")
	answer := c.FormValue("answer")
	displayOrderStr := c.FormValue("display_order")

	// Validate required fields
	if question == "" || answer == "" {
		return c.String(http.StatusBadRequest, "Pregunta y respuesta son requeridas")
	}

	// Parse display order
	displayOrder := int32(0)
	if displayOrderStr != "" {
		if order, err := strconv.Atoi(displayOrderStr); err == nil {
			displayOrder = int32(order)
		}
	}

	// Update FAQ
	err = h.queries.UpdateArticleFAQ(ctx, repository.UpdateArticleFAQParams{
		FaqID:        int32(faqID),
		Question:     question,
		Answer:       answer,
		DisplayOrder: displayOrder,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al actualizar FAQ: %v", err))
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListArticleFAQs(ctx, int32(articleID))
	return Render(c, http.StatusOK, view.AdminArticleFAQsGrid(int32(articleID), faqs))
}

// HandleArticleFAQDelete elimina una FAQ de artículo
func (h *Handler) HandleArticleFAQDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Get IDs
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	faqIDStr := c.Param("faqId")
	faqID, err := strconv.Atoi(faqIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de FAQ inválido")
	}

	// Delete FAQ
	err = h.queries.DeleteArticleFAQ(ctx, int32(faqID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al eliminar FAQ: %v", err))
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListArticleFAQs(ctx, int32(articleID))
	return Render(c, http.StatusOK, view.AdminArticleFAQsGrid(int32(articleID), faqs))
}
