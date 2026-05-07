package repository

import (
	"context"
	"fmt"
	"license-server/internals/domain"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(context.Context, domain.SignupRequest) (domain.User, error)
	GetById(context.Context, int) (domain.User, error)
	GetByEmail(context.Context, domain.SigninRequest) (domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, req domain.SignupRequest) (domain.User, error) {

	usr := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	err := gorm.G[domain.User](u.db).Create(ctx, &usr)

	if err != nil {
		log.Println("Failed to create user")
		return domain.User{}, fmt.Errorf("failed to create user", err.Error())
	}

	return usr, nil
}

func (u *userRepository) GetById(ctx context.Context, id int) (domain.User, error) {

	usr, err := gorm.G[domain.User](u.db).Where("id=?", id).First(ctx)
	if err != nil {
		return domain.User{}, nil
	}
	return usr, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, req domain.SigninRequest) (domain.User, error) {
	usr, err := gorm.G[domain.User](u.db).Where("email=?", req.Email).First(ctx)
	if err != nil {
		return domain.User{}, nil
	}
	return usr, nil
}
