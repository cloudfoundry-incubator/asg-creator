package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/cloudfoundry-incubator/asg-creator/asg"
	"github.com/cloudfoundry-incubator/asg-creator/iptools"

	"gopkg.in/yaml.v2"
)

var blacklistedIPs = []net.IP{
	{169, 254, 169, 254},
}

type Create struct {
	ExcludedNetworks []string `yaml:"excluded_networks"`
	ExcludedIPs      []string `yaml:"excluded_ips"`
}

func LoadCreateConfig(path string) (Create, error) {
	createConfig := new(Create)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return Create{}, err
	}

	err = yaml.Unmarshal(bs, createConfig)
	if err != nil {
		return Create{}, err
	}

	return *createConfig, nil
}

func (c *Create) PublicNetworksRules() []asg.Rule {
	var rules []asg.Rule

	ipRanges := make(chan iptools.IPRange)
	go func() {
		for _, ipRange := range iptools.PublicIPRanges() {
			ipRanges <- ipRange
		}
		close(ipRanges)
	}()

	filteredIPRanges := c.blacklistedIPFilter(c.ipFilter(c.networkFilter(ipRanges)))
	for ipRange := range filteredIPRanges {
		rules = append(rules, asg.Rule{
			Destination: fmt.Sprintf("%s-%s", ipRange.Start, ipRange.End),
			Protocol:    "all",
		})
	}

	return rules
}

func (c *Create) PrivateNetworksRules() []asg.Rule {
	var rules []asg.Rule

	ipRanges := make(chan iptools.IPRange)
	go func() {
		for _, ipNet := range iptools.PrivateIPNets() {
			ipRanges <- iptools.NewIPRangeFromIPNet(ipNet)
		}
		close(ipRanges)
	}()

	filteredIPRanges := c.blacklistedIPFilter(c.ipFilter(c.networkFilter(ipRanges)))
	for ipRange := range filteredIPRanges {
		rules = append(rules, asg.Rule{
			Destination: fmt.Sprintf("%s-%s", ipRange.Start, ipRange.End),
			Protocol:    "all",
		})
	}

	return rules
}

func (c *Create) ipFilter(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
	out := make(chan iptools.IPRange)
	go func() {
		for ipRange := range ipRanges {
			var ips []net.IP
			for i := range c.ExcludedIPs {
				ips = append(ips, net.ParseIP(c.ExcludedIPs[i]))
			}
			for _, newRange := range ipRange.SliceIPs(ips) {
				out <- newRange
			}
		}
		close(out)
	}()
	return out
}

func (c *Create) blacklistedIPFilter(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
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

func (c *Create) networkFilter(ipRanges <-chan iptools.IPRange) <-chan iptools.IPRange {
	out := make(chan iptools.IPRange)
	go func() {
		for ipRange := range ipRanges {
			if len(c.ExcludedNetworks) == 0 {
				out <- ipRange
				continue
			}
			for _, excludedNetwork := range c.ExcludedNetworks {
				_, excludedIPNet, err := net.ParseCIDR(excludedNetwork)
				if err != nil {
					log.Fatalf("non-CIDR given as network in config: %s", excludedNetwork)
				}

				if ipRange.OverlapsNet(excludedIPNet) {
					for _, newRange := range iptools.SliceNetFromRange(ipRange, excludedIPNet) {
						out <- newRange
					}
				} else {
					out <- ipRange
				}
			}
		}
		close(out)
	}()
	return out
}
