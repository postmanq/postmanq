package config

import "go.uber.org/config"

type Option config.YAMLOption

type LookupFunc = config.LookupFunc

func Static(val interface{}) Option {
	return config.Static(val)
}

func File(name string) Option {
	return config.File(name)
}

func Expand(lookup LookupFunc) Option {
	return config.Expand(lookup)
}

type ProviderFactory interface {
	Create(options ...Option) (Provider, error)
}

type Provider interface {
	Populate(interface{}) error
	PopulateByKey(string, interface{}) error
}
