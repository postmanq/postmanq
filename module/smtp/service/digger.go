package service

type DiggerStatus int

const (
	DiggerStatusProcess DiggerStatus = iota
	DiggerStatusComplete
	DiggerStatusInvalid
)

type Digger interface {
	GetStatus() DiggerStatus
}

type digger struct {
}
