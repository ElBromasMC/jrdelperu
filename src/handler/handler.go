package handler

import (
	"alc/repository"
	"alc/service"
)

type Handler struct {
	authService      *service.SessionAuthService
	fileService      *service.FileService
	recaptchaService *service.RecaptchaService
	emailService     *service.EmailService
	pdfService       *service.PDFService
	queries          *repository.Queries
}

func New(authService *service.SessionAuthService, fileService *service.FileService, recaptchaService *service.RecaptchaService, emailService *service.EmailService, pdfService *service.PDFService, queries *repository.Queries) Handler {
	return Handler{
		authService:      authService,
		fileService:      fileService,
		recaptchaService: recaptchaService,
		emailService:     emailService,
		pdfService:       pdfService,
		queries:          queries,
	}
}
