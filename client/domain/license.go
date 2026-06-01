package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LicenseClaims struct {
	LicenseID uint   `json:"license_id"`
	UserID    uint   `json:"user_id"`
	MachineID string `json:"machine_id"`
	Status    string `json:"status"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type License struct {
	ID            uint          `gorm:"primaryKey"`
	Key           string        `json:"key" gorm:"uniqueIndex; not null"`
	LicenseString string        `gorm:"uniqueIndex; not null"`
	MachineID     *string       `gorm:"uniqueIndex"`
	Status        LicenseStatus `gorm: "default:'pending'"`
	Duration      uint
	IssuedAt      time.Time
	ExpiresAt     time.Time
	RenewedAt     *time.Time
}

type LicenseStatus = string

const (
	statusActive  LicenseStatus = "active"
	statusRevoked LicenseStatus = "revoked"
	statusPending LicenseStatus = "pending"
)

type GenerateKeyResponse struct {
	License License
	Key     string
}

type ActivateLicenseRequest struct {
	FingerPrint string `json:"fingerPrint"`
	LicenseKey  string `json:"licenseKey"`
}

type GenerateKeyRequest struct {
	Duration uint `json:"duration"`
}
