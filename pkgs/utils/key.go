package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LicenseClaims struct {
	LicenseID   string `json:"license_id"`
	Email       string `json:"email"`
	Plan        string `json:"plan"`
	Fingerprint string `json:"fingerprint"`
	jwt.RegisteredClaims
}

type License struct {
	ID         string
	Email      string
	Plan       string
	Status     string
	MaxDevices int
	ExpiresAt  time.Time
}

func LoadPrivateKey() (ed25519.PrivateKey, error) {
	pemStr := os.Getenv("KEY_GEN_PRIVATE")
	if pemStr == "" {
		return nil, fmt.Errorf("KEY_GEN_PRIVATE not set")
	}

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	var ok bool
	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not Ed25519")
	}

	return privateKey, nil
}
