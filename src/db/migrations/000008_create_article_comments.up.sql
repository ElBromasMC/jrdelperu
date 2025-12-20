BEGIN;

-- Article comments table
CREATE TABLE IF NOT EXISTS article_comments (
    comment_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    article_id int NOT NULL REFERENCES blog_articles ON DELETE CASCADE,
    parent_id int REFERENCES article_comments ON DELETE CASCADE,  -- For admin replies
    author_name varchar(100) NOT NULL,
    author_email varchar(255),  -- Optional
    content text NOT NULL,
    is_admin boolean NOT NULL DEFAULT false,  -- True for admin replies
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- Indices
CREATE INDEX idx_article_comments_article ON article_comments(article_id, created_at);
CREATE INDEX idx_article_comments_parent ON article_comments(parent_id);

COMMIT;
