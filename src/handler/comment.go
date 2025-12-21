package handler

import (
	"alc/repository"
	"alc/service"
	"alc/view"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// getAdminSession returns the admin session data if logged in, nil otherwise
func (h *Handler) getAdminSession(c echo.Context) *service.SessionData {
	session, err := h.authService.GetSessionStore().Get(c.Request(), service.SessionName)
	if err != nil {
		return nil
	}
	sessionData, err := service.GetSessionData(session)
	if err != nil {
		return nil
	}
	return sessionData
}

// HandleCommentCreate handles public comment submission with reCAPTCHA validation
func (h *Handler) HandleCommentCreate(c echo.Context) error {
	ctx := c.Request().Context()

	// Get article by slug
	slug := c.Param("slug")
	article, err := h.queries.GetPublishedArticleBySlug(ctx, slug)
	if err != nil {
		return c.String(http.StatusNotFound, "Artículo no encontrado")
	}

	// Check if user is admin
	adminSession := h.getAdminSession(c)

	// Verify reCAPTCHA token (skip for admins)
	if adminSession == nil {
		recaptchaToken := c.FormValue("g-recaptcha-response")
		valid, err := h.recaptchaService.VerifyToken(recaptchaToken)
		if err != nil || !valid {
			if errors.Is(err, service.ErrRecaptchaNotFound) {
				return c.String(http.StatusBadRequest, "Por favor, completa el captcha")
			}
			if errors.Is(err, service.ErrRecaptchaLowScore) {
				return c.String(http.StatusBadRequest, "Verificación sospechosa. Intenta de nuevo.")
			}
			// Log the actual error for debugging
			c.Logger().Errorf("reCAPTCHA error: %v", err)
			return c.String(http.StatusBadRequest, "Verificación de captcha fallida. Intenta de nuevo.")
		}
	}

	// Get form values - for admins, use session data
	var authorName, authorEmail string
	if adminSession != nil {
		authorName = adminSession.Username
		authorEmail = adminSession.Email
	} else {
		authorName = c.FormValue("author_name")
		authorEmail = c.FormValue("author_email")
	}
	content := c.FormValue("content")
	parentIDStr := c.FormValue("parent_id")

	// Validate required fields
	if authorName == "" || content == "" {
		return c.String(http.StatusBadRequest, "Nombre y comentario son requeridos")
	}

	// Parse parent_id if present (for replies)
	var parentID pgtype.Int4
	if parentIDStr != "" {
		id, err := strconv.ParseInt(parentIDStr, 10, 32)
		if err == nil {
			parentID = pgtype.Int4{Int32: int32(id), Valid: true}
		}
	}

	// Build admin_id if logged in as admin
	var adminID pgtype.Int4
	if adminSession != nil {
		adminID = pgtype.Int4{Int32: adminSession.AdminID, Valid: true}
	}

	// Create comment
	createdComment, err := h.queries.CreateArticleComment(ctx, repository.CreateArticleCommentParams{
		ArticleID:   article.ArticleID,
		ParentID:    parentID,
		AuthorName:  authorName,
		AuthorEmail: pgtype.Text{String: authorEmail, Valid: authorEmail != ""},
		Content:     content,
		AdminID:     adminID,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al guardar el comentario")
	}

	// If this is a reply, fetch with admin info and return the reply card
	if parentID.Valid {
		reply, err := h.queries.GetArticleCommentWithAdmin(ctx, createdComment.CommentID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error al obtener el comentario")
		}
		// Convert to ListCommentRepliesWithAdminRow type
		replyRow := repository.ListCommentRepliesWithAdminRow{
			CommentID:     reply.CommentID,
			ArticleID:     reply.ArticleID,
			ParentID:      reply.ParentID,
			AuthorName:    reply.AuthorName,
			AuthorEmail:   reply.AuthorEmail,
			Content:       reply.Content,
			AdminID:       reply.AdminID,
			CreatedAt:     reply.CreatedAt,
			AdminUsername: reply.AdminUsername,
		}
		return Render(c, http.StatusOK, view.CommentReplyCard(replyRow, slug, adminSession))
	}

	// Fetch the comment with admin info
	comment, err := h.queries.GetArticleCommentWithAdmin(ctx, createdComment.CommentID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener el comentario")
	}

	// Convert to ListArticleCommentsWithAdminRow type
	commentRow := repository.ListArticleCommentsWithAdminRow{
		CommentID:     comment.CommentID,
		ArticleID:     comment.ArticleID,
		ParentID:      comment.ParentID,
		AuthorName:    comment.AuthorName,
		AuthorEmail:   comment.AuthorEmail,
		Content:       comment.Content,
		AdminID:       comment.AdminID,
		CreatedAt:     comment.CreatedAt,
		AdminUsername: comment.AdminUsername,
	}

	// Get reCAPTCHA site key for reply forms in the new comment
	recaptchaSiteKey := h.recaptchaService.GetSiteKey()
	return Render(c, http.StatusOK, view.CommentCard(commentRow, nil, slug, recaptchaSiteKey, adminSession))
}

// HandleAdminReplyCreate handles admin reply to a comment
func (h *Handler) HandleAdminReplyCreate(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	// Parse article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Parse comment ID
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de comentario inválido")
	}

	// Get reply content
	content := c.FormValue("content")
	if content == "" {
		return c.String(http.StatusBadRequest, "La respuesta no puede estar vacía")
	}

	// Create admin reply
	reply, err := h.queries.CreateAdminReply(ctx, repository.CreateAdminReplyParams{
		ArticleID:  int32(articleID),
		ParentID:   pgtype.Int4{Int32: int32(commentID), Valid: true},
		AuthorName: sessionData.Username,
		Content:    content,
		AdminID:    pgtype.Int4{Int32: sessionData.AdminID, Valid: true},
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al guardar la respuesta")
	}

	// Get the parent comment to return the updated card
	parentComment, err := h.queries.GetArticleComment(ctx, int32(commentID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al obtener el comentario")
	}

	// Get all replies for this comment
	replies, err := h.queries.ListCommentReplies(ctx, pgtype.Int4{Int32: int32(commentID), Valid: true})
	if err != nil {
		replies = []repository.ArticleComment{reply}
	}

	// Return the updated comment card with replies
	return Render(c, http.StatusOK, view.AdminCommentCard(int32(articleID), parentComment, replies))
}

// HandlePublicCommentDelete handles comment deletion from public blog pages (admin only)
func (h *Handler) HandlePublicCommentDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Check if user is admin
	adminSession := h.getAdminSession(c)
	if adminSession == nil {
		return c.String(http.StatusUnauthorized, "No autorizado")
	}

	// Parse comment ID
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de comentario inválido")
	}

	// Delete comment (CASCADE will delete replies)
	if err := h.queries.DeleteArticleComment(ctx, int32(commentID)); err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar el comentario")
	}

	// Return empty response (HTMX will remove the element)
	return c.NoContent(http.StatusOK)
}

// HandleCommentDelete handles comment deletion by admin
func (h *Handler) HandleCommentDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse article ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de artículo inválido")
	}

	// Parse comment ID
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de comentario inválido")
	}

	// Delete comment (CASCADE will delete replies)
	if err := h.queries.DeleteArticleComment(ctx, int32(commentID)); err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar el comentario")
	}

	// Return updated comments list
	comments, err := h.queries.ListArticleComments(ctx, int32(articleID))
	if err != nil {
		comments = []repository.ArticleComment{}
	}

	// Build replies map
	repliesMap := make(map[int32][]repository.ArticleComment)
	for _, comment := range comments {
		replies, _ := h.queries.ListCommentReplies(ctx, pgtype.Int4{Int32: comment.CommentID, Valid: true})
		repliesMap[comment.CommentID] = replies
	}

	return Render(c, http.StatusOK, view.AdminCommentsGrid(int32(articleID), comments, repliesMap))
}
