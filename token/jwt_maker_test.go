package token

import (
	"m1thrandir225/your_time/util"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)


func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))

	require.NoError(t, err)

	email := util.RandomEmail()
	duration := time.Minute

	//To compare the token later
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)


	token, err := maker.CreateToken(email, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)

	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)	
}



func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))

	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomEmail(), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(util.RandomEmail(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
	
}