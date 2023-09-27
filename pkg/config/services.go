package config

import "go.uber.org/config"

type Option config.YAMLOption

func Static(val interface{}) Option {
	return config.Static(val)
}

func File(name string) Option {
	return config.File(name)
}

type ProviderFactory interface {
	CreateFromMap(map[string]string) (Provider, error)
}

type Provider interface {
	Populate(interface{}) error
	PopulateByKey(string, interface{}) error
}
