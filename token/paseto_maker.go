package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKet []byte
}

func NewPasetoMaker(symmetricKet string) (Maker, error) {
	if len(symmetricKet) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKet: []byte(symmetricKet),
	}
	return maker, nil
}

func (maker PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	return maker.paseto.Encrypt(maker.symmetricKet, payload, nil)
}

func (maker PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKet, payload, nil)
	if err != nil {
		return nil, ErrIvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
