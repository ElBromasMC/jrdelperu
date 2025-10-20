
BEGIN;

-- Eliminar índices
DROP INDEX IF EXISTS idx_item_images_position;
DROP INDEX IF EXISTS idx_items_slug;
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_material_type;
DROP INDEX IF EXISTS idx_categories_tags_position;
DROP INDEX IF EXISTS idx_contact_submissions_is_read;
DROP INDEX IF EXISTS idx_contact_submissions_created_at;
DROP INDEX IF EXISTS idx_admins_email;
DROP INDEX IF EXISTS idx_admins_username;
DROP INDEX IF EXISTS idx_static_files_file_type;

-- Eliminar tablas en orden inverso de dependencias
DROP TABLE IF EXISTS item_images;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS category_features;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS categories_tags;
DROP TABLE IF EXISTS contact_submissions;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS static_files;

-- Eliminar tipo ENUM
DROP TYPE IF EXISTS material_type;

COMMIT;
