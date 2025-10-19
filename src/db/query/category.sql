
-- name: GetCategory :one
SELECT * FROM categories
WHERE category_id = $1;

