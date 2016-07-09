package iptools

import (
	"bytes"
	"fmt"
	"net"
)

type IPRange struct {
	Start net.IP
	End   net.IP
}

func NewIPRangeFromIPNet(ipNet *net.IPNet) IPRange {
	min, max := NetworkRange(ipNet)
	return IPRange{
		Start: min,
		End:   max,
	}
}

func (r *IPRange) String() string {
	if r.End == nil {
		return r.Start.String()
	}

	return fmt.Sprintf("%s-%s", r.Start, r.End)
}

func (r *IPRange) Contains(ip net.IP) bool {
	return bytes.Compare(r.Start.To4(), ip.To4()) <= 0 && bytes.Compare(ip.To4(), r.End.To4()) <= 0
}

func (r *IPRange) OverlapsNet(ipNet *net.IPNet) bool {
	min, max := NetworkRange(ipNet)
	return ipNet.Contains(r.Start) || ipNet.Contains(r.End) ||
		r.Contains(min) || r.Contains(max)
}

func (r *IPRange) EqualsNet(ipNet *net.IPNet) bool {
	min, max := NetworkRange(ipNet)
	return r.Start.Equal(min) && r.End.Equal(max)
}

func (r *IPRange) OverlapsRange(other IPRange) bool {
	return other.Contains(r.Start) || other.Contains(r.End) ||
		r.Contains(other.Start) || r.Contains(other.End)
}

func (r *IPRange) SliceIPs(ips []net.IP) []IPRange {
	ipRanges := []IPRange{*r}

	for i := range ips {
		ipRanges = SliceIPFromRanges(ipRanges, ips[i])
	}

	return ipRanges
}

func (r *IPRange) StartsAt(ip net.IP) bool {
	return r.Start.Equal(ip)
}

func (r *IPRange) EndsAt(ip net.IP) bool {
	return r.End.Equal(ip)
}

func (r *IPRange) SliceIP(ip net.IP) []IPRange {
	if !r.Contains(ip) {
		return []IPRange{*r}
	}

	if r.StartsAt(ip) {
		if r.End.Equal(Inc(ip)) {
			return []IPRange{
				{
					Start: r.End,
				},
			}
		}

		return []IPRange{
			{
				Start: Inc(ip),
				End:   r.End,
			},
		}
	}

	if r.EndsAt(ip) {
		if r.Start.Equal(Dec(ip)) {
			return []IPRange{
				{
					Start: r.Start,
				},
			}
		}

		return []IPRange{
			{
				Start: r.Start,
				End:   Dec(ip),
			},
		}
	}

	return []IPRange{
		IPRange{
			Start: r.Start,
			End:   Dec(ip),
		},
		IPRange{
			Start: Inc(ip),
			End:   r.End,
		},
	}
}
