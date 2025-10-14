package item

import "github.com/wb-go/wbf/dbpg"

type Postgres struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Postgres {
	return &Postgres{db: db}
}
