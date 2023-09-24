package db

import (
	"context"
	"testing"
	"time"

	"github.com/Foody-App-Tech/Main-server/helpers"
	"github.com/stretchr/testify/require"
)

func TestGenerateSession(t *testing.T) {

	user := createRandomUser(t)

	arg := GenerateSessionParams{UserIDSession: user.UserID, Valid: true, UserAgent: helpers.RandomString(10)}
	session, err := testQueries.GenerateSession(context.Background(), arg)

	// basic check
	require.NoError(t, err)
	require.NotEmpty(t, session)

	// update session
	time.Sleep(time.Second * 2)

	session2, err := testQueries.GenerateSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.NotEqual(t, session2.UpdatedAt, session.UpdatedAt)
	require.Equal(t, session2.CreatedAt, session.CreatedAt)

}
