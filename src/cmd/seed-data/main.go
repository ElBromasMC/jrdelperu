package main

import (
	"alc/config"
	"alc/repository"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Obtener la URL de PostgreSQL desde variable de entorno
	dbURL := os.Getenv("POSTGRESQL_URL")
	if dbURL == "" {
		log.Fatal("POSTGRESQL_URL environment variable is required")
	}

	// Conectar a la base de datos
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Crear queries
	queries := repository.New(dbpool)
	ctx := context.Background()

	// Crear archivo placeholder si no existe en la BD
	placeholderFile, err := queries.GetStaticFileByName(ctx, "placeholder.webp")
	if err != nil {
		// Si no existe, crearlo
		placeholderFile, err = queries.CreateStaticFile(ctx, repository.CreateStaticFileParams{
			FileName:      "placeholder.webp",
			FileType:      pgtype.Text{String: "image", Valid: true},
			MimeType:      pgtype.Text{String: "image/webp", Valid: true},
			FileSizeBytes: pgtype.Int8{Int64: 0, Valid: true},
			DisplayName:   pgtype.Text{String: "Placeholder Image", Valid: true},
		})
		if err != nil {
			log.Fatalf("Error creating placeholder file: %v", err)
		}
		log.Println("Created placeholder.webp file entry")
	}

	placeholderImageID := pgtype.Int4{Int32: placeholderFile.FileID, Valid: true}

	log.Println("Starting data seeding...")

	// Seed Vidrio Categories
	log.Println("\n=== Seeding Vidrio Categories ===")
	for _, cat := range config.VidrioCategories {
		category, err := queries.CreateCategory(ctx, repository.CreateCategoryParams{
			MaterialType:    repository.MaterialTypeVidrio,
			Slug:            cat.Slug,
			Name:            cat.Name,
			Description:     cat.Description,
			LongDescription: cat.LongDescription,
			ImageID:         placeholderImageID,
			TagID:           pgtype.Int4{},
			PdfID:           pgtype.Int4{},
		})
		if err != nil {
			log.Printf("Error creating vidrio category '%s': %v", cat.Name, err)
			continue
		}
		log.Printf("✓ Created category: %s (ID: %d)", category.Name, category.CategoryID)

		// Seed features for Vidrio Monolítico
		if cat.Id == 1 {
			log.Println("  Seeding features for Vidrio Monolítico...")
			for _, feat := range config.MonoliticoFeatures {
				feature, err := queries.CreateCategoryFeature(ctx, repository.CreateCategoryFeatureParams{
					CategoryID:  category.CategoryID,
					Name:        feat.Name,
					Description: feat.Description,
				})
				if err != nil {
					log.Printf("  Error creating feature '%s': %v", feat.Name, err)
					continue
				}
				log.Printf("  ✓ Created feature: %s (ID: %d)", feature.Name, feature.FeatureID)
			}

			// Seed items for Vidrio Monolítico
			log.Println("  Seeding items for Vidrio Monolítico...")
			for _, item := range config.MonoliticoItems {
				itemCreated, err := queries.CreateItem(ctx, repository.CreateItemParams{
					CategoryID:      category.CategoryID,
					Slug:            fmt.Sprintf("%d-mm", item.Id),
					Name:            item.Name,
					Description:     item.Description,
					LongDescription: item.Description,
					ImageID:         placeholderImageID,
				})
				if err != nil {
					log.Printf("  Error creating item '%s': %v", item.Name, err)
					continue
				}
				log.Printf("  ✓ Created item: %s (ID: %d)", itemCreated.Name, itemCreated.ItemID)
			}
		}
	}

	// Seed Aluminio Categories
	log.Println("\n=== Seeding Aluminio Categories ===")
	for _, cat := range config.AluminioCategories {
		category, err := queries.CreateCategory(ctx, repository.CreateCategoryParams{
			MaterialType:    repository.MaterialTypeAluminio,
			Slug:            cat.Slug,
			Name:            cat.Name,
			Description:     cat.Description,
			LongDescription: cat.Description,
			ImageID:         placeholderImageID,
			TagID:           pgtype.Int4{},
			PdfID:           pgtype.Int4{},
		})
		if err != nil {
			log.Printf("Error creating aluminio category '%s': %v", cat.Name, err)
			continue
		}
		log.Printf("✓ Created category: %s (ID: %d)", category.Name, category.CategoryID)

		// Seed items for Sistema de Fachadas ligeras
		if cat.Id == 1 {
			log.Println("  Seeding items for Sistema de Fachadas ligeras...")
			for _, item := range config.FachadasItems {
				itemCreated, err := queries.CreateItem(ctx, repository.CreateItemParams{
					CategoryID:      category.CategoryID,
					Slug:            item.Slug,
					Name:            item.Name,
					Description:     item.Name,
					LongDescription: item.LongDescription,
					ImageID:         placeholderImageID,
				})
				if err != nil {
					log.Printf("  Error creating item '%s': %v", item.Name, err)
					continue
				}
				log.Printf("  ✓ Created item: %s (ID: %d)", itemCreated.Name, itemCreated.ItemID)
			}
		}
	}

	// Seed uPVC Categories
	log.Println("\n=== Seeding uPVC Categories ===")
	for _, cat := range config.UPVCCategories {
		category, err := queries.CreateCategory(ctx, repository.CreateCategoryParams{
			MaterialType:    repository.MaterialTypeUpvc,
			Slug:            cat.Slug,
			Name:            cat.Name,
			Description:     cat.Description,
			LongDescription: cat.Description,
			ImageID:         placeholderImageID,
			TagID:           pgtype.Int4{},
			PdfID:           pgtype.Int4{},
		})
		if err != nil {
			log.Printf("Error creating uPVC category '%s': %v", cat.Name, err)
			continue
		}
		log.Printf("✓ Created category: %s (ID: %d)", category.Name, category.CategoryID)

		// Seed items for Lumina 60 Ventanas Articuladas
		if cat.Id == 1 {
			log.Println("  Seeding items for Lumina 60 Ventanas Articuladas...")
			for _, item := range config.LuminaItems {
				itemCreated, err := queries.CreateItem(ctx, repository.CreateItemParams{
					CategoryID:      category.CategoryID,
					Slug:            item.Slug,
					Name:            item.Name,
					Description:     item.Name,
					LongDescription: item.LongDescription,
					ImageID:         placeholderImageID,
				})
				if err != nil {
					log.Printf("  Error creating item '%s': %v", item.Name, err)
					continue
				}
				log.Printf("  ✓ Created item: %s (ID: %d)", itemCreated.Name, itemCreated.ItemID)
			}
		}
	}

	log.Println("\n=== Data seeding completed successfully! ===")
	log.Println("\nSummary:")

	// Count and display statistics
	vidrioCount, _ := queries.ListCategoriesByMaterialType(ctx, repository.MaterialTypeVidrio)
	aluminioCount, _ := queries.ListCategoriesByMaterialType(ctx, repository.MaterialTypeAluminio)
	upvcCount, _ := queries.ListCategoriesByMaterialType(ctx, repository.MaterialTypeUpvc)

	log.Printf("- Vidrio categories: %d", len(vidrioCount))
	log.Printf("- Aluminio categories: %d", len(aluminioCount))
	log.Printf("- uPVC categories: %d", len(upvcCount))

	allItems, _ := queries.ListAllItems(ctx)
	log.Printf("- Total items: %d", len(allItems))
}
