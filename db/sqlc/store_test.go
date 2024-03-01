package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(conn)

	account1 := accountFactory(
		store.Queries,
		CreateAccountParams{
			Owner:    "John",
			Currency: CurrencyUSD,
			Balance:  int64(200),
		},
	)

	account2 := accountFactory(
		store.Queries,
		CreateAccountParams{
			Owner:    "Joe",
			Currency: CurrencyUSD,
			Balance:  int64(200),
		},
	)

	transferAmounts := []int{100, 90, 80, 70, 60, 50, 40, 30, 20, 10}
	resChan := make(chan *TransferTxResult, len(transferAmounts))
	errChan := make(chan error, len(transferAmounts))

	for _, amount := range transferAmounts {
		go func(trAmount int) {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        int64(trAmount),
			})

			resChan <- res
			errChan <- err
		}(amount)
	}

	// check results
	var totalTrasnferedAmount int64
	for range transferAmounts {
		err := <-errChan
		res := <-resChan

		// check result
		if err != nil {
			require.Empty(t, res)
		} else {
			// validate result
			require.NotEmpty(t, res)

			// check transfer
			transfer := res.Transfer
			require.NotEmpty(t, transfer)
			require.Equal(t, transfer.FromAccountID, account1.ID)
			require.Equal(t, transfer.ToAccountID, account2.ID)
			// require.Equal(t, transfer.Amount, amount)
			require.NotZero(t, transfer.ID)
			require.NotZero(t, transfer.CreatedAt)

			// check entries
			fromEntry := res.FromEntry
			require.NotZero(t, fromEntry.ID)
			_, fromEntryErr := store.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, fromEntryErr)
			toEntry := res.ToEntry
			require.NotZero(t, toEntry.ID)
			_, toEntryErr := store.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, toEntryErr)

			// check accounts
			fromAccount := res.FromAccount
			require.NotEmpty(t, fromAccount)
			require.Equal(t, fromAccount.ID, account1.ID)

			toAccount := res.ToAccount
			require.NotEmpty(t, toAccount)
			require.Equal(t, toAccount.ID, account2.ID)

			// diff1 := account1.Balance - fromAccount.Balance
			// diff2 := toAccount.Balance - account2.Balance
			// require.Equal(t, diff1, diff2)
			// require.True(t, diff1 > 0)
			// require.True(t, diff1%int64(amount) == 0)

			totalTrasnferedAmount += transfer.Amount
		}
	}

	updatedAccount1, getAccount1Err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, getAccount1Err)

	updatedAccount2, getAccount2Err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, getAccount2Err)

	log.Println(">> final: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, (account1.Balance - totalTrasnferedAmount), updatedAccount1.Balance)
	require.Equal(t, (account2.Balance + totalTrasnferedAmount), updatedAccount2.Balance)
}
