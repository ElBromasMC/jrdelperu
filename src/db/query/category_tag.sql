
-- Crear un nuevo tag de categoría
-- name: CreateCategoryTag :one
INSERT INTO categories_tags (tag_name, position_num)
VALUES ($1, $2)
RETURNING *;

-- Obtener un tag por ID
-- name: GetCategoryTag :one
SELECT * FROM categories_tags
WHERE tag_id = $1;

-- Obtener un tag por nombre
-- name: GetCategoryTagByName :one
SELECT * FROM categories_tags
WHERE tag_name = $1;

-- Listar todos los tags ordenados por posición
-- name: ListCategoryTags :many
SELECT * FROM categories_tags
ORDER BY position_num;

-- Actualizar un tag
-- name: UpdateCategoryTag :exec
UPDATE categories_tags
SET tag_name = $2, position_num = $3, updated_at = NOW()
WHERE tag_id = $1;

-- Eliminar un tag
-- name: DeleteCategoryTag :exec
DELETE FROM categories_tags
WHERE tag_id = $1;
