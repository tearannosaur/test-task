package handlers

import (
	"app/repository"

	"go.uber.org/zap"
)

type Handler struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewHandler(r *repository.Repository, l *zap.Logger) *Handler {
	return &Handler{
		repo:   r,
		logger: l,
	}
}
