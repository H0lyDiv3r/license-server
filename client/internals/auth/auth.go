package auth

import (
	"bytes"
	"errors"
	"client/domain"
	"client/internals/state"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Auth struct {
	ctx   context.Context
	state *state.AppState
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Startup(ctx context.Context, state *state.AppState) {
	a.ctx = ctx
	a.state = state
}

func (a *Auth) Signin(req domain.LoginRequest) (*domain.LoginResponse, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := http.Post("http://localhost:3000/auth/signin", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errResp map[string]string
		_ = json.NewDecoder(res.Body).Decode(&errResp)
		if msg := errResp["error"]; msg != "" {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("signin failed with status %s", res.Status)
	}

	var result domain.LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	a.state.AuthToken = result.Token
	fmt.Println("token generated", result)

	return &result, nil
}

func (a *Auth) Signup(req domain.LoginRequest) (*domain.LoginResponse, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := http.Post("http://localhost:3000/auth/signup", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	fmt.Println("showing stuff", res)
	if res.StatusCode != http.StatusOK {
		var errResp map[string]string
		_ = json.NewDecoder(res.Body).Decode(&errResp)
		if msg := errResp["error"]; msg != "" {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("signup failed with status %s", res.Status)
	}

	r, err := a.Signin(req)
	if err != nil {
		return nil, err
	}

	result := &domain.LoginResponse{
		Token: r.Token,
	}

	return result, nil
}
