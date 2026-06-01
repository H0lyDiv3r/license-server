package domain

import "time"

type CreateCheckoutResponse struct {
	Url string `json:"url"`
}

type Payment struct {
	ID            uint      `gorm:"primaryKey"`
	SessionId     string    `gorm:"uniqueIndex; not null"`
	Amount        uint      `gorm:"not null"`
	Currency      string    `gorm:"not null"`
	Email         string    `gorm:"not null"`
	Status        string    `gorm:"not null"`
	PaymentMethod string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	MetaData      string
	License       License `gorm:"foreignKey:LicenseID"`
	LicenseID     uint    `gorm:"not null"`
	// User          User    `gorm:"foreignKey:UserID"`
}

type SavePaymentRequest struct {
	SessionId     string
	Amount        uint
	Currency      string
	Email         string
	Status        string
	PaymentMethod string
	Metadata      string
	LicenseId     uint
	// UserId        uint
}
