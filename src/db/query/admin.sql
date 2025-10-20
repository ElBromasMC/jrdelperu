
-- Crear un nuevo administrador
-- name: CreateAdmin :one
INSERT INTO admins (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- Obtener un administrador por ID
-- name: GetAdminByID :one
SELECT * FROM admins
WHERE admin_id = $1;

-- Obtener un administrador por username
-- name: GetAdminByUsername :one
SELECT * FROM admins
WHERE username = $1;

-- Obtener un administrador por email
-- name: GetAdminByEmail :one
SELECT * FROM admins
WHERE email = $1;

-- Listar todos los administradores activos
-- name: ListActiveAdmins :many
SELECT * FROM admins
WHERE is_active = true
ORDER BY username;

-- Actualizar contraseña de administrador
-- name: UpdateAdminPassword :exec
UPDATE admins
SET password_hash = $2, updated_at = NOW()
WHERE admin_id = $1;

-- Actualizar información de administrador
-- name: UpdateAdmin :exec
UPDATE admins
SET username = $2, email = $3, updated_at = NOW()
WHERE admin_id = $1;

-- Desactivar un administrador (soft delete)
-- name: DeactivateAdmin :exec
UPDATE admins
SET is_active = false, updated_at = NOW()
WHERE admin_id = $1;

-- Activar un administrador
-- name: ActivateAdmin :exec
UPDATE admins
SET is_active = true, updated_at = NOW()
WHERE admin_id = $1;

-- Eliminar un administrador (hard delete)
-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE admin_id = $1;
