-- Remove secondary image column from categories
ALTER TABLE categories DROP COLUMN IF EXISTS secondary_image_id;
