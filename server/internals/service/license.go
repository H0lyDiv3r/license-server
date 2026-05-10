package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"license-server/internals/domain"
	"license-server/internals/repository"
	"license-server/pkgs/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type LicenseService struct {
	repository repository.LicenseRepository
}

func NewLicenseService(repo repository.LicenseRepository) LicenseService {
	return LicenseService{repository: repo}
}

func (s *LicenseService) GenerateKey(ctx context.Context) (*domain.License, error) {

	b := make([]byte, 10)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	encoded = strings.ToUpper(encoded)

	key := encoded[0:4] + "-" + encoded[4:8] + "-" + encoded[8:12] + "-" + encoded[12:16]

	license, err := s.repository.StoreKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (s *LicenseService) ActivateLicense(ctx context.Context, req domain.ActivateLicenseRequest) (string, error) {

	fmt.Println("this is the request just incase", req)
	license, err := s.repository.GetLicenseByKey(ctx, req.LicenseKey)
	if err != nil {
		return "", fmt.Errorf("Key Not Found %s", err.Error())
	}

	token, err := IssueToken(*license, req.FingerPrint)
	if err != nil {
		return "", fmt.Errorf("failed to generate license token: %s", err.Error())
	}

	return token, nil
}

func IssueToken(license domain.License, fingerprint string) (string, error) {
	privateKey, err := utils.LoadPrivateKey()

	if err != nil {
		return "", err
	}
	claims := domain.LicenseClaims{
		LicenseID: license.ID,
		UserID:    license.UserID,
		MachineID: fingerprint,
		Status:    string(license.Status),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(license.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(license.ExpiresAt),
			NotBefore: jwt.NewNumericDate(license.IssuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(privateKey)

}
