package analytics

import (
	"context"
	"github.com/K1la/sales-tracker/internal/dto"
)

func (s *Service) GetAnalytics(ctx context.Context, query dto.AnalyticsQuery) (*dto.AnalyticsResponse, error) {
	data, err := s.db.GetAnalytics(ctx, query)
	if err != nil {
		return nil, err
	}

	resp := &dto.AnalyticsResponse{
		Sum:        data.Sum,
		Average:    data.Average,
		Count:      float64(data.Count),
		Median:     data.Median,
		Percentile: data.Percentile,
	}

	return resp, nil
}
