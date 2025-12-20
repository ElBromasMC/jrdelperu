package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// getGlobalFAQCategories converts pgtype.Text categories to []string
func getGlobalFAQCategories(h *Handler, ctx context.Context) []string {
	pgCategories, _ := h.queries.ListGlobalFAQCategories(ctx)
	categories := make([]string, 0, len(pgCategories))
	for _, cat := range pgCategories {
		if cat.Valid {
			categories = append(categories, cat.String)
		}
	}
	return categories
}

// HandleGlobalFAQsIndex muestra la lista de preguntas frecuentes globales
func (h *Handler) HandleGlobalFAQsIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	faqs, err := h.queries.ListAllGlobalFAQs(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar preguntas frecuentes")
	}

	// Get categories for filter and convert to []string
	categories := getGlobalFAQCategories(h, ctx)

	return Render(c, http.StatusOK, view.AdminGlobalFAQsPage(sessionData.Username, faqs, categories))
}

// HandleGlobalFAQCreate crea una nueva pregunta frecuente global
func (h *Handler) HandleGlobalFAQCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse form data
	category := c.FormValue("category")
	question := c.FormValue("question")
	answer := c.FormValue("answer")
	displayOrderStr := c.FormValue("display_order")
	isVisible := c.FormValue("is_visible") == "on"

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
	_, err := h.queries.CreateGlobalFAQ(ctx, repository.CreateGlobalFAQParams{
		Category: pgtype.Text{String: category, Valid: category != ""},
		Question: question,
		Answer:   answer,
		DisplayOrder: displayOrder,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al crear FAQ: %v", err))
	}

	// If not visible, update visibility
	if !isVisible {
		// Get the last created FAQ and set visibility
		faqs, _ := h.queries.ListAllGlobalFAQs(ctx)
		if len(faqs) > 0 {
			// Since we just created, the last one should be ours (sorted by display_order, faq_id)
			// Actually, we need to find the one we just created
			for _, faq := range faqs {
				if faq.Question == question && faq.Answer == answer {
					h.queries.SetGlobalFAQVisibility(ctx, repository.SetGlobalFAQVisibilityParams{
						FaqID:     faq.FaqID,
						IsVisible: false,
					})
					break
				}
			}
		}
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListAllGlobalFAQs(ctx)
	categories := getGlobalFAQCategories(h, ctx)
	return Render(c, http.StatusOK, view.AdminGlobalFAQsGrid(faqs, categories))
}

// HandleGlobalFAQUpdate actualiza una pregunta frecuente global
func (h *Handler) HandleGlobalFAQUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get FAQ ID
	faqIDStr := c.Param("id")
	faqID, err := strconv.Atoi(faqIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de FAQ inválido")
	}

	// Parse form data
	category := c.FormValue("category")
	question := c.FormValue("question")
	answer := c.FormValue("answer")
	displayOrderStr := c.FormValue("display_order")
	isVisible := c.FormValue("is_visible") == "on"

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
	err = h.queries.UpdateGlobalFAQ(ctx, repository.UpdateGlobalFAQParams{
		FaqID:        int32(faqID),
		Category:     pgtype.Text{String: category, Valid: category != ""},
		Question:     question,
		Answer:       answer,
		DisplayOrder: displayOrder,
		IsVisible:    isVisible,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al actualizar FAQ: %v", err))
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListAllGlobalFAQs(ctx)
	categories := getGlobalFAQCategories(h, ctx)
	return Render(c, http.StatusOK, view.AdminGlobalFAQsGrid(faqs, categories))
}

// HandleGlobalFAQDelete elimina una pregunta frecuente global
func (h *Handler) HandleGlobalFAQDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Get FAQ ID
	faqIDStr := c.Param("id")
	faqID, err := strconv.Atoi(faqIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de FAQ inválido")
	}

	// Delete FAQ
	err = h.queries.DeleteGlobalFAQ(ctx, int32(faqID))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al eliminar FAQ: %v", err))
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListAllGlobalFAQs(ctx)
	categories := getGlobalFAQCategories(h, ctx)
	return Render(c, http.StatusOK, view.AdminGlobalFAQsGrid(faqs, categories))
}

// HandleGlobalFAQToggleVisibility cambia la visibilidad de una FAQ
func (h *Handler) HandleGlobalFAQToggleVisibility(c echo.Context) error {
	ctx := c.Request().Context()

	// Get FAQ ID
	faqIDStr := c.Param("id")
	faqID, err := strconv.Atoi(faqIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de FAQ inválido")
	}

	// Get current FAQ
	faq, err := h.queries.GetGlobalFAQ(ctx, int32(faqID))
	if err != nil {
		return c.String(http.StatusNotFound, "FAQ no encontrada")
	}

	// Toggle visibility
	err = h.queries.SetGlobalFAQVisibility(ctx, repository.SetGlobalFAQVisibilityParams{
		FaqID:     int32(faqID),
		IsVisible: !faq.IsVisible,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al cambiar visibilidad: %v", err))
	}

	// Return updated FAQ list
	faqs, _ := h.queries.ListAllGlobalFAQs(ctx)
	categories := getGlobalFAQCategories(h, ctx)
	return Render(c, http.StatusOK, view.AdminGlobalFAQsGrid(faqs, categories))
}
