package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LicenseStatus = string

const (
	statusActive  LicenseStatus = "active"
	statusRevoked LicenseStatus = "revoked"
	statusPending LicenseStatus = "pending"
)

type License struct {
	ID        uint          `gorm:"primaryKey"`
	Key       string        `gorm:"uniqueIndex; not null"`
	UserID    uint          `gorm:"uniqueIndex; not null"`
	MachineID *string       `gorm:"uniqueIndex"`
	Status    LicenseStatus `gorm: "default:'pending'"`
	IssuedAt  time.Time
	ExpiresAt time.Time
	RenewedAt *time.Time
	User      User `gorm:"foreignKey:UserID"`
}

type ActivateLicenseRequest struct {
	FingerPrint string `json:"fingerPrint"`
	LicenseKey  string `json:"licenseKey"`
}

type LicenseClaims struct {
	LicenseID uint   `json:"license_id"`
	UserID    uint   `json:"user_id"`
	MachineID string `json:"machine_id"`
	Status    string `json:"status"`
	jwt.RegisteredClaims
}

type DecodeRequest struct {
	License string `json:"license"`
}

type GenerateKeyResponse struct {
	License License
	Key     string
}
