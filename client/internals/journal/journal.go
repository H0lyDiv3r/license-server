package journal

import (
	"client/domain"
	"client/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
)

type Journal struct {
	entry domain.JournalEntry
}

func New() *Journal {
	return &Journal{}
}

func (j *Journal) Read() error {
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
	producedHmac := utils.GenerateHmac(domain.HmacSecret, fmt.Sprintf("%s:%v", parsedEntry.LastSeen, parsedEntry.Tampered))
	if producedHmac != parsedEntry.Hmac {
		return fmt.Errorf("journal entry has been tampered with or is corrupted")
	}
	j.entry = parsedEntry
	return nil
}

func (j *Journal) Write(timestamp string, tampered bool) error {
	hmac := utils.GenerateHmac(domain.HmacSecret, fmt.Sprintf("%s:%v", timestamp, tampered))
	entry := domain.JournalEntry{
		LastSeen: timestamp,
		Tampered: tampered,
		Hmac:     hmac,
	}
	journalPath, err := utils.GetJournalPath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to compose journal entry: %w", err)
	}
	if err := os.WriteFile(journalPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write journal entry: %w", err)
	}
	j.entry = entry
	return nil
}

func (j *Journal) Entry() domain.JournalEntry {
	return j.entry
}
