BEGIN;

-- Blog articles table
CREATE TABLE IF NOT EXISTS blog_articles (
    article_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug varchar(255) UNIQUE NOT NULL,
    title varchar(255) NOT NULL,
    summary text NOT NULL,
    content text NOT NULL,
    cover_image_id int REFERENCES static_files ON DELETE SET NULL,
    author varchar(100) NOT NULL DEFAULT 'JR del Perú',
    is_published boolean NOT NULL DEFAULT false,
    published_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    -- Full-text search vector (Spanish language)
    search_vector tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('spanish', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('spanish', coalesce(summary, '')), 'B') ||
        setweight(to_tsvector('spanish', coalesce(content, '')), 'C')
    ) STORED
);

-- Article FAQ (per-article Q&A)
CREATE TABLE IF NOT EXISTS article_faqs (
    faq_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    article_id int NOT NULL REFERENCES blog_articles ON DELETE CASCADE,
    question text NOT NULL,
    answer text NOT NULL,
    display_order int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- Global FAQ (site-wide Q&A)
CREATE TABLE IF NOT EXISTS global_faqs (
    faq_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    category varchar(100),
    question text NOT NULL,
    answer text NOT NULL,
    display_order int NOT NULL DEFAULT 0,
    is_visible boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Indices for performance
CREATE INDEX idx_blog_articles_slug ON blog_articles(slug);
CREATE INDEX idx_blog_articles_published ON blog_articles(is_published, published_at DESC);
CREATE INDEX idx_blog_articles_search ON blog_articles USING GIN(search_vector);
CREATE INDEX idx_article_faqs_article ON article_faqs(article_id, display_order);
CREATE INDEX idx_global_faqs_order ON global_faqs(display_order);
CREATE INDEX idx_global_faqs_visible ON global_faqs(is_visible, display_order);

COMMIT;
