package handler

import (
	"alc/config"
	"alc/model"
	"alc/repository"
	"alc/service"
	"alc/view"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleStoreIndexShow(c echo.Context) error {
	return renderOK(c, view.StoreIndex())
}

func (h *Handler) HandleVidrioIndexShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get all vidrio categories from database
	dbCats, err := h.queries.ListCategoriesByMaterialType(ctx, repository.MaterialTypeVidrio)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar categorías")
	}

	// Convert to model.Category
	cats := make([]model.Category, 0, len(dbCats))
	for _, dbCat := range dbCats {
		// Get image if exists
		var dbImage *repository.StaticFile
		if dbCat.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
			if err == nil {
				dbImage = &img
			}
		}

		cats = append(cats, service.MapCategoryToModel(dbCat, dbImage))
	}

	// Get catalog PDF from site documents
	pdfURL := ""
	pdfName := ""
	doc, err := h.queries.GetSiteDocumentByKey(ctx, "catalogo_vidrios")
	if err == nil && doc.FileID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = doc.DisplayName
		}
	}

	return renderOK(c, view.StoreVidrioIndex(cats, pdfURL, pdfName))
}

func (h *Handler) HandleVidrioCategoryShow(c echo.Context) error {
	ctx := c.Request().Context()
	catSlug := c.Param("catSlug")

	// Get category by slug
	dbCat, err := h.queries.GetCategoryBySlug(ctx, repository.GetCategoryBySlugParams{
		MaterialType: repository.MaterialTypeVidrio,
		Slug:         catSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Get category image
	var catImage *repository.StaticFile
	if dbCat.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
		if err == nil {
			catImage = &img
		}
	}

	cat := service.MapCategoryToModel(dbCat, catImage)

	// Get items for this category
	dbItems, err := h.queries.ListItemsByCategory(ctx, dbCat.CategoryID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar items")
	}

	items := make([]model.Item, 0, len(dbItems))
	for _, dbItem := range dbItems {
		var itemImage *repository.StaticFile
		if dbItem.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, dbItem.ImageID.Int32)
			if err == nil {
				itemImage = &img
			}
		}
		items = append(items, service.MapItemToModel(dbItem, dbCat, itemImage, nil, catImage))
	}

	// Get features for this category
	dbFeatures, err := h.queries.ListCategoryFeatures(ctx, dbCat.CategoryID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar características")
	}

	features := make([]model.CategoryFeature, 0, len(dbFeatures))
	for _, dbFeature := range dbFeatures {
		features = append(features, service.MapCategoryFeatureToModel(dbFeature, dbCat, catImage))
	}

	return renderOK(c, view.StoreVidrioCategory(cat, items, features))
}

func (h *Handler) HandleAluminioIndexShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get all aluminio categories from database
	dbCats, err := h.queries.ListCategoriesByMaterialType(ctx, repository.MaterialTypeAluminio)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar categorías")
	}

	// Convert to model.Category
	cats := make([]model.Category, 0, len(dbCats))
	for _, dbCat := range dbCats {
		var dbImage *repository.StaticFile
		if dbCat.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
			if err == nil {
				dbImage = &img
			}
		}
		cats = append(cats, service.MapCategoryToModel(dbCat, dbImage))
	}

	// Get catalog PDF from site documents
	pdfURL := ""
	pdfName := ""
	doc, err := h.queries.GetSiteDocumentByKey(ctx, "catalogo_aluminios")
	if err == nil && doc.FileID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = doc.DisplayName
		}
	}

	return renderOK(c, view.StoreAluminioIndex(cats, pdfURL, pdfName))
}

func (h *Handler) HandleAluminioCategoryShow(c echo.Context) error {
	ctx := c.Request().Context()
	catSlug := c.Param("catSlug")

	// Get category by slug
	dbCat, err := h.queries.GetCategoryBySlug(ctx, repository.GetCategoryBySlugParams{
		MaterialType: repository.MaterialTypeAluminio,
		Slug:         catSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Get category image
	var catImage *repository.StaticFile
	if dbCat.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
		if err == nil {
			catImage = &img
		}
	}

	cat := service.MapCategoryToModel(dbCat, catImage)

	// Get items for this category
	dbItems, err := h.queries.ListItemsByCategory(ctx, dbCat.CategoryID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar items")
	}

	items := make([]model.Item, 0, len(dbItems))
	for _, dbItem := range dbItems {
		var itemImage *repository.StaticFile
		if dbItem.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, dbItem.ImageID.Int32)
			if err == nil {
				itemImage = &img
			}
		}
		items = append(items, service.MapItemToModel(dbItem, dbCat, itemImage, nil, catImage))
	}

	// Get PDF URL if exists
	pdfURL := ""
	pdfName := ""
	if dbCat.PdfID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, dbCat.PdfID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = pdf.DisplayName.String
		}
	}

	return renderOK(c, view.StoreAluminioCategory(cat, items, pdfURL, pdfName))
}

func (h *Handler) HandleAluminioItemShow(c echo.Context) error {
	ctx := c.Request().Context()
	catSlug := c.Param("catSlug")
	itemSlug := c.Param("itemSlug")

	// Get category by slug
	dbCat, err := h.queries.GetCategoryBySlug(ctx, repository.GetCategoryBySlugParams{
		MaterialType: repository.MaterialTypeAluminio,
		Slug:         catSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Get item by slug
	dbItem, err := h.queries.GetItemBySlug(ctx, repository.GetItemBySlugParams{
		CategoryID: dbCat.CategoryID,
		Slug:       itemSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Item no encontrado")
	}

	// Get category image
	var catImage *repository.StaticFile
	if dbCat.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
		if err == nil {
			catImage = &img
		}
	}

	// Get item image
	var itemImage *repository.StaticFile
	if dbItem.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbItem.ImageID.Int32)
		if err == nil {
			itemImage = &img
		}
	}

	// Get secondary item image
	var secondaryImage *repository.StaticFile
	if dbItem.SecondaryImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbItem.SecondaryImageID.Int32)
		if err == nil {
			secondaryImage = &img
		}
	}

	item := service.MapItemToModel(dbItem, dbCat, itemImage, secondaryImage, catImage)

	// Get item PDF URL if exists
	pdfURL := ""
	pdfName := ""
	if dbItem.PdfID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, dbItem.PdfID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = pdf.DisplayName.String
		}
	}

	return renderOK(c, view.StoreAluminioItem(item, pdfURL, pdfName))
}

func (h *Handler) HandleUPVCIndexShow(c echo.Context) error {
	ctx := c.Request().Context()

	// Get all uPVC categories from database
	dbCats, err := h.queries.ListCategoriesByMaterialType(ctx, repository.MaterialTypeUpvc)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar categorías")
	}

	// Convert to model.Category
	cats := make([]model.Category, 0, len(dbCats))
	for _, dbCat := range dbCats {
		var dbImage *repository.StaticFile
		if dbCat.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
			if err == nil {
				dbImage = &img
			}
		}
		cats = append(cats, service.MapCategoryToModel(dbCat, dbImage))
	}

	// Get catalog PDF from site documents
	pdfURL := ""
	pdfName := ""
	doc, err := h.queries.GetSiteDocumentByKey(ctx, "catalogo_upvc")
	if err == nil && doc.FileID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, doc.FileID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = doc.DisplayName
		}
	}

	// Get afiche PDF from site documents
	aficheURL := ""
	aficheName := ""
	aficheDoc, err := h.queries.GetSiteDocumentByKey(ctx, "afiche_upvc")
	if err == nil && aficheDoc.FileID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, aficheDoc.FileID.Int32)
		if err == nil {
			aficheURL = path.Join(config.PDFS_PATH, pdf.FileName)
			aficheName = aficheDoc.DisplayName
		}
	}

	return renderOK(c, view.StoreUPVCIndex(cats, pdfURL, pdfName, aficheURL, aficheName))
}

func (h *Handler) HandleUPVCCategoryShow(c echo.Context) error {
	ctx := c.Request().Context()
	catSlug := c.Param("catSlug")

	// Get category by slug
	dbCat, err := h.queries.GetCategoryBySlug(ctx, repository.GetCategoryBySlugParams{
		MaterialType: repository.MaterialTypeUpvc,
		Slug:         catSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Get category image
	var catImage *repository.StaticFile
	if dbCat.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
		if err == nil {
			catImage = &img
		}
	}

	cat := service.MapCategoryToModel(dbCat, catImage)

	// Get items for this category
	dbItems, err := h.queries.ListItemsByCategory(ctx, dbCat.CategoryID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar items")
	}

	items := make([]model.Item, 0, len(dbItems))
	for _, dbItem := range dbItems {
		var itemImage *repository.StaticFile
		if dbItem.ImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, dbItem.ImageID.Int32)
			if err == nil {
				itemImage = &img
			}
		}
		items = append(items, service.MapItemToModel(dbItem, dbCat, itemImage, nil, catImage))
	}

	// Get PDF URL if exists
	pdfURL := ""
	pdfName := ""
	if dbCat.PdfID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, dbCat.PdfID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = pdf.DisplayName.String
		}
	}

	return renderOK(c, view.StoreUPVCCategory(cat, items, pdfURL, pdfName))
}

func (h *Handler) HandleUPVCItemShow(c echo.Context) error {
	ctx := c.Request().Context()
	catSlug := c.Param("catSlug")
	itemSlug := c.Param("itemSlug")

	// Get category by slug
	dbCat, err := h.queries.GetCategoryBySlug(ctx, repository.GetCategoryBySlugParams{
		MaterialType: repository.MaterialTypeUpvc,
		Slug:         catSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Categoría no encontrada")
	}

	// Get item by slug
	dbItem, err := h.queries.GetItemBySlug(ctx, repository.GetItemBySlugParams{
		CategoryID: dbCat.CategoryID,
		Slug:       itemSlug,
	})
	if err != nil {
		return c.String(http.StatusNotFound, "Item no encontrado")
	}

	// Get category image
	var catImage *repository.StaticFile
	if dbCat.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbCat.ImageID.Int32)
		if err == nil {
			catImage = &img
		}
	}

	// Get item image
	var itemImage *repository.StaticFile
	if dbItem.ImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbItem.ImageID.Int32)
		if err == nil {
			itemImage = &img
		}
	}

	// Get secondary item image
	var secondaryImage *repository.StaticFile
	if dbItem.SecondaryImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, dbItem.SecondaryImageID.Int32)
		if err == nil {
			secondaryImage = &img
		}
	}

	item := service.MapItemToModel(dbItem, dbCat, itemImage, secondaryImage, catImage)

	// Get item PDF URL if exists
	pdfURL := ""
	pdfName := ""
	if dbItem.PdfID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, dbItem.PdfID.Int32)
		if err == nil {
			pdfURL = path.Join(config.PDFS_PATH, pdf.FileName)
			pdfName = pdf.DisplayName.String
		}
	}

	return renderOK(c, view.StoreUPVCItem(item, pdfURL, pdfName))
}
