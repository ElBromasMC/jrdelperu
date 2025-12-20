-- Add secondary image column to categories for natura-style rendering
ALTER TABLE categories ADD COLUMN secondary_image_id int REFERENCES static_files ON DELETE SET NULL;
