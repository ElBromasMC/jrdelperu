BEGIN;

-- Add admin_id column with foreign key
ALTER TABLE article_comments
ADD COLUMN admin_id int REFERENCES admins ON DELETE SET NULL;

-- Attempt to link existing admin comments by matching author_name to admin username
UPDATE article_comments ac
SET admin_id = a.admin_id
FROM admins a
WHERE ac.is_admin = true
AND LOWER(ac.author_name) = LOWER(a.username);

-- Create index for the new foreign key
CREATE INDEX idx_article_comments_admin ON article_comments(admin_id);

COMMIT;
