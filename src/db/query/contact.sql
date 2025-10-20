
-- Crear una nueva solicitud de contacto
-- name: CreateContactSubmission :one
INSERT INTO contact_submissions (full_name, email, phone, subject, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- Obtener una solicitud de contacto por ID
-- name: GetContactSubmission :one
SELECT * FROM contact_submissions
WHERE submission_id = $1;

-- Listar todas las solicitudes de contacto (paginadas)
-- name: ListContactSubmissions :many
SELECT * FROM contact_submissions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- Listar solicitudes no leídas
-- name: ListUnreadContactSubmissions :many
SELECT * FROM contact_submissions
WHERE is_read = false
ORDER BY created_at DESC;

-- Contar solicitudes no leídas
-- name: CountUnreadContactSubmissions :one
SELECT COUNT(*) FROM contact_submissions
WHERE is_read = false;

-- Marcar solicitud como leída
-- name: MarkContactSubmissionAsRead :exec
UPDATE contact_submissions
SET is_read = true
WHERE submission_id = $1;

-- Marcar solicitud como no leída
-- name: MarkContactSubmissionAsUnread :exec
UPDATE contact_submissions
SET is_read = false
WHERE submission_id = $1;

-- Eliminar una solicitud de contacto
-- name: DeleteContactSubmission :exec
DELETE FROM contact_submissions
WHERE submission_id = $1;
