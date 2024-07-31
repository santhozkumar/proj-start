package service

import (
	"context"
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/google/uuid"
)

func CreateFolder(ctx context.Context, userID uuid.UUID, folder dto.CreateFolderRequest) (db.CreateFolderRow, error) {
	return repository.CreateFolder(ctx, db.CreateFolderParams{
		Name:        folder.Name,
		Description: sql.NullString{String: folder.Description, Valid: true},
		CreatedBy:   uuid.NullUUID{UUID: userID, Valid: true},
	})
}

func GetFoldersForUser(ctx context.Context, userID uuid.UUID) ([]db.Folder, error) {
	return repository.GetFoldersForUser(ctx, userID)
}
