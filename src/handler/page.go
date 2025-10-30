package handler

import (
	"alc/repository"
	"alc/view"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get visible projects (limit to 3 for homepage)
	allProjects, err := h.queries.ListVisibleProjects(ctx)
	if err != nil {
		allProjects = []repository.Project{}
	}

	// Limit to 3 projects for homepage
	projects := allProjects
	if len(projects) > 3 {
		projects = projects[:3]
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
	return renderOK(c, view.Nosotros())
}

func (h *Handler) HandleDescargasShow(c echo.Context) error {
	return renderOK(c, view.Descargas())
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

	return renderOK(c, view.ProjectsPage(projects, projectImages))
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
