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
	LicenseString string        `gorm:"uniqueIndex; not null"`
	MachineID     string        `gorm:"uniqueIndex"`
	Status        LicenseStatus `gorm: "default:'pending'"`
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
