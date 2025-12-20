-- Crear un comentario público
-- name: CreateArticleComment :one
INSERT INTO article_comments (article_id, author_name, author_email, content, is_admin)
VALUES ($1, $2, $3, $4, false)
RETURNING *;

-- Crear respuesta de admin
-- name: CreateAdminReply :one
INSERT INTO article_comments (article_id, parent_id, author_name, content, is_admin)
VALUES ($1, $2, $3, $4, true)
RETURNING *;

-- Listar comentarios de un artículo (solo los principales, no respuestas)
-- name: ListArticleComments :many
SELECT * FROM article_comments
WHERE article_id = $1 AND parent_id IS NULL
ORDER BY created_at DESC;

-- Listar respuestas de un comentario
-- name: ListCommentReplies :many
SELECT * FROM article_comments
WHERE parent_id = $1
ORDER BY created_at ASC;

-- Obtener un comentario por ID
-- name: GetArticleComment :one
SELECT * FROM article_comments
WHERE comment_id = $1;

-- Eliminar un comentario (las respuestas se eliminan por CASCADE)
-- name: DeleteArticleComment :exec
DELETE FROM article_comments WHERE comment_id = $1;

-- Contar comentarios de un artículo
-- name: CountArticleComments :one
SELECT COUNT(*) FROM article_comments WHERE article_id = $1;

-- Listar todos los comentarios de un artículo (para admin)
-- name: ListAllArticleComments :many
SELECT * FROM article_comments
WHERE article_id = $1
ORDER BY created_at DESC;
