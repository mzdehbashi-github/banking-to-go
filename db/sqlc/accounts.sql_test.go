package db

import (
	"context"
	"database/sql"
	"gopsql/banking/util"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func accountFactory(queries Querier, optionalParams ...CreateAccountParams) Account {
	var arg CreateAccountParams

	if len(optionalParams) > 0 {
		arg = optionalParams[0]
	} else {
		user := userFactory(queries)
		arg = CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomMoney(),
			Currency: Currency(util.RandomCurrency()),
		}
	}

	account, err := queries.CreateAccount(context.Background(), arg)
	if err != nil {
		log.Fatal("error in creating account", err)
	}

	return account
}

func TestCreateAccount(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {
			user := userFactory(queries)
			arg := CreateAccountParams{
				Owner:    user.Username,
				Balance:  util.RandomMoney(),
				Currency: Currency(util.RandomCurrency()),
			}

			account, err := queries.CreateAccount(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, account)
			require.Equal(t, account.Owner, arg.Owner)
			require.Equal(t, account.Balance, arg.Balance)
			require.Equal(t, account.Currency, arg.Currency)

			require.NotZero(t, account.ID)
			require.NotZero(t, account.CreatedAt)
		},
	)
}

func TestGetAccount(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {
			account := accountFactory(queries)
			fetchedAccount, err := queries.GetAccount(context.Background(), account.ID)
			require.NoError(t, err)
			require.Equal(t, account.ID, fetchedAccount.ID)
		},
	)
}

func TestUpdateAccount(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {
			account := accountFactory(queries)
			updateArg := UpdateAccountParams{ID: account.ID, Balance: 255}
			updatedAccount, err := queries.UpdateAccount(context.Background(), updateArg)
			require.NoError(t, err)
			require.Equal(t, updatedAccount.Balance, updateArg.Balance)
		},
	)
}

func TestDeleteAccount(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {
			account := accountFactory(queries)
			err := queries.DeleteAccound(context.Background(), account.ID)
			require.NoError(t, err)
			fetchedAccount, getAccountErr := queries.GetAccount(context.Background(), account.ID)
			require.Error(t, getAccountErr)
			require.EqualError(t, getAccountErr, sql.ErrNoRows.Error())
			require.Zero(t, fetchedAccount.ID)
		},
	)
}

func TestListAccounts(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {

			for i := 0; i < 10; i++ {
				accountFactory(queries)
			}

			listArg := ListAccountsParams{Limit: 4, Offset: 5}
			accounts, err := queries.ListAccounts(context.Background(), listArg)
			require.NoError(t, err)
			require.Len(t, accounts, 4)
		},
	)
}
