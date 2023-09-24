package db

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLCStore provides all functions to execute db queries and transactions ----> IMPORTANT : we implement store to re-use neccesary code for transactions
type Store interface {
	Querier
	CreateNewBusinessTx(ctx context.Context, arg CreateNewBusinessTxParams) (CreateNewBusinessTxResult, error)
}

type SQLCStore struct {
	*Queries
	db *sql.DB
}

// create new SQLCStore with db reference
func NewStore(db *sql.DB) *SQLCStore {
	return &SQLCStore{db: db, Queries: New(db)}
}

// execTX executes a function that implements transaction in the database
func (store *SQLCStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type CreateNewBusinessTxParams struct {
	CreateBusinessParams
	AddBusinessMemberParams
}

type CreateNewBusinessTxResult struct {
	BusinessId int64 `json:"business_id"`
}

// TransferTx
func (store *SQLCStore) CreateNewBusinessTx(ctx context.Context, arg CreateNewBusinessTxParams) (CreateNewBusinessTxResult, error) {
	var result CreateNewBusinessTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		newBusiness, err := q.CreateBusiness(ctx, arg.CreateBusinessParams)
		if err != nil {
			return err
		}
		result.BusinessId = newBusiness.BusinessID

		_, err = q.AddBusinessMember(ctx, arg.AddBusinessMemberParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
