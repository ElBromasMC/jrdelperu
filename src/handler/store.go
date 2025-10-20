package handler

import (
	"alc/model"
	"alc/repository"
	"alc/service"
	"alc/view"
	"net/http"

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

	return renderOK(c, view.StoreVidrioIndex(cats))
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
		items = append(items, service.MapItemToModel(dbItem, dbCat, itemImage, catImage))
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

	return renderOK(c, view.StoreAluminioIndex(cats))
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
		items = append(items, service.MapItemToModel(dbItem, dbCat, itemImage, catImage))
	}

	// Get PDF URL if exists
	pdfURL := ""
	if dbCat.PdfID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, dbCat.PdfID.Int32)
		if err == nil {
			pdfURL = "/uploads/pdfs/" + pdf.FileName
		}
	}

	return renderOK(c, view.StoreAluminioCategory(cat, items, pdfURL))
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

	item := service.MapItemToModel(dbItem, dbCat, itemImage, catImage)

	// Get item images (gallery)
	dbItemImages, err := h.queries.ListItemImages(ctx, dbItem.ItemID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar galería")
	}

	images := make([]model.Image, 0, len(dbItemImages))
	for _, dbImg := range dbItemImages {
		images = append(images, service.MapImageToModel(dbImg))
	}

	return renderOK(c, view.StoreAluminioItem(item, images))
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

	return renderOK(c, view.StoreUPVCIndex(cats))
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
		items = append(items, service.MapItemToModel(dbItem, dbCat, itemImage, catImage))
	}

	// Get PDF URL if exists
	pdfURL := ""
	if dbCat.PdfID.Valid {
		pdf, err := h.queries.GetStaticFile(ctx, dbCat.PdfID.Int32)
		if err == nil {
			pdfURL = "/uploads/pdfs/" + pdf.FileName
		}
	}

	return renderOK(c, view.StoreUPVCCategory(cat, items, pdfURL))
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

	item := service.MapItemToModel(dbItem, dbCat, itemImage, catImage)

	// Get item images (gallery)
	dbItemImages, err := h.queries.ListItemImages(ctx, dbItem.ItemID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al cargar galería")
	}

	images := make([]model.Image, 0, len(dbItemImages))
	for _, dbImg := range dbItemImages {
		images = append(images, service.MapImageToModel(dbImg))
	}

	return renderOK(c, view.StoreUPVCItem(item, images))
}
