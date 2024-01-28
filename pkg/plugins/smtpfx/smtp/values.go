package smtp

import "time"

type MailEvent struct {
}

type MxRecord struct {
	Host      string
	Priority  uint16
	CreatedAt time.Time
}
