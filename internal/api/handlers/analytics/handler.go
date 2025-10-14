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

func New(s Service, l zerolog.Logger) *Handler {
	return &Handler{service: s, log: l}
}

type Service interface {
	GetAnalytics(ctx context.Context, q dto.AnalyticsQuery) (*dto.AnalyticsResponse, error)
}
