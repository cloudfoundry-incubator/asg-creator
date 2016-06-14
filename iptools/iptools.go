package iptools

import (
	"bytes"
	"net"
)

func PrivateIPNets() []*net.IPNet {
	return []*net.IPNet{
		{
			IP:   net.IP{10, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		},
		{
			IP:   net.IP{172, 16, 0, 0},
			Mask: net.IPMask{255, 240, 0, 0},
		},
		{
			IP:   net.IP{192, 168, 0, 0},
			Mask: net.IPMask{255, 255, 0, 0},
		},
	}
}

func PublicIPRanges() []IPRange {
	nets := []net.IP{
		net.IP{0, 0, 0, 0},
	}

	for _, ipNet := range PrivateIPNets() {
		minusOne := Dec(ipNet.IP)
		nets = append(nets, minusOne)

		_, max := NetworkRange(ipNet)
		plusOne := Inc(max)
		nets = append(nets, plusOne)
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

	min, max := NetworkRange(ipNet)

	var ipRanges []IPRange
	var minusOne net.IP
	if bytes.Compare(ipRange.Start, min) == -1 && bytes.Compare(min, ipRange.End) == -1 {
		minusOne = Dec(min)
		ipRanges = append(ipRanges, IPRange{
			Start: ipRange.Start,
			End:   minusOne,
		})
	}

	var plusOne net.IP
	if bytes.Compare(ipRange.Start, max) == -1 && bytes.Compare(max, ipRange.End) == -1 {
		plusOne = Inc(max)
		ipRanges = append(ipRanges, IPRange{
			Start: plusOne,
			End:   ipRange.End,
		})
	}

	return ipRanges
}
