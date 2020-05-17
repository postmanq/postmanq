package module

type Email struct {
	Sender        string
	Recipient     string
	RecipientHost string
	Body          []byte
}
