package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalida key size: must be exact %d characters", chacha20poly1305.KeySize)
	}

	pMaker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return pMaker, nil
}

// CreateToken create a new token for a specific username and duration
func (pMaker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	iat := time.Now()
	exp := iat.Add(duration)
	jsonToken := paseto.JSONToken{
		Subject:    tokenId.String(),
		IssuedAt:   iat,
		Expiration: exp,
	}
	jsonToken.Set("username", username)

	return pMaker.paseto.Encrypt(pMaker.symmetricKey, jsonToken, nil)
}

// VerifyToken checks if the token is valid or not
func (pMaker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	jsonToken := &paseto.JSONToken{}
	err := pMaker.paseto.Decrypt(token, pMaker.symmetricKey, jsonToken, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = jsonToken.Validate()
	if err != nil {
		if strings.ContainsAny(errors.Unwrap(err).Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	id, err := uuid.Parse(jsonToken.Subject)
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        id,
		Username:  jsonToken.Get("username"),
		IssuedAt:  jsonToken.IssuedAt,
		ExpiredAt: jsonToken.Expiration,
	}, nil
}
