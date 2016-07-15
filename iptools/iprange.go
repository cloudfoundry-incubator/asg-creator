package iptools

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type IPRange struct {
	Start net.IP
	End   net.IP
}

func (r *IPRange) UnmarshalYAML(tag string, value interface{}) error {
	switch data := value.(type) {
	case string:
		dataWithoutSpaces := strings.Replace(data, " ", "", -1)
		idx := strings.IndexAny(dataWithoutSpaces, "-")
		if idx == -1 {
			return fmt.Errorf("invalid range given (missing hyphen): '%v'", data)
		}

		*r = IPRange{
			Start: net.ParseIP(dataWithoutSpaces[:idx]),
			End:   net.ParseIP(dataWithoutSpaces[idx+1:]),
		}

	default:
		return fmt.Errorf("invalid range given: '%#v'", data)
	}

	return nil
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

func (r *IPRange) EqualsRange(other IPRange) bool {
	return r.Start.Equal(other.Start) && r.End.Equal(other.End)
}

func (r *IPRange) OverlapsRange(other IPRange) bool {
	return other.Contains(r.Start) || other.Contains(r.End) ||
		r.Contains(other.Start) || r.Contains(other.End)
}

func (r *IPRange) SliceIPs(ips []net.IP) []IPRange {
	rs := []IPRange{*r}

	for i := range ips {
		var newRanges []IPRange
		for j := range rs {
			newRanges = append(newRanges, rs[j].SliceIP(ips[i])...)
		}
		rs = newRanges
	}

	return rs
}

func (r *IPRange) SliceRanges(ipRanges []IPRange) []IPRange {
	rs := []IPRange{*r}

	for i := range ipRanges {
		var newRanges []IPRange
		for j := range rs {
			newRanges = append(newRanges, rs[j].SliceRange(ipRanges[i])...)
		}
		rs = newRanges
	}

	return rs
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

	if r.Start.Equal(Dec(ip)) {
		a = IPRange{
			Start: r.Start,
		}
	}

	if r.End.Equal(Inc(ip)) {
		b = IPRange{
			Start: r.Start,
		}
	}

	return []IPRange{
		a,
		b,
	}
}

func (r *IPRange) SliceRange(other IPRange) []IPRange {
	if !r.OverlapsRange(other) {
		return []IPRange{*r}
	}

	if r.EqualsRange(other) {
		return nil
	}

	thisStart := r.Start.To4()
	thisEnd := r.End.To4()

	otherStart := other.Start.To4()
	otherEnd := other.End.To4()

	var ipRanges []IPRange
	if bytes.Compare(thisStart, otherStart) == -1 {
		ipRanges = append(ipRanges, IPRange{
			Start: thisStart,
			End:   Dec(otherStart),
		})
	}

	if bytes.Compare(thisEnd, otherEnd) == 1 {
		ipRanges = append(ipRanges, IPRange{
			Start: Inc(otherEnd),
			End:   thisEnd,
		})
	}

	return ipRanges
}
