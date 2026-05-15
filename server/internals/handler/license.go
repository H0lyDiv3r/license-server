package handler

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"license-server/internals/domain"
	"license-server/internals/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const publicKey = "0ddc979bbf017e8627161321a0e193d90865d119a8fe87e62854aa32c4db017b"

type LicenseHandler struct {
	LicenseService service.LicenseService
	UserService    service.UserService
}

func NewLicenseHandler(licenseService service.LicenseService, userService service.UserService) LicenseHandler {
	return LicenseHandler{UserService: userService, LicenseService: licenseService}
}

func (h *LicenseHandler) GenerateLicense(w http.ResponseWriter, r *http.Request) {

	var req domain.GenerateKeyRequest
	json.NewDecoder(r.Body).Decode(&req)

	license, err := h.LicenseService.GenerateKey(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Errorf("error generating key: %s ", err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(license)
}

func (h *LicenseHandler) ActivateLicense(w http.ResponseWriter, r *http.Request) {

	var req domain.ActivateLicenseRequest
	json.NewDecoder(r.Body).Decode(&req)

	license, err := h.LicenseService.ActivateLicense(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to generate a token: %s", err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(license)

}

func (h *LicenseHandler) DecodeLicense(w http.ResponseWriter, r *http.Request) {
	var req domain.DecodeRequest
	json.NewDecoder(r.Body).Decode(&req)

	publickey, err := hex.DecodeString(publicKey)
	if err != nil {
		fmt.Println("cant read pubkey", err)
		return
	}

	pubKey := ed25519.PublicKey(publickey)
	token, err := jwt.ParseWithClaims(req.License, &domain.LicenseClaims{}, func(t *jwt.Token) (interface{}, error) {
		// explicitly reject any algorithm that isn't EdDSA
		// this prevents algorithm switching attacks
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	},
		// we handle expiry ourselves using trusted time
		// so we tell jwt not to validate it
		jwt.WithoutClaimsValidation(),
	)

	if err != nil {
		fmt.Println("issue here", err)
	}

	claims, ok := token.Claims.(*domain.LicenseClaims)
	if !ok {
		fmt.Println("issue here", err)
	}

	usr, _ := h.UserService.GetUser(r.Context(), int(claims.UserID))

	fmt.Println("claims returned", usr)

	w.Write([]byte("done"))

}
