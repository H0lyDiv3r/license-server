package domain

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex; not null"`
	createdAt time.Time
}
