package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute DB queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, errRb)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParam contains the input parameters of the transfer transaction
type TransferTxParam struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to another
//
// It creates a transfer record, add both accounts entries and update both accounts balance within a single db transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey) // get the tx name for debug logging
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: arg.FromAccountId, Amount: -arg.Amount})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: arg.ToAccountId, Amount: arg.Amount})
		if err != nil {
			return err
		}

		//avoid deadlocks on concurrent update of Account
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, addMoneyParams{
				account1Id:     arg.FromAccountId,
				account1Amount: -arg.Amount,
				account2Id:     arg.ToAccountId,
				account2Amount: arg.Amount})
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, addMoneyParams{
				account1Id:     arg.ToAccountId,
				account1Amount: arg.Amount,
				account2Id:     arg.FromAccountId,
				account2Amount: -arg.Amount})
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

type addMoneyParams struct {
	account1Id     int64
	account1Amount int64
	account2Id     int64
	account2Amount int64
}

func addMoney(ctx context.Context, q *Queries, arg addMoneyParams) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     arg.account1Id,
		Amount: arg.account1Amount,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     arg.account2Id,
		Amount: arg.account2Amount,
	})
	if err != nil {
		return
	}
	return
}
