package server

type Server interface {
	Register(descriptor Descriptor) error
	Start() error
}
