package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// HandleDocumentsIndex muestra la lista de documentos del sitio
func (h *Handler) HandleDocumentsIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Obtener todos los documentos
	documents, err := h.queries.ListSiteDocuments(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener documentos")
	}

	// Obtener PDFs disponibles
	pdfs, err := h.queries.ListPDFs(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener PDFs")
	}

	// Crear mapa de archivos para mostrar info del PDF asignado
	filesMap := make(map[int32]repository.StaticFile)
	for _, doc := range documents {
		if doc.FileID.Valid {
			file, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
			if err == nil {
				filesMap[doc.FileID.Int32] = file
			}
		}
	}

	return Render(c, http.StatusOK, view.AdminDocumentsPage(sessionData.Username, documents, pdfs, filesMap))
}

// HandleDocumentUpdate actualiza un documento del sitio
func (h *Handler) HandleDocumentUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseInt(documentIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de documento inválido")
	}

	// Parsear formulario
	displayName := c.FormValue("display_name")
	if displayName == "" {
		return c.String(http.StatusBadRequest, "Nombre es requerido")
	}

	// Parsear file_id opcional
	var fileID pgtype.Int4
	if fileIDStr := c.FormValue("file_id"); fileIDStr != "" {
		if id, err := strconv.ParseInt(fileIDStr, 10, 32); err == nil {
			fileID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Actualizar documento
	err = h.queries.UpdateSiteDocument(ctx, repository.UpdateSiteDocumentParams{
		DocumentID:  int32(documentID),
		DisplayName: displayName,
		FileID:      fileID,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al actualizar documento")
	}

	// Obtener documento actualizado
	documents, err := h.queries.ListSiteDocuments(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener documentos")
	}

	// Obtener PDFs disponibles
	pdfs, err := h.queries.ListPDFs(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener PDFs")
	}

	// Crear mapa de archivos
	filesMap := make(map[int32]repository.StaticFile)
	for _, doc := range documents {
		if doc.FileID.Valid {
			file, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
			if err == nil {
				filesMap[doc.FileID.Int32] = file
			}
		}
	}

	// Retornar la tabla actualizada
	return Render(c, http.StatusOK, view.AdminDocumentsTable(documents, pdfs, filesMap))
}
