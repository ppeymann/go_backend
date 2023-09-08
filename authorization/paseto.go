package authorization

import (
	"errors"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var ErrInvalidToken = errors.New("provided token is not valid")

// PasetoMaker is PASETO token maker.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker and initialize properties
func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: symmetric key must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (m *PasetoMaker) CreateToken(claims *Claims) (string, error) {
	return m.paseto.Encrypt(m.symmetricKey, claims, nil)
}

// VerifyToken checks if the token string is valid or not and parse it to Claims struct
func (m *PasetoMaker) VerifyToken(token string) (*Claims, error) {
	claims := &Claims{}

	err := m.paseto.Decrypt(token, m.symmetricKey, claims, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
