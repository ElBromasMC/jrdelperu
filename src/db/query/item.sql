
-- Crear un nuevo item
-- name: CreateItem :one
INSERT INTO items (category_id, slug, name, description, long_description, image_id)
VALUES ($1, $2, $3, $4, $5, $6)
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
    updated_at = NOW()
WHERE item_id = $1;

-- Actualizar la imagen principal de un item
-- name: UpdateItemImage :exec
UPDATE items
SET image_id = $2, updated_at = NOW()
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
