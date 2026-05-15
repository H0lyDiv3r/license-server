package main

import (
	"client/internals/license"
	"client/internals/state"
	"client/internals/store"
	"context"
	"fmt"
	"time"
)

const publicKey = "0ddc979bbf017e8627161321a0e193d90865d119a8fe87e62854aa32c4db017b"

// App struct
type App struct {
	ctx     context.Context
	state   *state.AppState
	store   *store.Store
	license *license.License
}

// NewApp creates a new App application struct
func NewApp(state *state.AppState, store *store.Store, license *license.License) *App {
	return &App{state: state, store: store, license: license}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.CheckLicense()
}

func (a *App) IsLoggedIn() bool {
	return a.state.AuthToken != ""
}

func (a *App) GetState() *state.AppState {
	return &state.AppState{
		AuthToken:    a.state.AuthToken,
		ValidLicense: a.state.ValidLicense,
	}
}

func (a *App) CheckLicense() error {
	// fetch license
	license, err := a.store.GetLicense(a.ctx)

	// if the license doesnt exist; send user to login and force them to buy a subscription
	if err != nil || license == nil {
		return fmt.Errorf("License not found: %w", err)
	}

	// if license exists
	// decode license and check for health and tampering;
	parsedLicense, err := a.license.DecodeLicense(*license)
	if err != nil {
		return fmt.Errorf("license cant be decoded: %w", err)
	}
	fmt.Println("here is license", license)
	fmt.Println("claims", parsedLicense)
	fmt.Println("this is the time now", time.Now().Before(license.ExpiresAt))

	// we are still trusting the user at this point
	// if healthy: check expiry. if expired send to auth and force license renewal
	// check if expiry date is bigger than date now, expiry
	// check if time now is bigger than issuedAt; if it is smaller it means user rolled back time.

	if time.Now().After(parsedLicense.ExpiresAt.Time) || time.Now().Before(parsedLicense.IssuedAt.Time) {

		fmt.Println("this license is so expired man", time.Now().After(parsedLicense.ExpiresAt.Time))
		fmt.Println("somebody rolled back time", time.Now().Before(parsedLicense.IssuedAt.Time))
		a.state.ValidLicense = false
	}

	return nil
}
