package db

import (
	"context"
	"m1thrandir225/your_time/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User  {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams {
		FirstName: util.RandomString(6),
		LastName: util.RandomString(6),
		Email: util.RandomString(6),
		Password: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)

	require.Equal(t, arg.LastName, user.LastName)

	require.Equal(t, arg.Email, user.Email)

	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.ID)

	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	user := createRandomUser(t)

	user2, err := testQueries.GetUserByID(context.Background(), user.ID)

	require.NoError(t, err)

	require.NotEmpty(t, user2)

	require.Equal(t, user.ID, user2.ID)

	require.Equal(t, user.FirstName, user2.FirstName)

	require.Equal(t, user.LastName, user2.LastName)

	require.Equal(t, user.Email, user2.Email)

	require.Equal(t, user.Password, user2.Password)

	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user.UpdatedAt, user2.UpdatedAt, time.Second)	
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user.Email)

	require.NoError(t, err)

	require.NotEmpty(t, user2)

	require.Equal(t, user.ID, user2.ID)

	require.Equal(t, user.FirstName, user2.FirstName)

	require.Equal(t, user.LastName, user2.LastName)

	require.Equal(t, user.Email, user2.Email)

	require.Equal(t, user.Password, user2.Password)

	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)

	require.WithinDuration(t, user.UpdatedAt, user2.UpdatedAt, time.Second)

}
