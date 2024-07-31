package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (store *SQLStore) UpdateAuthorName(ctx context.Context, id uuid.UUID, name string) (err error) {
	err = store.exec_tx(ctx, func(q *Queries) error {
		var err error
		r, err := q.GetAuthor(ctx, id)
		if err != nil {
			return err
		}
		if err := q.UpdateAuthor(ctx, UpdateAuthorParams{
			ID:   r.ID,
			Name: name,
		}); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (store *SQLStore) exec_tx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
	}
	qtx := store.Queries.WithTx(tx)
	err = fn(qtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %s; rollback err: %s", err, rbErr)
		}
	}
	return tx.Commit()
}
