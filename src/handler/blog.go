package handler

import (
	"alc/config"
	"alc/model"
	"alc/repository"
	"alc/view"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// HandleBlogIndex muestra la página principal del blog con lista de artículos
func (h *Handler) HandleBlogIndex(c echo.Context) error {
	ctx := c.Request().Context()

	// Get pagination params
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := int32(9) // Articles per page
	offset := int32((page - 1)) * limit

	// Get published articles
	articles, err := h.queries.ListPublishedArticlesPaginated(ctx, repository.ListPublishedArticlesPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		articles = []repository.ListPublishedArticlesPaginatedRow{}
	}

	// Get total count for pagination
	totalCount, _ := h.queries.CountPublishedArticles(ctx)
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	// Load cover images
	imagesMap := make(map[int32]repository.StaticFile)
	for _, article := range articles {
		if article.CoverImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, article.CoverImageID.Int32)
			if err == nil {
				imagesMap[article.CoverImageID.Int32] = img
			}
		}
	}

	// SEO meta
	meta := model.DefaultPageMeta(
		"Blog",
		"Artículos sobre vidrio, aluminio, uPVC y soluciones arquitectónicas. Información técnica y novedades de JR del Perú.",
	)
	baseURL := getBaseURL()
	meta.Canonical = baseURL + "/blog"

	return Render(c, http.StatusOK, view.BlogIndexPage(meta, articles, imagesMap, page, totalPages))
}

// HandleBlogSearch busca artículos en el blog
func (h *Handler) HandleBlogSearch(c echo.Context) error {
	ctx := c.Request().Context()

	query := c.QueryParam("q")
	if query == "" {
		// Redirect to blog index if no query
		return c.Redirect(http.StatusFound, "/blog")
	}

	// Search articles
	results, err := h.queries.SearchArticles(ctx, query)
	if err != nil {
		results = []repository.SearchArticlesRow{}
	}

	// Load cover images
	imagesMap := make(map[int32]repository.StaticFile)
	for _, article := range results {
		if article.CoverImageID.Valid {
			img, err := h.queries.GetStaticFile(ctx, article.CoverImageID.Int32)
			if err == nil {
				imagesMap[article.CoverImageID.Int32] = img
			}
		}
	}

	// SEO meta
	meta := model.DefaultPageMeta(
		"Buscar: "+query,
		"Resultados de búsqueda para '"+query+"' en el blog de JR del Perú.",
	)
	baseURL := getBaseURL()
	meta.Canonical = baseURL + "/blog/buscar?q=" + query

	return Render(c, http.StatusOK, view.BlogSearchPage(meta, query, results, imagesMap))
}

// HandleBlogArticle muestra un artículo individual del blog
func (h *Handler) HandleBlogArticle(c echo.Context) error {
	ctx := c.Request().Context()

	slug := c.Param("slug")
	if slug == "" {
		return c.Redirect(http.StatusFound, "/blog")
	}

	// Get article
	article, err := h.queries.GetPublishedArticleBySlug(ctx, slug)
	if err != nil {
		return c.Redirect(http.StatusFound, "/blog")
	}

	// Check if user is admin
	adminSession := h.getAdminSession(c)

	// Get article FAQs
	faqs, err := h.queries.ListArticleFAQs(ctx, article.ArticleID)
	if err != nil {
		faqs = []repository.ArticleFaq{}
	}

	// Get cover image
	var coverImage *repository.StaticFile
	if article.CoverImageID.Valid {
		img, err := h.queries.GetStaticFile(ctx, article.CoverImageID.Int32)
		if err == nil {
			coverImage = &img
		}
	}

	// Get comments with admin info
	comments, err := h.queries.ListArticleCommentsWithAdmin(ctx, article.ArticleID)
	if err != nil {
		comments = []repository.ListArticleCommentsWithAdminRow{}
	}

	// Build replies map with admin info
	repliesMap := make(map[int32][]repository.ListCommentRepliesWithAdminRow)
	for _, comment := range comments {
		replies, _ := h.queries.ListCommentRepliesWithAdmin(ctx, pgtype.Int4{Int32: comment.CommentID, Valid: true})
		repliesMap[comment.CommentID] = replies
	}

	// Get reCAPTCHA site key
	recaptchaSiteKey := h.recaptchaService.GetSiteKey()

	// Build SEO meta
	baseURL := getBaseURL()
	meta := model.PageMeta{
		Title:       article.Title,
		Description: article.Summary,
		OGType:      "article",
		Author:      article.Author,
		Canonical:   baseURL + "/blog/" + article.Slug,
	}

	if article.PublishedAt.Valid {
		meta.PublishedAt = article.PublishedAt.Time.Format(time.RFC3339)
	}

	if coverImage != nil {
		meta.OGImage = baseURL + path.Join(config.IMAGES_PATH, coverImage.FileName)
	}

	return Render(c, http.StatusOK, view.BlogArticlePage(meta, article, faqs, coverImage, comments, repliesMap, recaptchaSiteKey, adminSession))
}

// HandleFAQPage muestra la página pública de preguntas frecuentes
func (h *Handler) HandleFAQPage(c echo.Context) error {
	ctx := c.Request().Context()

	// Get visible global FAQs
	faqs, err := h.queries.ListVisibleGlobalFAQs(ctx)
	if err != nil {
		faqs = []repository.GlobalFaq{}
	}

	// Group by category
	faqsByCategory := make(map[string][]repository.GlobalFaq)
	var categories []string
	categorySet := make(map[string]bool)

	for _, faq := range faqs {
		cat := "General"
		if faq.Category.Valid && faq.Category.String != "" {
			cat = faq.Category.String
		}
		if !categorySet[cat] {
			categorySet[cat] = true
			categories = append(categories, cat)
		}
		faqsByCategory[cat] = append(faqsByCategory[cat], faq)
	}

	// SEO meta
	meta := model.DefaultPageMeta(
		"Preguntas Frecuentes",
		"Respuestas a las preguntas más frecuentes sobre nuestros productos y servicios de vidrio, aluminio y uPVC.",
	)
	baseURL := getBaseURL()
	meta.Canonical = baseURL + "/preguntas-frecuentes"

	return Render(c, http.StatusOK, view.FAQPage(meta, categories, faqsByCategory))
}

// getBaseURL returns the base URL for the site
func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://jrdelperu.com"
	}
	return baseURL
}
