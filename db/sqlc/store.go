package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) executeTransaction(ctx context.Context, fn func(*Queries) error) error {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(transaction)
	err = fn(q)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction Err: %v, rollback Err: %v", err, rollbackErr)
		}
	}
	return transaction.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTransactionParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error) {
	var result TransferTransactionResult

	err := store.executeTransaction(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	fromAmount int64,
	toAccountID int64,
	toAmount int64,
) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}

	toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAmount,
	})
	return
}
