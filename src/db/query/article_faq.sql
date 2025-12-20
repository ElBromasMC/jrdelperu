
-- Crear una pregunta frecuente para un artículo
-- name: CreateArticleFAQ :one
INSERT INTO article_faqs (article_id, question, answer, display_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- Obtener una pregunta frecuente por ID
-- name: GetArticleFAQ :one
SELECT * FROM article_faqs
WHERE faq_id = $1;

-- Listar preguntas frecuentes de un artículo
-- name: ListArticleFAQs :many
SELECT * FROM article_faqs
WHERE article_id = $1
ORDER BY display_order, faq_id;

-- Actualizar una pregunta frecuente
-- name: UpdateArticleFAQ :exec
UPDATE article_faqs
SET question = $2,
    answer = $3,
    display_order = $4
WHERE faq_id = $1;

-- Eliminar una pregunta frecuente
-- name: DeleteArticleFAQ :exec
DELETE FROM article_faqs
WHERE faq_id = $1;

-- Eliminar todas las preguntas frecuentes de un artículo
-- name: DeleteArticleFAQsByArticle :exec
DELETE FROM article_faqs
WHERE article_id = $1;

-- Contar preguntas frecuentes de un artículo
-- name: CountArticleFAQs :one
SELECT COUNT(*) FROM article_faqs
WHERE article_id = $1;

-- Obtener el máximo display_order de un artículo
-- name: GetMaxArticleFAQOrder :one
SELECT COALESCE(MAX(display_order), 0) FROM article_faqs
WHERE article_id = $1;

