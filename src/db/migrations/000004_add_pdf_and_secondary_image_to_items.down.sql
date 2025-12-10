BEGIN;

-- Remove indexes
DROP INDEX IF EXISTS idx_items_pdf_id;
DROP INDEX IF EXISTS idx_items_secondary_image_id;

-- Remove columns
ALTER TABLE items DROP COLUMN IF EXISTS pdf_id;
ALTER TABLE items DROP COLUMN IF EXISTS secondary_image_id;

COMMIT;
