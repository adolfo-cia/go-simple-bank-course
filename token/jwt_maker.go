package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

type JWTPayload struct {
	Username string
	jwt.RegisteredClaims
}

func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken create a new token for a specific username and duration
func (jwtM *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	iat := time.Now()
	exp := iat.Add(duration)
	payload := JWTPayload{
		username,
		jwt.RegisteredClaims{
			Subject:   tokenId.String(),
			IssuedAt:  jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(jwtM.secretKey))
}

// VerifyToken checks if the token is valid or not
func (jwtM *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtM.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenUnverifiable) {
			return nil, ErrInvalidToken
		}

		return nil, err
	}

	jwtPayload, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	id, err := uuid.Parse(jwtPayload.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        id,
		Username:  jwtPayload.Username,
		IssuedAt:  jwtPayload.IssuedAt.Time,
		ExpiredAt: jwtPayload.ExpiresAt.Time,
	}, nil
}
