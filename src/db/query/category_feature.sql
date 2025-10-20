
-- Crear una nueva característica de categoría
-- name: CreateCategoryFeature :one
INSERT INTO category_features (category_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- Obtener una característica por ID
-- name: GetCategoryFeature :one
SELECT * FROM category_features
WHERE feature_id = $1;

-- Listar características de una categoría
-- name: ListCategoryFeatures :many
SELECT * FROM category_features
WHERE category_id = $1
ORDER BY name;

-- Actualizar una característica
-- name: UpdateCategoryFeature :exec
UPDATE category_features
SET name = $2, description = $3, updated_at = NOW()
WHERE feature_id = $1;

-- Eliminar una característica
-- name: DeleteCategoryFeature :exec
DELETE FROM category_features
WHERE feature_id = $1;

-- Eliminar todas las características de una categoría
-- name: DeleteAllCategoryFeatures :exec
DELETE FROM category_features
WHERE category_id = $1;
