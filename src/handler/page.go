package handler

import (
	"alc/config"
	"alc/model"
	"alc/repository"
	"alc/view"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get visible projects
	projects, err := h.queries.ListVisibleProjects(ctx)
	if err != nil {
		projects = []repository.Project{}
	}

	// Get featured images for each project
	projectImages := make(map[int32]repository.ListProjectImagesRow)
	for _, proj := range projects {
		featuredImg, err := h.queries.GetFeaturedProjectImage(ctx, proj.ProjectID)
		if err == nil {
			// Convert GetFeaturedProjectImageRow to ListProjectImagesRow
			projectImages[proj.ProjectID] = repository.ListProjectImagesRow{
				ProjectID:    featuredImg.ProjectID,
				ImageID:      featuredImg.ImageID,
				DisplayOrder: featuredImg.DisplayOrder,
				IsFeatured:   featuredImg.IsFeatured,
				FileName:     featuredImg.FileName,
				FileType:     featuredImg.FileType,
				MimeType:     featuredImg.MimeType,
				DisplayName:  featuredImg.DisplayName,
			}
		}
	}

	return renderOK(c, view.Index(projects, projectImages))
}

func (h *Handler) HandleNosotrosShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get visible projects (limit to 3 for homepage)
	projects, err := h.queries.ListVisibleProjects(ctx)
	if err != nil {
		projects = []repository.Project{}
	}

	// Get featured images for each project
	projectImages := make(map[int32]repository.ListProjectImagesRow)
	for _, proj := range projects {
		featuredImg, err := h.queries.GetFeaturedProjectImage(ctx, proj.ProjectID)
		if err == nil {
			// Convert GetFeaturedProjectImageRow to ListProjectImagesRow
			projectImages[proj.ProjectID] = repository.ListProjectImagesRow{
				ProjectID:    featuredImg.ProjectID,
				ImageID:      featuredImg.ImageID,
				DisplayOrder: featuredImg.DisplayOrder,
				IsFeatured:   featuredImg.IsFeatured,
				FileName:     featuredImg.FileName,
				FileType:     featuredImg.FileType,
				MimeType:     featuredImg.MimeType,
				DisplayName:  featuredImg.DisplayName,
			}
		}
	}

	return renderOK(c, view.Nosotros(projects, projectImages))
}

func (h *Handler) HandleDescargasShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get all site documents
	docs, err := h.queries.ListSiteDocuments(ctx)
	if err != nil {
		docs = []repository.SiteDocument{}
	}

	// Build document info with PDF URLs
	docInfos := make([]model.SiteDocumentInfo, 0, len(docs))
	for _, doc := range docs {
		info := model.SiteDocumentInfo{
			Key:         doc.DocumentKey,
			DisplayName: doc.DisplayName,
			HasFile:     false,
		}
		if doc.FileID.Valid {
			pdf, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
			if err == nil {
				info.URL = path.Join(config.PDFS_PATH, pdf.FileName)
				info.HasFile = true
			}
		}
		docInfos = append(docInfos, info)
	}

	// Helper to build PDF info from categories
	buildCategoryPDFs := func(materialType repository.MaterialType) []model.PDFInfo {
		cats, err := h.queries.ListCategoriesWithPDFByMaterialType(ctx, materialType)
		if err != nil {
			return nil
		}
		pdfs := make([]model.PDFInfo, 0, len(cats))
		for _, cat := range cats {
			displayName := cat.Name
			if cat.PdfDisplayName.Valid && cat.PdfDisplayName.String != "" {
				displayName = cat.PdfDisplayName.String
			}
			pdfs = append(pdfs, model.PDFInfo{
				Name:        cat.Name,
				DisplayName: displayName,
				URL:         path.Join(config.PDFS_PATH, cat.PdfFileName),
			})
		}
		return pdfs
	}

	// Helper to build PDF info from items
	buildItemPDFs := func(materialType repository.MaterialType) []model.PDFInfo {
		items, err := h.queries.ListItemsWithPDFByMaterialType(ctx, materialType)
		if err != nil {
			return nil
		}
		pdfs := make([]model.PDFInfo, 0, len(items))
		for _, item := range items {
			displayName := item.Name
			if item.PdfDisplayName.Valid && item.PdfDisplayName.String != "" {
				displayName = item.PdfDisplayName.String
			}
			pdfs = append(pdfs, model.PDFInfo{
				Name:        item.CategoryName + " - " + item.Name,
				DisplayName: displayName,
				URL:         path.Join(config.PDFS_PATH, item.PdfFileName),
			})
		}
		return pdfs
	}

	// Get PDFs from categories and items by material type
	aluminioCatPDFs := buildCategoryPDFs(repository.MaterialTypeAluminio)
	aluminioItemPDFs := buildItemPDFs(repository.MaterialTypeAluminio)
	upvcCatPDFs := buildCategoryPDFs(repository.MaterialTypeUpvc)
	upvcItemPDFs := buildItemPDFs(repository.MaterialTypeUpvc)
	vidrioCatPDFs := buildCategoryPDFs(repository.MaterialTypeVidrio)
	vidrioItemPDFs := buildItemPDFs(repository.MaterialTypeVidrio)

	return renderOK(c, view.Descargas(
		docInfos,
		aluminioCatPDFs, aluminioItemPDFs,
		upvcCatPDFs, upvcItemPDFs,
		vidrioCatPDFs, vidrioItemPDFs,
	))
}

func (h *Handler) HandleGaleriaShow(c echo.Context) error {
	return renderOK(c, view.Galeria())
}

func (h *Handler) HandleContactoShow(c echo.Context) error {
	return renderOK(c, view.Contacto())
}

// HandleProjectsShow muestra la página de todos los proyectos
func (h *Handler) HandleProjectsShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get all visible projects
	projects, err := h.queries.ListVisibleProjects(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener proyectos")
	}

	// Get all images for each project
	projectImages := make(map[int32][]repository.ListProjectImagesRow)
	for _, proj := range projects {
		images, err := h.queries.ListProjectImages(ctx, proj.ProjectID)
		if err == nil {
			projectImages[proj.ProjectID] = images
		}
	}

	// Get brochure PDF from site documents
	brochureURL := ""
	brochureName := ""
	doc, err := h.queries.GetSiteDocumentByKey(ctx, "brochure_empresa")
	if err == nil && doc.FileID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
		if err == nil {
			brochureURL = path.Join(config.PDFS_PATH, pdf.FileName)
			brochureName = doc.DisplayName
		}
	}

	return renderOK(c, view.ProjectsPage(projects, projectImages, brochureURL, brochureName))
}

// HandleProjectDetailShow muestra el detalle de un proyecto
func (h *Handler) HandleProjectDetailShow(c echo.Context) error {
	ctx := c.Request().Context()
	slug := c.Param("slug")

	// Get project by slug
	project, err := h.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		return c.String(http.StatusNotFound, "Proyecto no encontrado")
	}

	// Check if project is visible
	if !project.IsVisible.Valid || !project.IsVisible.Bool {
		return c.String(http.StatusNotFound, "Proyecto no encontrado")
	}

	// Get all images for the project
	images, err := h.queries.ListProjectImages(ctx, project.ProjectID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener imágenes del proyecto")
	}

	return renderOK(c, view.ProjectDetailPage(project, images))
}
