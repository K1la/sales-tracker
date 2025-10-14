package analytics

import (
	"context"
	"github.com/K1la/sales-tracker/internal/dto"
	"github.com/K1la/sales-tracker/internal/model"
	"github.com/rs/zerolog"
)

type Service struct {
	db  Repo
	log zerolog.Logger
}

func New(d Repo, l zerolog.Logger) *Service {
	return &Service{db: d, log: l}
}

type Repo interface {
	GetAnalytics(ctx context.Context, q dto.AnalyticsQuery) (*model.Aggregated, error)
}
