
-- Listar todos los documentos del sitio (para admin)
-- name: ListSiteDocuments :many
SELECT * FROM site_documents
ORDER BY document_id;

-- Obtener un documento por su clave
-- name: GetSiteDocumentByKey :one
SELECT * FROM site_documents
WHERE document_key = $1;

-- Obtener múltiples documentos por sus claves
-- name: GetSiteDocumentsByKeys :many
SELECT * FROM site_documents
WHERE document_key = ANY($1::varchar[])
ORDER BY document_id;

-- Actualizar el archivo de un documento
-- name: UpdateSiteDocumentFile :exec
UPDATE site_documents
SET file_id = $2, updated_at = NOW()
WHERE document_id = $1;

-- Actualizar el nombre de visualización de un documento
-- name: UpdateSiteDocumentDisplayName :exec
UPDATE site_documents
SET display_name = $2, updated_at = NOW()
WHERE document_id = $1;

-- Actualizar documento completo (archivo y nombre)
-- name: UpdateSiteDocument :exec
UPDATE site_documents
SET display_name = $2, file_id = $3, updated_at = NOW()
WHERE document_id = $1;
