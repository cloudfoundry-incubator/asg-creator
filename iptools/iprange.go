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
	return bytes.Compare(r.Start, ip) == 0
}

func (r *IPRange) EndsAt(ip net.IP) bool {
	return bytes.Compare(r.End, ip) == 0
}

func (r *IPRange) SliceIP(ip net.IP) []IPRange {
	if !r.Contains(ip) {
		return []IPRange{*r}
	}

	if r.StartsAt(ip) {
		return []IPRange{
			{
				Start: Inc(ip),
				End:   r.End,
			},
		}
	}

	if r.EndsAt(ip) {
		return []IPRange{
			{
				Start: r.Start,
				End:   Dec(ip),
			},
		}
	}

	a := IPRange{
		Start: r.Start,
		End:   Dec(ip),
	}

	b := IPRange{
		Start: Inc(ip),
		End:   r.End,
	}

	if bytes.Compare(r.Start.To4(), Dec(ip).To4()) == 0 {
		a = IPRange{
			Start: r.Start,
		}
	}

	if bytes.Compare(r.End.To4(), Inc(ip).To4()) == 0 {
		b = IPRange{
			Start: r.Start,
		}
	}

	return []IPRange{
		a,
		b,
	}
}
