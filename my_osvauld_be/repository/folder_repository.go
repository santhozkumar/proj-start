package repository

import (
	"context"
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/google/uuid"
)

func CreateFolder(ctx context.Context, args db.CreateFolderParams) (db.CreateFolderRow, error) {
	return database.Store.CreateFolder(ctx, args)
}
func GetFoldersForUser(ctx context.Context, userID uuid.UUID) ([]db.Folder, error) {
	return database.Store.GetFoldersForUser(ctx, uuid.NullUUID{UUID: userID, Valid: true})
}
