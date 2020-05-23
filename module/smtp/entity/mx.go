package entity

import "net"

type MX struct {
	*net.MX
	IP net.IP
}

type MXs []MX

func (s MXs) Equal(other []MX) bool {
	if len(s) != len(other) {
		return false
	}

	for i, item := range s {
		otherItem := other[i]

		if item.MX != otherItem.MX || !item.IP.Equal(otherItem.IP) {
			return false
		}
	}

	return true
}
