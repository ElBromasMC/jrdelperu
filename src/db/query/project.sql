
-- Crear un nuevo proyecto
-- name: CreateProject :one
INSERT INTO projects (slug, description, location, period, area_m2, service, display_order, is_visible)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- Obtener un proyecto por ID
-- name: GetProject :one
SELECT * FROM projects
WHERE project_id = $1;

-- Obtener un proyecto por slug
-- name: GetProjectBySlug :one
SELECT * FROM projects
WHERE slug = $1;

-- Listar proyectos visibles (para la vista pública)
-- name: ListVisibleProjects :many
SELECT * FROM projects
WHERE is_visible = true
ORDER BY display_order, created_at DESC;

-- Listar todos los proyectos (para el admin)
-- name: ListAllProjects :many
SELECT * FROM projects
ORDER BY display_order, created_at DESC;

-- Actualizar un proyecto
-- name: UpdateProject :exec
UPDATE projects
SET slug = $2,
    description = $3,
    location = $4,
    period = $5,
    area_m2 = $6,
    service = $7,
    display_order = $8,
    is_visible = $9,
    updated_at = NOW()
WHERE project_id = $1;

-- Actualizar orden de visualización
-- name: UpdateProjectDisplayOrder :exec
UPDATE projects
SET display_order = $2, updated_at = NOW()
WHERE project_id = $1;

-- Actualizar visibilidad
-- name: UpdateProjectVisibility :exec
UPDATE projects
SET is_visible = $2, updated_at = NOW()
WHERE project_id = $1;

-- Eliminar un proyecto
-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;

-- Contar todos los proyectos
-- name: CountProjects :one
SELECT COUNT(*) FROM projects;

-- Contar proyectos visibles
-- name: CountVisibleProjects :one
SELECT COUNT(*) FROM projects
WHERE is_visible = true;
