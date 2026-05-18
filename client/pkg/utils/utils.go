package utils

import (
	"client/domain"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateHmac(secret string, payload string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(payload))
	signature := h.Sum(nil)

	hmacString := hex.EncodeToString(signature)

	return hmacString
}

func WriteJournalEntry(secret string, timestamp string) error {

	journalPath, err := GetJournalPath()
	if err != nil {
		return err
	}

	entry := domain.JournalEntry{
		LastSeen: timestamp,
		Hmac:     GenerateHmac(secret, timestamp),
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

func GetJournalPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to create journal file: %w", err)
	}

	dir := filepath.Join(home, ".config", "secure_desktop")
	journalPath := filepath.Join(dir, "journal.json")

	return journalPath, nil
}
