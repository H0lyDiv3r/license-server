package domain

import "time"

type LicenseStatus = string

const (
	statusActive  LicenseStatus = "active"
	statusRevoked LicenseStatus = "revoked"
)

type License struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex; not null"`
	UserID    uint   `gorm:"uniqueIndex; not null"`
	MachineID string `gorm:"uniqueIndex; not null "`
	// Status    LicenseStatus `gorm: "default:'active'"`
	IssuedAt  time.Time
	ExpiresAt time.Time
	RenewedAt *time.Time
	User      User `gorm:"foreignKey:UserID"`
}
