package model

import "time"

type Item struct {
	ID        int64
	Type      string // income(доход) / expense(расход)
	Amount    float64
	Date      string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Aggregated struct {
	Sum        float64
	Average    float64
	Count      int64
	Median     float64
	Percentile float64
}
