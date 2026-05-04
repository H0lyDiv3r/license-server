package domain

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex; not null"`
	Password  string `gorm:"not null"`
	createdAt time.Time
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPayload struct {
	UserId float64 `json:"user_id"`
	Email  string  `json:"email"`
}
