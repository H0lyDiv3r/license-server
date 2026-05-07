package license

import (
	"client/domain"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type License struct {
	ctx context.Context
}

const publicKey = "0ddc979bbf017e8627161321a0e193d90865d119a8fe87e62854aa32c4db017b"

func NewLicense() *License {
	return &License{}
}

func (l *License) Startup(ctx context.Context) {
	l.ctx = ctx
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
