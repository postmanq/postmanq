package entity

import "net"

type MX struct {
	*net.MX
	IP net.IP
}
