package entity

type MailServerStatus int

const (
	MailServerStatusScannable MailServerStatus = iota
	MailServerStatusScanned
	MailServerStatusInvalid
)

type MailServer struct {
	MXs    []MX
	Status MailServerStatus
}
