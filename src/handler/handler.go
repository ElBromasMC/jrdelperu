package handler

import (
	"alc/repository"
	"alc/service"
)

type Handler struct {
	authService *service.SessionAuthService
	fileService *service.FileService
	queries     *repository.Queries
}

func New(authService *service.SessionAuthService, fileService *service.FileService, queries *repository.Queries) Handler {
	return Handler{
		authService: authService,
		fileService: fileService,
		queries:     queries,
	}
}
