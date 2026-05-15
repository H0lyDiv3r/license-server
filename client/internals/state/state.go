package state

type AppState struct {
	AuthToken    string `json:"authToken"`
	ValidLicense bool   `json:"validLicense"`
}
