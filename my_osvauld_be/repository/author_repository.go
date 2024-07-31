package repository

import (
	"context"
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/google/uuid"
)

func GetAuthor(ctx context.Context, id uuid.UUID) (db.Author, error) {
	return database.Store.GetAuthor(ctx, id)
}

func CreateAuthor(ctx context.Context, args db.CreateAuthorParams) (db.Author, error) {
	return database.Store.CreateAuthor(ctx, args)
}
