
-- Crear un nuevo item
-- name: CreateItem :one
INSERT INTO items (category_id, slug, name, description, long_description, image_id, secondary_image_id, pdf_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- Obtener un item por ID
-- name: GetItem :one
SELECT * FROM items
WHERE item_id = $1;

-- Obtener un item por slug y categoría
-- name: GetItemBySlug :one
SELECT * FROM items
WHERE category_id = $1 AND slug = $2;

-- Listar items de una categoría
-- name: ListItemsByCategory :many
SELECT * FROM items
WHERE category_id = $1
ORDER BY name;

-- Listar todos los items
-- name: ListAllItems :many
SELECT * FROM items
ORDER BY category_id, name;

-- Actualizar un item
-- name: UpdateItem :exec
UPDATE items
SET slug = $2,
    name = $3,
    description = $4,
    long_description = $5,
    image_id = $6,
    secondary_image_id = $7,
    pdf_id = $8,
    updated_at = NOW()
WHERE item_id = $1;

-- Actualizar la imagen principal de un item
-- name: UpdateItemImage :exec
UPDATE items
SET image_id = $2, updated_at = NOW()
WHERE item_id = $1;

-- Actualizar la imagen secundaria de un item
-- name: UpdateItemSecondaryImage :exec
UPDATE items
SET secondary_image_id = $2, updated_at = NOW()
WHERE item_id = $1;

-- Actualizar el PDF de un item
-- name: UpdateItemPDF :exec
UPDATE items
SET pdf_id = $2, updated_at = NOW()
WHERE item_id = $1;

-- Eliminar un item
-- name: DeleteItem :exec
DELETE FROM items
WHERE item_id = $1;

-- Contar items por categoría
-- name: CountItemsByCategory :one
SELECT COUNT(*) FROM items
WHERE category_id = $1;

-- Contar todos los items
-- name: CountAllItems :one
SELECT COUNT(*) FROM items;

-- Listar items con PDF por tipo de material de la categoría
-- name: ListItemsWithPDFByMaterialType :many
SELECT i.*, c.name as category_name, sf.file_id as pdf_file_id, sf.file_name as pdf_file_name, sf.display_name as pdf_display_name
FROM items i
JOIN categories c ON i.category_id = c.category_id
JOIN static_files sf ON i.pdf_id = sf.file_id
WHERE c.material_type = $1 AND i.pdf_id IS NOT NULL
ORDER BY c.name, i.name;
