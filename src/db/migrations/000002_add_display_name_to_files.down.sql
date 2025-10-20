
BEGIN;

-- Remove display_name column from static_files
ALTER TABLE static_files
DROP COLUMN IF EXISTS display_name;

COMMIT;
