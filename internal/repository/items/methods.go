package items

import (
	"context"
	"database/sql"
	"errors"
	"github.com/K1la/sales-tracker/internal/dto"
	"github.com/K1la/sales-tracker/internal/model"
	"time"
)

var (
	ErrItemNotFound = errors.New("no item found")
)

func (r *Postgres) Create(ctx context.Context, req *model.Item) error {
	query := `
	INSERT INTO items
	(type, amount, date, category, description)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		req.Type,
		req.Amount,
		req.Date,
		req.Category,
		req.Description,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *Postgres) GetAll(ctx context.Context, params dto.GetItemsParams) ([]dto.ItemResponse, error) {
	baseQuery := `
	SELECT *
	FROM items
	`

	args := []interface{}{}
	conditiond := ""
	if params.From != "" && params.To != "" {
		conditiond = " WHERE date BETWEEN $1 AND $2"
		args = append(args, params.From, params.To)
	}

	orderBy := " ORDER BY date DESC"
	if len(params.SortBy) > 0 {
		orderBy = " ORDER By " + params.SortBy[0]
	}

	query := baseQuery + orderBy + conditiond
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.ItemResponse
	for rows.Next() {
		var it dto.ItemResponse
		var date time.Time
		if err = rows.Scan(
			&it.ID,
			&it.Type,
			&it.Amount,
			&date,
			&it.Category,
			&it.Description,
			&it.CreatedAt,
			&it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		it.Date = date.Format(time.DateOnly)
		items = append(items, it)
	}
	return items, nil
}

func (r *Postgres) GetByID(ctx context.Context, id string) (*dto.ItemResponse, error) {
	query := `
	SELECT * 
	FROM items
	WHERE id = $1
	`

	var it dto.ItemResponse
	var date time.Time
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&it.ID,
		&it.Type,
		&it.Amount,
		&date,
		&it.Category,
		&it.Description,
		&it.CreatedAt,
		&it.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	it.Date = date.Format(time.DateOnly)
	return &it, nil
}

func (r *Postgres) Update(ctx context.Context, req *dto.ItemResponse) error {
	query := `
	UPDATE items
	SET type = $1,
	    amount = $2,
	    date = $3,
	    category = $4,
	    description = $5
	WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		req.Type,
		req.Amount,
		req.Date,
		req.Category,
		req.Description,
		req.ID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrItemNotFound
	}
	return err
}

func (r *Postgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM items WHERE id = $1`, id)
	return err
}
