package config

import (
	"io/ioutil"
	"net"

	"github.com/cloudfoundry-incubator/asg-creator/asg"
	"github.com/cloudfoundry-incubator/asg-creator/iptools"
	"github.com/cloudfoundry-incubator/candiedyaml"
)

const protocolAll = "all"

var blacklistedIPs = []net.IP{
	{169, 254, 169, 254},
}

type Create struct {
	IncludedNetworks       []iptools.IPRange `yaml:"included_networks"`
	ExcludedCIDRRanges     []iptools.IPRange `yaml:"excluded_networks"`
	ExcludedSingleIPRanges []iptools.IPRange `yaml:"excluded_ips"`
	ExcludedRanges         []iptools.IPRange `yaml:"excluded_ranges"`
}

func LoadCreateConfig(path string) (Create, error) {
	createConfig := new(Create)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return Create{}, err
	}

	err = candiedyaml.Unmarshal(bs, createConfig)
	if err != nil {
		return Create{}, err
	}

	return *createConfig, nil
}

func (c *Create) IncludedNetworksRules() []asg.Rule {
	ipRanges := make(chan iptools.IPRange)
	go func() {
		for i := range c.IncludedNetworks {
			ipRanges <- c.IncludedNetworks[i]
		}
		close(ipRanges)
	}()
	return c.rulesForRanges(ipRanges)
}

func (c *Create) PublicNetworksRules() []asg.Rule {
	ipRanges := make(chan iptools.IPRange)
	go func() {
		for _, ipRange := range iptools.PublicIPRanges() {
			ipRanges <- ipRange
		}
		close(ipRanges)
	}()

	return c.rulesForRanges(ipRanges)
}

func (c *Create) PrivateNetworksRules() []asg.Rule {
	ipRanges := make(chan iptools.IPRange)
	go func() {
		for _, ipNet := range iptools.PrivateIPNets() {
			ipRanges <- iptools.NewIPRangeFromIPNet(ipNet)
		}
		close(ipRanges)
	}()

	return c.rulesForRanges(ipRanges)
}

func (c *Create) rulesForRanges(ipRangesCh chan iptools.IPRange) []asg.Rule {
	var rules []asg.Rule

	ipRanges := c.filterBlacklistedIPs(c.filterExcludedSingleIPRanges(c.filterExcludedRanges(c.filterExcludedCIDRRanges(ipRangesCh))))
	for ipRange := range ipRanges {
		rules = append(rules, asg.Rule{
			Destination: ipRange.String(),
			Protocol:    protocolAll,
		})
	}

	return rules
}

func (c *Create) filterExcludedSingleIPRanges(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
	out := make(chan iptools.IPRange)
	go func() {
		for ipRange := range ipRanges {
			if len(c.ExcludedSingleIPRanges) == 0 {
				out <- ipRange
				continue
			}

			var ips []net.IP
			for i := range c.ExcludedSingleIPRanges {
				ips = append(ips, c.ExcludedSingleIPRanges[i].Start)
			}
			for _, newRange := range ipRange.SliceIPs(ips) {
				out <- newRange
			}
		}
		close(out)
	}()
	return out
}

func (c *Create) filterBlacklistedIPs(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
	out := make(chan iptools.IPRange)
	go func() {
		for ipRange := range ipRanges {
			for _, newRange := range ipRange.SliceIPs(blacklistedIPs) {
				out <- newRange
			}
		}
		close(out)
	}()
	return out
}

func (c *Create) filterExcludedCIDRRanges(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
	out := make(chan iptools.IPRange)
	go func() {
		for ipRange := range ipRanges {
			if len(c.ExcludedCIDRRanges) == 0 {
				out <- ipRange
				continue
			}
			for _, excludedRange := range c.ExcludedCIDRRanges {
				for _, newRange := range ipRange.SliceRange(excludedRange) {
					out <- newRange
				}
			}
		}
		close(out)
	}()
	return out
}

func (c *Create) filterExcludedRanges(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
	out := make(chan iptools.IPRange)
	go func() {
		for ipRange := range ipRanges {
			if len(c.ExcludedRanges) == 0 {
				out <- ipRange
				continue
			}

			for _, newRange := range ipRange.SliceRanges(c.ExcludedRanges) {
				out <- newRange
			}
		}
		close(out)
	}()
	return out
}
