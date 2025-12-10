package service

import (
	"alc/model"
	"alc/repository"
)

// MapCategoryToModel convierte una categoría de la base de datos al modelo
func MapCategoryToModel(dbCat repository.Category, dbImage *repository.StaticFile) model.Category {
	cat := model.Category{
		Id:              int(dbCat.CategoryID),
		Slug:            dbCat.Slug,
		Name:            dbCat.Name,
		Description:     dbCat.Description,
		LongDescription: dbCat.LongDescription,
	}

	// Map material type
	switch dbCat.MaterialType {
	case repository.MaterialTypeVidrio:
		cat.Type = model.VidrioType
	case repository.MaterialTypeAluminio:
		cat.Type = model.AluminioType
	case repository.MaterialTypeUpvc:
		cat.Type = model.UPVCType
	}

	// Map image if exists
	if dbImage != nil {
		cat.Img = model.Image{
			Id:       int(dbImage.FileID),
			Filename: dbImage.FileName,
		}
	}

	return cat
}

// MapItemToModel convierte un item de la base de datos al modelo
func MapItemToModel(dbItem repository.Item, dbCat repository.Category, dbImage *repository.StaticFile, dbSecondaryImage *repository.StaticFile, dbCatImage *repository.StaticFile) model.Item {
	item := model.Item{
		Id:              int(dbItem.ItemID),
		Slug:            dbItem.Slug,
		Name:            dbItem.Name,
		Description:     dbItem.Description,
		LongDescription: dbItem.LongDescription,
		Category:        MapCategoryToModel(dbCat, dbCatImage),
	}

	// Map image if exists
	if dbImage != nil {
		item.Img = model.Image{
			Id:       int(dbImage.FileID),
			Filename: dbImage.FileName,
		}
	}

	// Map secondary image if exists
	if dbSecondaryImage != nil {
		item.SecondaryImg = model.Image{
			Id:       int(dbSecondaryImage.FileID),
			Filename: dbSecondaryImage.FileName,
		}
	}

	return item
}

// MapCategoryFeatureToModel convierte una característica de categoría al modelo
func MapCategoryFeatureToModel(dbFeature repository.CategoryFeature, dbCat repository.Category, dbCatImage *repository.StaticFile) model.CategoryFeature {
	return model.CategoryFeature{
		Id:          int(dbFeature.FeatureID),
		Name:        dbFeature.Name,
		Description: dbFeature.Description,
		Category:    MapCategoryToModel(dbCat, dbCatImage),
	}
}

// MapImageToModel convierte un archivo estático al modelo de imagen
func MapImageToModel(dbFile repository.StaticFile) model.Image {
	return model.Image{
		Id:       int(dbFile.FileID),
		Filename: dbFile.FileName,
	}
}
