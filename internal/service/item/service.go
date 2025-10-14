package item

import "context"

type Service struct {
	db Repo
}

func New(d Repo) *Service {
	return &Service{db: d}
}

type Repo interface {
	// TODO: дописать параметры
	CreateItem(ctx context.Context)
}
