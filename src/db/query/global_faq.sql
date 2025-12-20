
-- Crear una pregunta frecuente global
-- name: CreateGlobalFAQ :one
INSERT INTO global_faqs (category, question, answer, display_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- Obtener una pregunta frecuente global por ID
-- name: GetGlobalFAQ :one
SELECT * FROM global_faqs
WHERE faq_id = $1;

-- Listar preguntas frecuentes globales visibles (para público)
-- name: ListVisibleGlobalFAQs :many
SELECT * FROM global_faqs
WHERE is_visible = true
ORDER BY display_order, faq_id;

-- Listar preguntas frecuentes globales por categoría
-- name: ListVisibleGlobalFAQsByCategory :many
SELECT * FROM global_faqs
WHERE is_visible = true AND category = $1
ORDER BY display_order, faq_id;

-- Listar todas las categorías de preguntas frecuentes
-- name: ListGlobalFAQCategories :many
SELECT DISTINCT category FROM global_faqs
WHERE category IS NOT NULL AND category != ''
ORDER BY category;

-- Listar todas las preguntas frecuentes globales (para admin)
-- name: ListAllGlobalFAQs :many
SELECT * FROM global_faqs
ORDER BY display_order, faq_id;

-- Actualizar una pregunta frecuente global
-- name: UpdateGlobalFAQ :exec
UPDATE global_faqs
SET category = $2,
    question = $3,
    answer = $4,
    display_order = $5,
    is_visible = $6,
    updated_at = NOW()
WHERE faq_id = $1;

-- Cambiar visibilidad de una pregunta frecuente global
-- name: SetGlobalFAQVisibility :exec
UPDATE global_faqs
SET is_visible = $2, updated_at = NOW()
WHERE faq_id = $1;

-- Eliminar una pregunta frecuente global
-- name: DeleteGlobalFAQ :exec
DELETE FROM global_faqs
WHERE faq_id = $1;

-- Contar preguntas frecuentes globales visibles
-- name: CountVisibleGlobalFAQs :one
SELECT COUNT(*) FROM global_faqs
WHERE is_visible = true;

-- Contar todas las preguntas frecuentes globales
-- name: CountAllGlobalFAQs :one
SELECT COUNT(*) FROM global_faqs;

-- Obtener el máximo display_order global
-- name: GetMaxGlobalFAQOrder :one
SELECT COALESCE(MAX(display_order), 0) FROM global_faqs;

