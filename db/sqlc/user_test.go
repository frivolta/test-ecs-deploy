package db

import (
	"birdie/util"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Email:    util.RandomEmail(),
		Role:     RoleTEACHER,
		FullName: util.RandomName(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Role, user.Role)
	require.Equal(t, arg.Email, user.Email)
	require.NotEmpty(t, user.ID)
	require.NotZero(t, user.UpdatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Role, user2.Role)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Email)
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user1.Email)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	args := UpdateUserParams{
		Email: user.Email,
		Role:  RolePARENT,
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, updatedUser.Role, args.Role)
	require.WithinDuration(t, updatedUser.UpdatedAt, time.Now(), time.Second)

	updatedUser, err = testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: "fake_email@email.com",
		Role:  RoleTEACHER,
	})
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, updatedUser)
}
