package types

type Constructor interface{}

type StartPlugin interface {
	Start() error
}

type StopPlugin interface {
	Stop() error
}
