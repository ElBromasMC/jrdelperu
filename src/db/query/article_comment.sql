-- Crear un comentario (admin o público, puede ser respuesta si parent_id está presente)
-- name: CreateArticleComment :one
INSERT INTO article_comments (article_id, parent_id, author_name, author_email, content, admin_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- Crear respuesta de admin (legacy, usar CreateArticleComment con admin_id)
-- name: CreateAdminReply :one
INSERT INTO article_comments (article_id, parent_id, author_name, content, admin_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- Listar comentarios de un artículo con info de admin (solo principales)
-- name: ListArticleCommentsWithAdmin :many
SELECT
    c.comment_id,
    c.article_id,
    c.parent_id,
    c.author_name,
    c.author_email,
    c.content,
    c.admin_id,
    c.created_at,
    a.username AS admin_username
FROM article_comments c
LEFT JOIN admins a ON c.admin_id = a.admin_id
WHERE c.article_id = $1 AND c.parent_id IS NULL
ORDER BY c.created_at DESC;

-- Listar respuestas de un comentario con info de admin
-- name: ListCommentRepliesWithAdmin :many
SELECT
    c.comment_id,
    c.article_id,
    c.parent_id,
    c.author_name,
    c.author_email,
    c.content,
    c.admin_id,
    c.created_at,
    a.username AS admin_username
FROM article_comments c
LEFT JOIN admins a ON c.admin_id = a.admin_id
WHERE c.parent_id = $1
ORDER BY c.created_at ASC;

-- Obtener un comentario con info de admin
-- name: GetArticleCommentWithAdmin :one
SELECT
    c.comment_id,
    c.article_id,
    c.parent_id,
    c.author_name,
    c.author_email,
    c.content,
    c.admin_id,
    c.created_at,
    a.username AS admin_username
FROM article_comments c
LEFT JOIN admins a ON c.admin_id = a.admin_id
WHERE c.comment_id = $1;

-- Listar comentarios de un artículo (solo los principales, sin JOIN - legacy)
-- name: ListArticleComments :many
SELECT * FROM article_comments
WHERE article_id = $1 AND parent_id IS NULL
ORDER BY created_at DESC;

-- Listar respuestas de un comentario (sin JOIN - legacy)
-- name: ListCommentReplies :many
SELECT * FROM article_comments
WHERE parent_id = $1
ORDER BY created_at ASC;

-- Obtener un comentario por ID (sin JOIN - legacy)
-- name: GetArticleComment :one
SELECT * FROM article_comments
WHERE comment_id = $1;

-- Eliminar un comentario (las respuestas se eliminan por CASCADE)
-- name: DeleteArticleComment :exec
DELETE FROM article_comments WHERE comment_id = $1;

-- Contar comentarios de un artículo
-- name: CountArticleComments :one
SELECT COUNT(*) FROM article_comments WHERE article_id = $1;

-- Listar todos los comentarios de un artículo con info de admin (para admin panel)
-- name: ListAllArticleCommentsWithAdmin :many
SELECT
    c.comment_id,
    c.article_id,
    c.parent_id,
    c.author_name,
    c.author_email,
    c.content,
    c.admin_id,
    c.created_at,
    a.username AS admin_username
FROM article_comments c
LEFT JOIN admins a ON c.admin_id = a.admin_id
WHERE c.article_id = $1
ORDER BY c.created_at DESC;
