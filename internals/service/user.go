package service

import (
	"context"
	"fmt"
	"license-server/internals/domain"
	"license-server/internals/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{repository: repo}
}

func (u *UserService) Signup(ctx context.Context, req *domain.SignupRequest) error {

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashed)
	user, err := u.repository.Create(ctx, *req)

	if err != nil {
		return err
	}

	fmt.Println("signed up", user)
	return nil
}

func (u *UserService) Signin(ctx context.Context, req *domain.SigninRequest) (string, error) {
	usr, err := u.repository.GetByEmail(ctx, *req)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("unauthorized", err.Error())
	}

	fmt.Println("signin")
	return generateJwt(usr)
}

func (u *UserService) GetUser(ctx context.Context, id int) (domain.User, error) {
	usr, err := u.repository.GetById(ctx, id)

	if err != nil {
		return domain.User{}, err
	}

	return usr, nil
}

func generateJwt(usr domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": usr.ID,
		"email":   usr.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
