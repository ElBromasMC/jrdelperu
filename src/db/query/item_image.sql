
-- Agregar una imagen a un item
-- name: AddItemImage :exec
INSERT INTO item_images (item_id, image_id, position_num)
VALUES ($1, $2, $3);

-- Listar imágenes de un item ordenadas por posición
-- name: ListItemImages :many
SELECT sf.* FROM static_files sf
JOIN item_images ii ON sf.file_id = ii.image_id
WHERE ii.item_id = $1
ORDER BY ii.position_num;

-- Obtener la cantidad de imágenes de un item
-- name: CountItemImages :one
SELECT COUNT(*) FROM item_images
WHERE item_id = $1;

-- Actualizar la posición de una imagen en un item
-- name: UpdateItemImagePosition :exec
UPDATE item_images
SET position_num = $3
WHERE item_id = $1 AND image_id = $2;

-- Eliminar una imagen de un item
-- name: RemoveItemImage :exec
DELETE FROM item_images
WHERE item_id = $1 AND image_id = $2;

-- Eliminar todas las imágenes de un item
-- name: RemoveAllItemImages :exec
DELETE FROM item_images
WHERE item_id = $1;

-- Verificar si una imagen está asociada a un item
-- name: ItemImageExists :one
SELECT EXISTS(
    SELECT 1 FROM item_images
    WHERE item_id = $1 AND image_id = $2
);
