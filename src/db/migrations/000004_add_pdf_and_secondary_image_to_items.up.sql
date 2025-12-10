BEGIN;

-- Add pdf_id to items table (like categories have)
ALTER TABLE items
ADD COLUMN pdf_id int REFERENCES static_files ON DELETE SET NULL;

-- Add secondary_image_id to items table for the second image display
ALTER TABLE items
ADD COLUMN secondary_image_id int REFERENCES static_files ON DELETE SET NULL;

-- Index for faster lookups
CREATE INDEX idx_items_pdf_id ON items(pdf_id);
CREATE INDEX idx_items_secondary_image_id ON items(secondary_image_id);

COMMIT;
