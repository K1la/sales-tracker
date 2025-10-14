package analytics

import (
	"context"
	"github.com/K1la/sales-tracker/internal/dto"
	"github.com/rs/zerolog"
)

type Handler struct {
	service Service
	log     zerolog.Logger
}

func New(s Service) *Handler {
	return &Handler{service: s}
}

type Service interface {
	GetAnalytics(ctx context.Context, query dto.AnalyticsQuery) (*dto.AnalyticsResponse, error)
}
