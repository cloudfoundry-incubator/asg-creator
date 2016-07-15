package iptools

import (
	"bytes"
	"net"
)

func PrivateIPRanges() []IPRange {
	return []IPRange{
		{
			Start: net.IP{10, 0, 0, 0},
			End:   net.IP{10, 255, 255, 255},
		},
		{
			Start: net.IP{172, 16, 0, 0},
			End:   net.IP{172, 31, 255, 255},
		},
		{
			Start: net.IP{192, 168, 0, 0},
			End:   net.IP{192, 168, 255, 255},
		},
	}
}

func PublicIPRanges() []IPRange {
	nets := []net.IP{
		net.IP{0, 0, 0, 0},
	}

	for _, ipRange := range PrivateIPRanges() {
		nets = append(nets, Dec(ipRange.Start), Inc(ipRange.End))
	}

	nets = append(nets, net.IP{255, 255, 255, 255})

	ranges := []IPRange{}
	for i := 0; i < len(nets); i += 2 {
		ranges = append(ranges, IPRange{
			Start: nets[i],
			End:   nets[i+1],
		})
	}

	return ranges
}

func NetworkRange(ipNet *net.IPNet) (net.IP, net.IP) {
	// Inspired by https://github.com/docker/libnetwork/blob/master/netutils/utils.go
	ip4 := ipNet.IP.To4()
	min := ip4.Mask(ipNet.Mask)
	max := make([]byte, len(min))

	for i := range min {
		max[i] = ip4[i] | ^ipNet.Mask[i]
	}

	return min, max
}

func Dec(ip net.IP) net.IP {
	ipc := CopyIP(ip)
	for j := len(ipc) - 1; j >= 0; j-- {
		ipc[j]--
		if ipc[j] != 255 {
			break
		}
	}

	return ipc
}

func Inc(ip net.IP) net.IP {
	ipc := CopyIP(ip)
	// Hat tip to Russ Cox
	// https://groups.google.com/forum/#!topic/golang-nuts/zlcYA4qk-94
	for j := len(ipc) - 1; j >= 0; j-- {
		ipc[j]++
		if ipc[j] > 0 {
			break
		}
	}

	return ipc
}

func NetworkOverlaps(left, right *net.IPNet) bool {
	return left.Contains(right.IP) || right.Contains(left.IP)
}

func CopyIP(from net.IP) net.IP {
	to := make(net.IP, len(from))
	copy(to, from)
	return to
}

func SliceNetFromNet(netX, netY *net.IPNet) []IPRange {
	min, max := NetworkRange(netX)
	ipRange := IPRange{Start: min, End: max}
	return SliceNetFromRange(ipRange, netY)
}

func SliceNetFromRange(ipRange IPRange, ipNet *net.IPNet) []IPRange {
	if !ipRange.OverlapsNet(ipNet) {
		return []IPRange{ipRange}
	}

	if ipRange.EqualsNet(ipNet) {
		return nil
	}

	rangeStart := ipRange.Start.To4()
	rangeEnd := ipRange.End.To4()

	min, max := NetworkRange(ipNet)
	netMin := min.To4()
	netMax := max.To4()

	var ipRanges []IPRange
	if bytes.Compare(rangeStart, netMin) == -1 {
		ipRanges = append(ipRanges, IPRange{
			Start: ipRange.Start,
			End:   Dec(min),
		})
	}

	if bytes.Compare(rangeEnd, netMax) == 1 {
		ipRanges = append(ipRanges, IPRange{
			Start: Inc(max),
			End:   ipRange.End,
		})
	}

	return ipRanges
}
