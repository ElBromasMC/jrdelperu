
-- Crear una nueva reclamación (Libro de Reclamaciones)
-- name: CreateComplaint :one
INSERT INTO complaints (
    full_name, document_number, address, phone, email,
    good_type, good_description, claim_type, detail, request, registered_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- Obtener una reclamación por ID
-- name: GetComplaint :one
SELECT * FROM complaints
WHERE complaint_id = $1;

-- Listar reclamaciones (paginadas)
-- name: ListComplaints :many
SELECT * FROM complaints
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- Buscar reclamaciones por nombre, documento o correo (paginadas)
-- name: SearchComplaints :many
SELECT * FROM complaints
WHERE full_name ILIKE '%' || @search::text || '%'
   OR document_number ILIKE '%' || @search::text || '%'
   OR email ILIKE '%' || @search::text || '%'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- Actualizar las observaciones de la empresa y el estado de una reclamación
-- name: UpdateComplaintCompanyNotes :exec
UPDATE complaints
SET company_notes = $2,
    is_resolved = $3,
    updated_at = NOW()
WHERE complaint_id = $1;

-- Eliminar una reclamación
-- name: DeleteComplaint :exec
DELETE FROM complaints
WHERE complaint_id = $1;

-- Contar todas las reclamaciones
-- name: CountComplaints :one
SELECT COUNT(*) FROM complaints;

-- Contar reclamaciones que coinciden con una búsqueda
-- name: CountSearchComplaints :one
SELECT COUNT(*) FROM complaints
WHERE full_name ILIKE '%' || @search::text || '%'
   OR document_number ILIKE '%' || @search::text || '%'
   OR email ILIKE '%' || @search::text || '%';

-- Contar reclamaciones sin resolver
-- name: CountUnresolvedComplaints :one
SELECT COUNT(*) FROM complaints
WHERE is_resolved = false;
