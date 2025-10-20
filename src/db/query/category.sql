
-- Crear una nueva categoría
-- name: CreateCategory :one
INSERT INTO categories (material_type, slug, name, description, long_description, image_id, tag_id, pdf_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- Obtener una categoría por ID
-- name: GetCategory :one
SELECT * FROM categories
WHERE category_id = $1;

-- Obtener una categoría por slug y tipo de material
-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE material_type = $1 AND slug = $2;

-- Listar categorías por tipo de material
-- name: ListCategoriesByMaterialType :many
SELECT * FROM categories
WHERE material_type = $1
ORDER BY name;

-- Listar todas las categorías
-- name: ListAllCategories :many
SELECT * FROM categories
ORDER BY material_type, name;

-- Listar categorías por tag
-- name: ListCategoriesByTag :many
SELECT c.* FROM categories c
JOIN categories_tags t ON c.tag_id = t.tag_id
WHERE t.tag_id = $1
ORDER BY c.name;

-- Actualizar una categoría
-- name: UpdateCategory :exec
UPDATE categories
SET slug = $2,
    name = $3,
    description = $4,
    long_description = $5,
    image_id = $6,
    tag_id = $7,
    pdf_id = $8,
    updated_at = NOW()
WHERE category_id = $1;

-- Actualizar la imagen de una categoría
-- name: UpdateCategoryImage :exec
UPDATE categories
SET image_id = $2, updated_at = NOW()
WHERE category_id = $1;

-- Actualizar el PDF de una categoría
-- name: UpdateCategoryPDF :exec
UPDATE categories
SET pdf_id = $2, updated_at = NOW()
WHERE category_id = $1;

-- Eliminar una categoría
-- name: DeleteCategory :exec
DELETE FROM categories
WHERE category_id = $1;

-- Contar todas las categorías
-- name: CountCategories :one
SELECT COUNT(*) FROM categories;

