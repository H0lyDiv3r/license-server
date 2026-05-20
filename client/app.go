package main

import (
	"client/domain"
	"client/internals/journal"
	"client/internals/license"
	"client/internals/state"
	"client/internals/store"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/denisbrodbeck/machineid"
)

const publicKey = "0ddc979bbf017e8627161321a0e193d90865d119a8fe87e62854aa32c4db017b"

// App struct
type App struct {
	ctx       context.Context
	state     *state.AppState
	store     *store.Store
	license   *license.License
	startTime time.Time
	lastSeen  time.Time
	journal   *journal.Journal
}

// NewApp creates a new App application struct
func NewApp(state *state.AppState, store *store.Store, license *license.License, journal *journal.Journal) *App {
	return &App{state: state, store: store, license: license, journal: journal}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.startTime = time.Now()
	err := a.CheckLicense()
	if err != nil {
		fmt.Println("license has been messed with so user needs to renew")
		a.state.ValidLicense = false
	}
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

	err := a.journal.Read()
	if err != nil {
		return err
	}

	lastSeenDateInt, err := strconv.ParseInt(a.journal.Entry().LastSeen, 10, 64)
	if err == nil {
		a.lastSeen = time.Unix(lastSeenDateInt, 0)
	}

	if a.journal.Entry().Tampered {
		return fmt.Errorf("license has been tampered with. you need to renew license")
	}

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

	err = a.CheckDeviceFingerPrint(parsedLicense)
	if err != nil {

		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		a.journal.Write(timeStr, true)
		return err
	}

	err = a.CheckLicenseTime(parsedLicense)
	if err != nil {

		fmt.Println("HERE WE ARE AN ERROR. TIME HAS BEEN ROLLED BACK: %w", err)
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		err := a.journal.Write(timeStr, true)
		if err != nil {
			fmt.Println("there is an error writing journal")
		}
		return err
	}

	return nil
}

func (a *App) CheckDeviceFingerPrint(parsedLicense *domain.LicenseClaims) error {

	machineID, err := machineid.ProtectedID("secure_desktop")
	if err != nil || machineID != parsedLicense.MachineID {
		return fmt.Errorf("machine finger print doesnt match the one stored in license: %w", err)
	}
	return nil
}

func (a *App) CheckLicenseTime(parsedLicense *domain.LicenseClaims) error {

	lastSeenDateInt, err := strconv.ParseInt(a.journal.Entry().LastSeen, 10, 64)
	if err != nil {
		return fmt.Errorf("last seen value is not readable: %w", err)
	}

	// if journal entry is not corrupt compare the last seen date with os date and take the latest of the two
	lastSeenDate := time.Unix(int64(lastSeenDateInt), 0)
	a.lastSeen = lastSeenDate
	if time.Now().Before(lastSeenDate) {
		return fmt.Errorf("time rollback detected. user clock and journal entry dont add up")
	}

	// if healthy: check expiry. if expired send to auth and force license renewal
	// check if expiry date is bigger than date now, expiry
	if time.Now().After(parsedLicense.ExpiresAt.Time) {
		return fmt.Errorf("license is expired. you need to renew it.")
	}

	// check if time now is bigger than issuedAt; if it is smaller it means user rolled back time.
	if time.Now().Before(parsedLicense.IssuedAt.Time) {
		return fmt.Errorf("your system clock and the app clock are out of sync")
	}

	return nil
}

func (a *App) OnShutDown() {

	var trustedTime time.Time
	// in the case of time rollback
	// check if journal time is ahead of current startTime
	// if time has been rolled back use journal time, else use os time
	// store truested time. ie startTime + elapsed time so that time can only go foreward

	if a.startTime.Before(a.lastSeen) {
		trustedTime = a.lastSeen.Add(time.Since(a.startTime))
	} else {
		trustedTime = a.startTime.Add(time.Since(a.startTime))
	}

	lastSeen := trustedTime.Unix()
	lastSeenStr := strconv.FormatInt(lastSeen, 10)

	err := a.journal.Write(lastSeenStr, a.journal.Entry().Tampered)
	if err != nil {
		log.Println("shutting down with grace")
	}
}
