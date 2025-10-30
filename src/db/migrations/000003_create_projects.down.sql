
BEGIN;

-- Eliminar índices
DROP INDEX IF EXISTS idx_project_images_featured;
DROP INDEX IF EXISTS idx_project_images_display_order;
DROP INDEX IF EXISTS idx_projects_is_visible;
DROP INDEX IF EXISTS idx_projects_display_order;
DROP INDEX IF EXISTS idx_projects_slug;

-- Eliminar tablas
DROP TABLE IF EXISTS project_images;
DROP TABLE IF EXISTS projects;

COMMIT;
