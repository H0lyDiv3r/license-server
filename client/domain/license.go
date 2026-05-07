package domain

import "github.com/golang-jwt/jwt/v5"

type LicenseClaims struct {
	LicenseID uint   `json:"license_id"`
	UserID    uint   `json:"user_id"`
	MachineID string `json:"machine_id"`
	Status    string `json:"status"`
	jwt.RegisteredClaims
}
