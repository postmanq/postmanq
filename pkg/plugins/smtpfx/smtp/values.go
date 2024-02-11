package smtp

import (
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"sync"
	"time"
)

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

type EmailAddress struct {
	Address   string
	LocalPart string
	Domain    string
}

type RecipientDescriptor struct {
	Servers    collection.Slice[*ServerDescriptor]
	ModifiedAt time.Time
}

type ServerDescriptor struct {
	MxRecord             MxRecord
	Clients              collection.Slice[Client]
	ModifiedAt           time.Time
	hasMaxCountOfClients bool
	mtx                  sync.Mutex
}

func (d *ServerDescriptor) HasMaxCountOfClients() bool {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	return d.hasMaxCountOfClients
}

func (d *ServerDescriptor) SetMaxCountOfClientsOn() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.hasMaxCountOfClients = true
}

func (d *ServerDescriptor) SetMaxCountOfClientsOff() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.hasMaxCountOfClients = false
}

type ClientStatus int

const (
	ClientStatusUndefined ClientStatus = iota
	ClientStatusFree
	ClientStatusBusy
)
