package store

import (
	"github.com/tidwall/buntdb"
)

type Config struct {
	DBFile string `yaml:"db_file"`
}

type Datastore struct {
	db *buntdb.DB
}

func New(cfg Config) (*Datastore, error) {
	dbFile := cfg.DBFile
	if dbFile == "" {
		dbFile = ":memory:"
	}
	db, err := buntdb.Open(dbFile)
	if err != nil {
		return nil, err
	}
	must(db.CreateIndex(IndexCmtCreatedAt, "cmt:*:model", buntdb.IndexJSONCaseSensitive("createdAt")))
	must(db.CreateIndex(IndexIssueLastChangedAt, "is:*:lastChangedAt", buntdb.IndexString))
	st := &Datastore{
		db: db,
	}
	return st, nil
}
