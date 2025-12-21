BEGIN;

DROP INDEX IF EXISTS idx_article_comments_admin;
ALTER TABLE article_comments DROP COLUMN IF EXISTS admin_id;

COMMIT;
