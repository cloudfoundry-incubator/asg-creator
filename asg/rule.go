package asg

import (
	"bytes"
	"net"
	"strings"
)

type Rule struct {
	Protocol    string `json:"protocol"`
	Destination string `json:"destination"`
	Ports       string `json:"ports,omitempty"`
	Type        string `json:"type,omitempty"`
	Code        string `json:"code,omitempty"`
	Log         bool   `json:"log,omitempty"`
}

func (r Rule) Contains(ipString string) bool {
	ip := net.ParseIP(ipString)

	dip := net.ParseIP(r.Destination)
	if dip != nil {
		return dip.Equal(ip)
	}

	dip, dipNet, err := net.ParseCIDR(r.Destination)
	if err == nil {
		return dipNet.Contains(ip)
	}

	dips := strings.Split(r.Destination, "-")
	minDip := net.ParseIP(dips[0])
	maxDip := net.ParseIP(dips[1])

	return bytes.Compare(ip, minDip) >= 0 && bytes.Compare(ip, maxDip) <= 0
}
