package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/Foody-App-Tech/Main-server/helpers"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{SocialID: strconv.Itoa(helpers.RandomInt(0, 1000000)),
		Username: helpers.RandomUserName(), Email: helpers.RandomEmail(),
		Picture:  helpers.RandomPicture(),
		Provider: helpers.RandomProvider()}

	user, err := testQueries.CreateUser(context.Background(), arg)

	// basic check
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// 	// is this equal to arg provided
	require.Equal(t, arg.SocialID, user.SocialID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Picture, user.Picture)
	require.Equal(t, arg.Provider, user.Provider)

	// not zero
	require.NotZero(t, user.UserID)
	require.NotZero(t, user.RegisteredAt)
	return user
}

func TestUpdateUser(t *testing.T) {

	// case 1
	userData := createRandomUser(t)
	oldUserData := createRandomUser(t)

	arg := UpdateUserParams{Username: userData.Username, Email: userData.Email, Pictue: userData.Picture, UserID: userData.UserID, OldUsername: oldUserData.Username, OldEmail: oldUserData.Email, OldPicture: oldUserData.Picture}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	// check no err
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	// check that values are different
	require.NotEqual(t, updatedUser.Username, oldUserData.Username)
	require.NotEqual(t, updatedUser.Email, oldUserData.Email)
	require.NotEqual(t, updatedUser.Picture, oldUserData.Picture)

}
