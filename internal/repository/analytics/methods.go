package analytics

import (
	"context"
	"database/sql"
	"errors"
	"github.com/K1la/sales-tracker/internal/dto"
	"github.com/K1la/sales-tracker/internal/model"
)

func (r *Postgres) GetAnalytics(ctx context.Context, q dto.AnalyticsQuery) (*model.Aggregated, error) {
	query := `
	SELECT
		COALESCE(SUM(amount), 0),
		COALESCE(AVG(amount), 0),
		COUNT(*),
		COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0),
		COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0)
	FROM items
	WHERE date BETWEEN $1 AND $2;
	`

	var agg model.Aggregated
	err := r.db.QueryRowContext(ctx, query, q.From, q.To).Scan(
		&agg.Sum, &agg.Average, &agg.Count, &agg.Median, &agg.Percentile,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return &agg, nil
}
