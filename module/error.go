package module

type ErrorCode int

const (
	ErrorCodeUnknown = 0
)

type Error struct {
	Code    ErrorCode
	Message string
}
