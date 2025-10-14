package items

import (
	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
)

type Postgres struct {
	db  *dbpg.DB
	log zerolog.Logger
}

func New(db *dbpg.DB, l zerolog.Logger) *Postgres {
	return &Postgres{db: db, log: l}
}
