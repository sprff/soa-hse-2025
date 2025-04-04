package psql

import (
	"fmt"
	"postservice/internal/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var _ storage.CommonRepository = &PsqlStorage{}

type PsqlStorage struct {
	db *sqlx.DB
}

func NewPsqlStorage() (*PsqlStorage, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		"user", "password", "postservice-database", 5432, "db"))
	if err != nil {
		return nil, fmt.Errorf("can't connect: %w", err)
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS posts
(
    id VARCHAR(63) PRIMARY KEY,
    title TEXT,
    description TEXT,
	creator_id VARCHAR(63),
	is_private BOOLEAN,
	tags TEXT[],
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);
`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare tables: %w", err)
	}
	return &PsqlStorage{db: db}, nil
}
