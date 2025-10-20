package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// HandleFilesIndex muestra la lista de archivos
func (h *Handler) HandleFilesIndex(c echo.Context) error {
	sessionData := c.Get("session").(*service.SessionData)

	// Obtener lista de imágenes y PDFs
	images, err := h.queries.ListImages(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes")
	}

	pdfs, err := h.queries.ListPDFs(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener PDFs")
	}

	return Render(c, http.StatusOK, view.AdminFilesIndex(sessionData.Username, images, pdfs))
}

// HandleFileUpload maneja la subida de archivos
func (h *Handler) HandleFileUpload(c echo.Context) error {
	ctx := c.Request().Context()

	// Obtener archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "No se proporcionó archivo")
	}

	// Obtener display name personalizado (opcional)
	customDisplayName := c.FormValue("display_name")

	// Subir archivo (tipo se infiere, display name puede ser personalizado o auto-generado)
	result, err := h.fileService.UploadFile(ctx, file, customDisplayName)
	if err != nil {
		if errors.Is(err, service.ErrInvalidFileType) {
			return c.String(http.StatusBadRequest, "Tipo de archivo no permitido. Solo imágenes (JPG, PNG, WebP, GIF) y PDFs.")
		}
		if errors.Is(err, service.ErrFileTooLarge) {
			return c.String(http.StatusBadRequest, "Archivo demasiado grande. Máximo: 10 MB para imágenes, 20 MB para PDFs.")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error al subir archivo: %v", err))
	}

	// Obtener archivo completo para renderizar
	uploadedFile, err := h.queries.GetStaticFile(ctx, result.FileID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener archivo subido")
	}

	// Obtener conteos actualizados
	imagesCount, err := h.queries.CountImages(ctx)
	if err != nil {
		imagesCount = 0
	}

	pdfsCount, err := h.queries.CountPDFs(ctx)
	if err != nil {
		pdfsCount = 0
	}

	// Retornar respuesta HTMX con datos del archivo (para auto-select en formularios)
	triggerData := fmt.Sprintf(`{"fileUploaded":{"fileId":%d,"fileType":"%s","displayName":"%s"}}`,
		result.FileID, result.FileType, result.DisplayName)
	c.Response().Header().Set("HX-Trigger", triggerData)

	// Renderizar respuesta con OOB swaps:
	// 1. Mensaje de éxito para el target original
	// 2. Nueva FileCard para el grid correspondiente
	// 3. Contadores actualizados
	return Render(c, http.StatusOK, view.FileUploadResponse(
		result.DisplayName,
		uploadedFile,
		result.FileType,
		int(imagesCount),
		int(pdfsCount),
	))
}

// HandleFileGet devuelve una tarjeta de archivo en modo visualización
func (h *Handler) HandleFileGet(c echo.Context) error {
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de archivo inválido")
	}

	// Obtener archivo
	ctx := c.Request().Context()
	file, err := h.queries.GetStaticFile(ctx, int32(fileID))
	if err != nil {
		return c.String(http.StatusNotFound, "Archivo no encontrado")
	}

	// Determinar tipo de archivo
	fileType := "image"
	if file.FileType.Valid && file.FileType.String == "pdf" {
		fileType = "pdf"
	}

	return Render(c, http.StatusOK, view.FileCard(file, fileType))
}

// HandleFileEdit devuelve una tarjeta de archivo en modo edición
func (h *Handler) HandleFileEdit(c echo.Context) error {
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de archivo inválido")
	}

	// Obtener archivo
	ctx := c.Request().Context()
	file, err := h.queries.GetStaticFile(ctx, int32(fileID))
	if err != nil {
		return c.String(http.StatusNotFound, "Archivo no encontrado")
	}

	// Determinar tipo de archivo
	fileType := "image"
	if file.FileType.Valid && file.FileType.String == "pdf" {
		fileType = "pdf"
	}

	return Render(c, http.StatusOK, view.FileCardEdit(file, fileType))
}

// HandleFileUpdateDisplayName actualiza el nombre de visualización de un archivo
func (h *Handler) HandleFileUpdateDisplayName(c echo.Context) error {
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de archivo inválido")
	}

	displayName := c.FormValue("display_name")
	if displayName == "" {
		return c.String(http.StatusBadRequest, "Nombre de visualización requerido")
	}

	// Actualizar display name
	ctx := c.Request().Context()
	err = h.queries.UpdateStaticFileDisplayName(ctx, repository.UpdateStaticFileDisplayNameParams{
		FileID:      int32(fileID),
		DisplayName: pgtype.Text{String: displayName, Valid: true},
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al actualizar nombre")
	}

	// Obtener archivo actualizado para devolver la tarjeta
	file, err := h.queries.GetStaticFile(ctx, int32(fileID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener archivo")
	}

	// Determinar tipo de archivo para la tarjeta
	fileType := "image"
	if file.FileType.Valid && file.FileType.String == "pdf" {
		fileType = "pdf"
	}

	return Render(c, http.StatusOK, view.FileCard(file, fileType))
}

// HandleFileDelete elimina un archivo
func (h *Handler) HandleFileDelete(c echo.Context) error {
	ctx := c.Request().Context()
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de archivo inválido")
	}

	// Eliminar archivo
	if err := h.fileService.DeleteFile(ctx, int32(fileID)); err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			return c.String(http.StatusNotFound, "Archivo no encontrado")
		}
		return c.String(http.StatusInternalServerError, "Error al eliminar archivo")
	}

	// Obtener conteos actualizados
	imagesCount, err := h.queries.CountImages(ctx)
	if err != nil {
		imagesCount = 0
	}

	pdfsCount, err := h.queries.CountPDFs(ctx)
	if err != nil {
		pdfsCount = 0
	}

	// Retornar respuesta HTMX con OOB swaps para actualizar contadores
	c.Response().Header().Set("HX-Trigger", "fileDeleted")
	return Render(c, http.StatusOK, view.FileDeleteResponse(int(imagesCount), int(pdfsCount)))
}
