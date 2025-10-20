
-- Crear un nuevo archivo estático
-- name: CreateStaticFile :one
INSERT INTO static_files (file_name, file_type, mime_type, file_size_bytes, display_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- Obtener un archivo por ID
-- name: GetStaticFile :one
SELECT * FROM static_files
WHERE file_id = $1;

-- Obtener un archivo por nombre
-- name: GetStaticFileByName :one
SELECT * FROM static_files
WHERE file_name = $1;

-- Listar todos los archivos de un tipo específico
-- name: ListStaticFilesByType :many
SELECT * FROM static_files
WHERE file_type = $1
ORDER BY created_at DESC;

-- Listar todas las imágenes
-- name: ListImages :many
SELECT * FROM static_files
WHERE file_type = 'image'
ORDER BY created_at DESC;

-- Listar todos los PDFs
-- name: ListPDFs :many
SELECT * FROM static_files
WHERE file_type = 'pdf'
ORDER BY created_at DESC;

-- Actualizar el nombre de visualización de un archivo
-- name: UpdateStaticFileDisplayName :exec
UPDATE static_files
SET display_name = $2
WHERE file_id = $1;

-- Eliminar un archivo
-- name: DeleteStaticFile :exec
DELETE FROM static_files
WHERE file_id = $1;

-- Contar imágenes
-- name: CountImages :one
SELECT COUNT(*) FROM static_files
WHERE file_type = 'image';

-- Contar PDFs
-- name: CountPDFs :one
SELECT COUNT(*) FROM static_files
WHERE file_type = 'pdf';
