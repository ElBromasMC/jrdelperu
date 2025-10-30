
-- Agregar una imagen a un proyecto
-- name: AddProjectImage :exec
INSERT INTO project_images (project_id, image_id, display_order, is_featured)
VALUES ($1, $2, $3, $4);

-- Remover una imagen de un proyecto
-- name: RemoveProjectImage :exec
DELETE FROM project_images
WHERE project_id = $1 AND image_id = $2;

-- Listar imágenes de un proyecto
-- name: ListProjectImages :many
SELECT pi.project_id, pi.image_id, pi.display_order, pi.is_featured,
       sf.file_name, sf.file_type, sf.mime_type, sf.display_name
FROM project_images pi
JOIN static_files sf ON pi.image_id = sf.file_id
WHERE pi.project_id = $1
ORDER BY pi.display_order, sf.created_at;

-- Actualizar orden de imagen
-- name: UpdateProjectImageOrder :exec
UPDATE project_images
SET display_order = $3
WHERE project_id = $1 AND image_id = $2;

-- Desmarcar todas las imágenes destacadas de un proyecto
-- name: UnsetAllFeaturedProjectImages :exec
UPDATE project_images
SET is_featured = false
WHERE project_id = $1;

-- Marcar imagen como destacada
-- name: SetFeaturedProjectImage :exec
UPDATE project_images
SET is_featured = true
WHERE project_id = $1 AND image_id = $2;

-- Actualizar estado de imagen destacada
-- name: UpdateProjectImageFeatured :exec
UPDATE project_images
SET is_featured = $3
WHERE project_id = $1 AND image_id = $2;

-- Obtener la imagen destacada de un proyecto
-- name: GetFeaturedProjectImage :one
SELECT pi.project_id, pi.image_id, pi.display_order, pi.is_featured,
       sf.file_name, sf.file_type, sf.mime_type, sf.display_name
FROM project_images pi
JOIN static_files sf ON pi.image_id = sf.file_id
WHERE pi.project_id = $1 AND pi.is_featured = true
LIMIT 1;

-- Contar imágenes de un proyecto
-- name: CountProjectImages :one
SELECT COUNT(*) FROM project_images
WHERE project_id = $1;

-- Eliminar todas las imágenes de un proyecto
-- name: RemoveAllProjectImages :exec
DELETE FROM project_images
WHERE project_id = $1;
