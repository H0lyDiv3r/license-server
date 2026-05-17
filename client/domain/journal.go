package domain

type JournalEntry struct {
	LastSeen string `json:"last_seen"`
	Hmac     string `json:"hmac"`
}

const HmacSecret = ` Like two ships in the night in foggy weathe
Just a-waitin' for fresh winds to blow
Maybe we're losin' one another
I could be wrong I don't know
Like two dopes in the boat, without a paddle
Just a-wonderin' why it don't go
We could be losin' one another
I could be wrong I don't know `
