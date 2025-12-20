package model

import "time"

// Article representa un artículo del blog
type Article struct {
	ID          int
	Slug        string
	Title       string
	Summary     string
	Content     string // HTML from Quill editor
	CoverImage  Image
	Author      string
	IsPublished bool
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FAQs        []ArticleFAQ
}

// ArticleFAQ representa una pregunta frecuente de un artículo
type ArticleFAQ struct {
	ID           int
	ArticleID    int
	Question     string
	Answer       string
	DisplayOrder int
}

// GlobalFAQ representa una pregunta frecuente global del sitio
type GlobalFAQ struct {
	ID           int
	Category     string
	Question     string
	Answer       string
	DisplayOrder int
	IsVisible    bool
}

// PageMeta contiene metadatos SEO para las páginas
type PageMeta struct {
	Title       string
	Description string
	OGImage     string
	OGType      string // "article" or "website"
	Canonical   string
	PublishedAt string // ISO 8601 format for articles
	Author      string
}

// DefaultPageMeta retorna metadatos por defecto para páginas estándar
func DefaultPageMeta(title, description string) PageMeta {
	return PageMeta{
		Title:       title,
		Description: description,
		OGType:      "website",
		Author:      "JR del Perú",
	}
}

// ArticlePageMeta retorna metadatos SEO para un artículo
func ArticlePageMeta(article Article, baseURL string) PageMeta {
	meta := PageMeta{
		Title:       article.Title,
		Description: article.Summary,
		OGType:      "article",
		Author:      article.Author,
		Canonical:   baseURL + "/blog/" + article.Slug,
	}

	if !article.PublishedAt.IsZero() {
		meta.PublishedAt = article.PublishedAt.Format(time.RFC3339)
	}

	if article.CoverImage.Id != 0 {
		meta.OGImage = baseURL + "/uploads/images/" + article.CoverImage.Filename
	}

	return meta
}
