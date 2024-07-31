package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Store interface {
	Querier
	UpdateAuthorName(ctx context.Context, id uuid.UUID, name string) (err error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(conn *sql.DB) Store {
	return &SQLStore{
		db:      conn,
		Queries: New(conn),
	}
}
