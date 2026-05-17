package license

import (
	"bytes"
	"client/domain"
	"client/internals/state"
	"client/internals/store"
	"client/pkg/utils"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/golang-jwt/jwt/v5"
)

type License struct {
	ctx   context.Context
	state *state.AppState
	store *store.Store
}

const (
	publicKey = "0ddc979bbf017e8627161321a0e193d90865d119a8fe87e62854aa32c4db017b"
)

func NewLicense(state *state.AppState, store *store.Store) *License {
	return &License{state: state, store: store}
}

func (l *License) Startup(ctx context.Context) {
	l.ctx = ctx
}

func (l *License) GenerateLicense(Duration domain.GenerateKeyRequest) (*domain.License, error) {
	r, err := json.Marshal(Duration)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "http://localhost:3000/license/generate", bytes.NewReader(r))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+l.state.AuthToken)

	res, err := http.DefaultClient.Do(req)
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
		return nil, fmt.Errorf("signup failed with status %s", res.Status)
	}

	var response domain.License
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to read response body %s", err)
	}

	l.store.StoreLicence(l.ctx, response)
	return &response, nil
}

func (l *License) ActivateLicense(key string) error {

	machineID, err := machineid.ProtectedID("secure_desktop")
	if err != nil {
		return fmt.Errorf("cant read machine id", err)
	}

	payload, err := json.Marshal(domain.ActivateLicenseRequest{LicenseKey: key, FingerPrint: machineID})
	if err != nil {
		return fmt.Errorf("cant read request data", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/license/activate", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+l.state.AuthToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var errResp map[string]string
		_ = json.NewDecoder(res.Body).Decode(&errResp)
		if msg := errResp["error"]; msg != "" {
			return errors.New(msg)
		}
		return fmt.Errorf("activation failed with status %s", err)
	}

	var response struct {
		License domain.License
		Token   string `json:"token"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to parse response body", err)
	}

	// create journal file
	err = l.InitializeJournal()
	if err != nil {
		return fmt.Errorf("failed to initialize journal: %w", err)
	}

	fmt.Println("token generated", response)
	_, err = l.store.UpdateLicense(l.ctx, domain.License{MachineID: &machineID, LicenseString: response.Token, Status: "active", Key: key})
	if err != nil {
		fmt.Errorf("updating issue: %w", err)
	}
	return nil
}

func (l *License) DecodeLicense(license domain.License) (*domain.LicenseClaims, error) {

	pkey, err := hex.DecodeString(publicKey)
	if err != nil {
		l.state.ValidLicense = false
		return nil, fmt.Errorf("malformed public Key: %w", err)
	}
	pubKey := ed25519.PublicKey(pkey)
	token, err := jwt.ParseWithClaims(license.LicenseString, &domain.LicenseClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			l.state.ValidLicense = false
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	},
		jwt.WithoutClaimsValidation(),
	)

	if err != nil {
		l.state.ValidLicense = false
		return nil, fmt.Errorf("malformed public Key: %w", err)
	}

	claims, ok := token.Claims.(*domain.LicenseClaims)
	if !ok {
		l.state.ValidLicense = false
		return nil, fmt.Errorf("cant parse claims: %w", err)
	}
	fmt.Println("i made it here dor some reason", claims)
	l.state.ValidLicense = true
	return claims, nil
}

func (l *License) WriteJournalEntry(secret string, timeStamp int64) {
	// hmacString := utils.GenerateHmac(hmacSecret, string(timeStamp))
	fmt.Println("time now,", time.Now().Unix())
}

func (l *License) InitializeJournal() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to create journal file: %w", err)
	}

	dir := filepath.Join(home, ".config", "secure_desktop")
	journalPath := filepath.Join(dir, "journal.json")

	lastSeen := time.Now().Unix()
	lastSeenStr := strconv.FormatInt(lastSeen, 10)
	entry := domain.JournalEntry{
		LastSeen: string(lastSeenStr),
		Hmac:     utils.GenerateHmac(domain.HmacSecret, string(lastSeenStr)),
	}

	JournalEntry, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to compose journal entry: %w", err)
	}

	if err := os.WriteFile(journalPath, JournalEntry, 0600); err != nil {
		return fmt.Errorf("failed to write journal entry: %w", err)
	}

	return nil
}
