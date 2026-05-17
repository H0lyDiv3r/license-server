package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHmac(secret string, payload string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(payload))
	signature := h.Sum(nil)

	hmacString := hex.EncodeToString(signature)

	return hmacString
}
