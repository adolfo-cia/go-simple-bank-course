package token

import (
	"testing"
	"time"

	"github.com/adolfo-cia/go-simple-bank-course/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	pMaker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	experiedAt := issuedAt.Add(duration)

	token, err := pMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, experiedAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	pMaker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := pMaker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	pMaker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := pMaker.CreateToken(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	token = token + "a"

	payload, err := pMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
