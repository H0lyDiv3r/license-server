package license

import (
	"bytes"
	"client/domain"
	"client/internals/state"
	"client/internals/store"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/denisbrodbeck/machineid"
	"github.com/golang-jwt/jwt/v5"
)

type License struct {
	ctx   context.Context
	state *state.AppState
	store *store.Store
}

const publicKey = "0ddc979bbf017e8627161321a0e193d90865d119a8fe87e62854aa32c4db017b"

func NewLicense(state *state.AppState, store *store.Store) *License {
	return &License{state: state, store: store}
}

func (l *License) Startup(ctx context.Context) {
	l.ctx = ctx
}

func (l *License) GenerateLicense() (*domain.License, error) {
	req, err := http.NewRequest("POST", "http://localhost:3000/license/generate", nil)
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
		Token string `json:"token"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to parse response body", err)
	}

	fmt.Println("token generated", response)
	_, err = l.store.UpdateLicense(l.ctx, domain.License{MachineID: &machineID, LicenseString: response.Token, Status: "active", Key: key})
	if err != nil {
		fmt.Println("updating issue", err)
	}
	return nil
}

func (l *License) DecodeLicense(license string) *domain.LicenseClaims {

	pkey, err := hex.DecodeString(publicKey)
	if err != nil {
		fmt.Println("cant read pubkey", err)
		return &domain.LicenseClaims{}
	}
	pubKey := ed25519.PublicKey(pkey)
	token, err := jwt.ParseWithClaims(license, &domain.LicenseClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	},
		jwt.WithoutClaimsValidation(),
	)

	if err != nil {
		fmt.Println("issue here", err)
	}

	claims, ok := token.Claims.(*domain.LicenseClaims)
	if !ok {
		fmt.Println("another issue here", err)
	}
	fmt.Println("claims", claims)

	return claims
}

// func (h *LicenseHandler) DecodeLicense(w http.ResponseWriter, r *http.Request) {
// 	var req domain.DecodeRequest
// 	json.NewDecoder(r.Body).Decode(&req)

// 	publickey, err := hex.DecodeString(publicKey)

// 	pubKey := ed25519.PublicKey(publickey)
// 	token, err := jwt.ParseWithClaims(req.License, &domain.LicenseClaims{}, func(t *jwt.Token) (interface{}, error) {
// 		// explicitly reject any algorithm that isn't EdDSA
// 		// this prevents algorithm switching attacks
// 		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
// 		}
// 		return pubKey, nil
// 	},
// 		// we handle expiry ourselves using trusted time
// 		// so we tell jwt not to validate it
// 		jwt.WithoutClaimsValidation(),
// 	)

// 	if err != nil {
// 		fmt.Println("issue here", err)
// 	}

// 	claims, ok := token.Claims.(*domain.LicenseClaims)
// 	if !ok {
// 		fmt.Println("issue here", err)
// 	}

// 	usr, _ := h.UserService.GetUser(r.Context(), int(claims.UserID))

// 	fmt.Println("claims returned", usr)

// 	w.Write([]byte("done"))

// }
