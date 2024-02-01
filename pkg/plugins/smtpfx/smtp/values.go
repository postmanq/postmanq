package smtp

import "time"

type Config struct {
	Hostname     string         `yaml:"hostname"`
	IPs          []*string      `yaml:"ips"`
	DkimSelector string         `yaml:"dkimSelector"`
	TLS          *TLSConfig     `yaml:"tls"`
	Timeout      *TimeoutConfig `yaml:"timeout"`
}

type TLSConfig struct {
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"privateKey"`
}

type TimeoutConfig struct {
	Connection time.Duration `yaml:"connection"`
	Hello      time.Duration `yaml:"hello"`
	Mail       time.Duration `yaml:"mail"`
	Rcpt       time.Duration `yaml:"rcpt"`
	Data       time.Duration `yaml:"data"`
}

type MxRecord struct {
	Host      string
	Priority  uint16
	CreatedAt time.Time
}
