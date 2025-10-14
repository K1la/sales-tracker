package dto

import "time"

// Item dto

type CreateItem struct {
	Type     string  `json:"type" binding:"required,oneof=income expense"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Date     string  `json:"date" binding:"required,datetime=2006-01-02"`
	Category string  `json:"category" binding:"required"`
}

type ItemResponse struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Date      string    `json:"date"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateItem struct {
	Type     *string  `json:"type" binding:"omitempty,oneof=income expense"`
	Amount   *float64 `json:"amount" binding:"omitempty,gt=0"`
	Date     *string  `json:"date" binding:"omitempty,datetime=2006-01-02"`
	Category *string  `json:"category" binding:"omitempty"`
}

type GetItemsParams struct {
	From   string   `form:"from"`
	To     string   `form:"to"`
	SortBy []string `form:"sort_by"`
}

// Analytics dto

type AnalyticsQuery struct {
	From string `form:"from" binding:"required,datetime=2006-01-02"`
	To   string `form:"to" binding:"required,datetime=2006-01-02"`
}

type AnalyticsResponse struct {
	Sum        float64 `json:"sum"`
	Average    float64 `json:"average"`
	Count      float64 `json:"count"`
	Median     float64 `json:"median"`
	Percentile float64 `json:"percentile"`
}
