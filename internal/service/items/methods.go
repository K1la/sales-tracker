package items

import (
	"context"
	"fmt"
	"github.com/K1la/sales-tracker/internal/dto"
	"github.com/K1la/sales-tracker/internal/model"
	"time"
)

func (s *Service) Create(ctx context.Context, req dto.CreateItem) (*dto.ItemResponse, error) {
	parsedDate, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	item := model.Item{
		Type:        req.Type,
		Amount:      req.Amount,
		Date:        parsedDate,
		Category:    req.Category,
		Description: req.Description,
	}

	if err = s.db.Create(ctx, &item); err != nil {
		return nil, err
	}

	return &dto.ItemResponse{
		ID:          item.ID,
		Type:        item.Type,
		Amount:      item.Amount,
		Date:        item.Date.Format(time.DateOnly),
		Category:    item.Category,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (s *Service) GetAll(ctx context.Context, params dto.GetItemsParams) ([]dto.ItemResponse, error) {
	items, err := s.db.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ItemResponse, 0, len(items))
	for _, it := range items {
		res = append(res, dto.ItemResponse{
			ID:          it.ID,
			Type:        it.Type,
			Amount:      it.Amount,
			Date:        it.Date,
			Category:    it.Category,
			Description: it.Description,
			CreatedAt:   it.CreatedAt,
			UpdatedAt:   it.UpdatedAt,
		})
	}

	return res, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*dto.ItemResponse, error) {
	item, err := s.db.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ItemResponse{
		ID:          item.ID,
		Type:        item.Type,
		Amount:      item.Amount,
		Date:        item.Date,
		Category:    item.Category,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, id string, req dto.UpdateItem) (*dto.ItemResponse, error) {
	item, err := s.db.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Type != nil {
		item.Type = *req.Type
	}
	if req.Amount != nil {
		item.Amount = *req.Amount
	}
	if req.Date != nil {
		item.Date = *req.Date
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	item.UpdatedAt = time.Now()

	if err = s.db.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.db.Delete(ctx, id)
}
