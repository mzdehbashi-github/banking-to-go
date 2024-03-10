package db

import (
	"context"
	"fmt"
	"gopsql/banking/util"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func userFactory(queries Querier, optionalParams ...CreateUserParams) User {
	var arg CreateUserParams

	if len(optionalParams) > 0 {
		arg = optionalParams[0]
	} else {
		username := util.RandomOwner()
		arg = CreateUserParams{
			Username:       username,
			HashedPassword: "123456",
			FullName:       "John Doe",
			Email:          fmt.Sprintf("%v@c.com", username),
		}
	}

	account, err := queries.CreateUser(context.Background(), arg)
	if err != nil {
		log.Fatal("error in creating account", err)
	}

	return account
}

func TestCreateUser(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {
			arg := CreateUserParams{
				Username:       "myusername",
				HashedPassword: "123456",
				FullName:       "John Doe",
				Email:          "johndoe@c.com",
			}

			user, err := queries.CreateUser(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, user)
			require.Equal(t, user.Username, arg.Username)
			require.Equal(t, user.FullName, arg.FullName)
			require.Equal(t, user.Email, arg.Email)

			require.NotZero(t, user.CreatedAt)
			require.NotZero(t, user.PasswordChangedAt)
		},
	)
}

func TestGetUser(t *testing.T) {
	withTransaction(
		t,
		func(queries Querier) {
			user := userFactory(queries)
			fetchedAccount, err := queries.GetUser(context.Background(), user.Username)
			require.NoError(t, err)
			require.Equal(t, user.Username, fetchedAccount.Username)
		},
	)
}
