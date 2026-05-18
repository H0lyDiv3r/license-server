package main

import (
	"client/domain"
	"client/internals/license"
	"client/internals/state"
	"client/internals/store"
	"client/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/denisbrodbeck/machineid"
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

	err = a.CheckDeviceFingerPrint(parsedLicense)
	if err != nil {
		// emit an error
		a.state.ValidLicense = false
	}

	err = a.CheckLicenseTime(parsedLicense)
	if err != nil {
		// emit an error
		a.state.ValidLicense = false
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

	// at this point we shouldnt really trus the users os time. since it can easily be manipulated
	// we will check if time has been altered by checking last enty journal and os time. take the later CheckLicenseTime
	// open journal
	journalPath, err := utils.GetJournalPath()
	if err != nil {
		return err
	}

	file, err := os.ReadFile(journalPath)
	if err != nil {
		return fmt.Errorf("journal doesnt exist or is corrupted: %w", err)
	}

	var parsedEntry domain.JournalEntry
	if err := json.Unmarshal(file, &parsedEntry); err != nil {
		return fmt.Errorf("journal doesnt exist or is corrupted: %w", err)
	}
	lastSeenDateInt, err := strconv.ParseInt(parsedEntry.LastSeen, 10, 64)

	// read journal produce hmac for the lastSeen value and see if it has been tampered with
	producedHmac := utils.GenerateHmac(domain.HmacSecret, parsedEntry.LastSeen)
	if producedHmac != parsedEntry.Hmac || err != nil {
		return fmt.Errorf("journal entry has been tampered with or is corrupted")
	}

	// if journal entry is not corrupt compare the last seen date with os date and take the latest of the two
	lastSeenDate := time.Unix(int64(lastSeenDateInt), 0)
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
	lastSeen := time.Now().Unix()
	lastSeenStr := strconv.FormatInt(lastSeen, 10)

	err := utils.WriteJournalEntry(domain.HmacSecret, lastSeenStr)
	if err != nil {
		log.Println("shutting down with grace")
	}
}
