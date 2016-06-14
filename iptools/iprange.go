package iptools

import (
	"bytes"
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

func (r *IPRange) Contains(ip net.IP) bool {
	return bytes.Compare(r.Start.To4(), ip.To4()) == -1 && bytes.Compare(ip.To4(), r.End.To4()) == -1
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

func (r *IPRange) SliceIP(ip net.IP) []IPRange {
	oneLess := Dec(ip)
	oneMore := Inc(ip)

	if bytes.Compare(r.Start, ip) == 0 {
		return []IPRange{
			{
				Start: oneMore,
				End:   r.End,
			},
		}
	}

	if bytes.Compare(r.End, ip) == 0 {
		return []IPRange{
			{
				Start: r.Start,
				End:   oneLess,
			},
		}
	}

	return []IPRange{
		{
			Start: r.Start,
			End:   net.IP(oneLess),
		},
		{
			Start: net.IP(oneMore),
			End:   r.End,
		},
	}
}
