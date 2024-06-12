package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmeizh/simplebank/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomString(6)
	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, newFullName, user.FullName)
	require.Equal(t, oldUser.Email, user.Email)
	require.Equal(t, oldUser.HashedPassword, user.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, newEmail, user.Email)
	require.Equal(t, oldUser.FullName, user.FullName)
	require.Equal(t, oldUser.HashedPassword, user.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(10)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, newHashedPassword, user.HashedPassword)
	require.Equal(t, oldUser.FullName, user.FullName)
	require.Equal(t, oldUser.Email, user.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomString(6)
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(10)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, newHashedPassword, user.HashedPassword)
	require.Equal(t, newFullName, user.FullName)
	require.Equal(t, newEmail, user.Email)
}
