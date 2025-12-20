
-- Crear un nuevo artículo
-- name: CreateArticle :one
INSERT INTO blog_articles (slug, title, summary, content, cover_image_id, author)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at;

-- Obtener un artículo por ID
-- name: GetArticle :one
SELECT article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at FROM blog_articles
WHERE article_id = $1;

-- Obtener un artículo publicado por slug
-- name: GetPublishedArticleBySlug :one
SELECT article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at FROM blog_articles
WHERE slug = $1 AND is_published = true;

-- Obtener un artículo por slug (para admin)
-- name: GetArticleBySlug :one
SELECT article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at FROM blog_articles
WHERE slug = $1;

-- Listar artículos publicados (para público)
-- name: ListPublishedArticles :many
SELECT article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at FROM blog_articles
WHERE is_published = true
ORDER BY published_at DESC;

-- Listar artículos publicados con paginación
-- name: ListPublishedArticlesPaginated :many
SELECT article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at FROM blog_articles
WHERE is_published = true
ORDER BY published_at DESC
LIMIT $1 OFFSET $2;

-- Listar todos los artículos (para admin)
-- name: ListAllArticles :many
SELECT article_id, slug, title, summary, content, cover_image_id, author, is_published, published_at, created_at, updated_at FROM blog_articles
ORDER BY created_at DESC;

-- Actualizar un artículo
-- name: UpdateArticle :exec
UPDATE blog_articles
SET slug = $2,
    title = $3,
    summary = $4,
    content = $5,
    cover_image_id = $6,
    author = $7,
    updated_at = NOW()
WHERE article_id = $1;

-- Publicar un artículo
-- name: PublishArticle :exec
UPDATE blog_articles
SET is_published = true,
    published_at = NOW(),
    updated_at = NOW()
WHERE article_id = $1;

-- Despublicar un artículo
-- name: UnpublishArticle :exec
UPDATE blog_articles
SET is_published = false,
    updated_at = NOW()
WHERE article_id = $1;

-- Actualizar imagen de portada
-- name: UpdateArticleCoverImage :exec
UPDATE blog_articles
SET cover_image_id = $2, updated_at = NOW()
WHERE article_id = $1;

-- Eliminar un artículo
-- name: DeleteArticle :exec
DELETE FROM blog_articles
WHERE article_id = $1;

-- Contar artículos publicados
-- name: CountPublishedArticles :one
SELECT COUNT(*) FROM blog_articles
WHERE is_published = true;

-- Contar todos los artículos
-- name: CountAllArticles :one
SELECT COUNT(*) FROM blog_articles;

-- Buscar artículos (full-text search)
-- name: SearchArticles :many
SELECT article_id, slug, title, summary, cover_image_id, author, published_at,
       ts_rank(search_vector, plainto_tsquery('spanish', $1)) AS rank
FROM blog_articles
WHERE is_published = true AND search_vector @@ plainto_tsquery('spanish', $1)
ORDER BY rank DESC
LIMIT 20;

