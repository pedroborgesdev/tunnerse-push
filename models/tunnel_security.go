package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"math/big"
)

type PublicKeyJSON struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type EncryptedPayload struct {
	Nonce string `json:"nonce"`
	Body  string `json:"ciphertext"`
}

// ToECDSA converte PublicKeyJSON em *ecdsa.PublicKey
func (pk *PublicKeyJSON) ToECDSA() (*ecdsa.PublicKey, error) {
	xBytes, err := base64.StdEncoding.DecodeString(pk.X)
	if err != nil {
		return nil, errors.New("invalid X base64")
	}

	yBytes, err := base64.StdEncoding.DecodeString(pk.Y)
	if err != nil {
		return nil, errors.New("invalid Y base64")
	}

	pub := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}

	if !pub.Curve.IsOnCurve(pub.X, pub.Y) {
		return nil, errors.New("point is not on curve")
	}

	return pub, nil
}
