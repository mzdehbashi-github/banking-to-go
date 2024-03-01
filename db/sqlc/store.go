package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

func (s *Store) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, beginErr := s.db.BeginTx(ctx, nil)

	if beginErr != nil {
		log.Fatal("error beginning transaction", beginErr)
	}

	// Rollback changes at the end of the test
	defer tx.Rollback()

	queries := New(tx)

	// Execute the function with the transaction
	err := fn(queries)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("transaction err: %v rollback error: %v", err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {

	var result TransferTxResult

	err := s.execTransaction(ctx, func(queries *Queries) error {
		transfer, createTransferErr := queries.CreateTransfer(
			ctx,
			CreateTransferParams(arg),
		)

		if createTransferErr != nil {
			log.Println("error in creating transfer", createTransferErr)
			return createTransferErr
		}

		fromEntry, createFromEntryErr := queries.CreateEntry(
			ctx,
			CreateEntryParams{
				AccountID: arg.FromAccountID,
				Amount:    -arg.Amount,
			},
		)
		if createFromEntryErr != nil {
			log.Println("error in creating entry", createFromEntryErr)
			return createFromEntryErr
		}

		toEntry, createToEntryErr := queries.CreateEntry(
			ctx,
			CreateEntryParams{
				AccountID: arg.ToAccountID,
				Amount:    arg.Amount,
			},
		)

		if createToEntryErr != nil {
			log.Println("error in creating entry", createToEntryErr)
			return createToEntryErr
		}

		accounts, updateAccountsErr := queries.UpdateAccountsBalance(
			ctx,
			UpdateAccountsBalanceParams{
				Amount:        arg.Amount,
				Fromaccountid: arg.FromAccountID,
				Toaccountid:   arg.ToAccountID,
			},
		)

		if updateAccountsErr != nil {
			return updateAccountsErr
		}

		log.Println("accounts, ", accounts)
		if len(accounts) != 2 {
			errMessage := fmt.Sprintf("error in updating accounts %d", len(accounts))
			log.Println(errMessage)
			return fmt.Errorf(errMessage)
		}

		log.Printf("%v %v", arg.FromAccountID, arg.ToAccountID)
		if accounts[0].ID == arg.FromAccountID {
			result.FromAccount = Account(accounts[0])
			result.ToAccount = Account(accounts[1])
		} else {
			result.FromAccount = Account(accounts[1])
			result.ToAccount = Account(accounts[0])
		}

		result.Transfer = transfer
		result.FromEntry = fromEntry
		result.ToEntry = toEntry
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
