
BEGIN;

-- Add display_name column to static_files
ALTER TABLE static_files
ADD COLUMN display_name varchar(255);

-- Set default display name based on file_name for existing files
UPDATE static_files
SET display_name = file_name
WHERE display_name IS NULL;

COMMIT;
