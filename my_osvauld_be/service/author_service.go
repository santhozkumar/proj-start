package service

import (
	"context"
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/google/uuid"
)

func GetAuthor(ctx context.Context, id uuid.UUID) (db.Author, error) {
	return repository.GetAuthor(ctx, id)
}

func CreateAuthor(ctx context.Context, user dto.CreateAuthor) (db.Author, error) {
	return repository.CreateAuthor(ctx, db.CreateAuthorParams{
		Name: user.Name,
		Bio:  sql.NullString{String: user.Bio, Valid: true},
	})
}
