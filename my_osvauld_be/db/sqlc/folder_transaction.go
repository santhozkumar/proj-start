package db

import (
	"context"
)


func (store *SQLStore) CreateFolderTransaction(ctx  context.Context, arg CreateFolderParams) (CreateFolderRow, error) {
    folderData, err := store.CreateFolder(ctx, arg)



    return
}
