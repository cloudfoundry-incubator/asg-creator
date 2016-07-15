package config

import (
	"io/ioutil"
	"net"

	"github.com/cloudfoundry-incubator/asg-creator/asg"
	"github.com/cloudfoundry-incubator/asg-creator/iptools"
	"github.com/cloudfoundry-incubator/candiedyaml"
)

const protocolAll = "all"

var linkLocalIPRange = iptools.IPRange{
	Start: net.IP{169, 254, 0, 0},
	End:   net.IP{169, 254, 255, 255},
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
	return c.rulesFromRanges(c.IncludedNetworks)
}

func (c *Create) PublicNetworksRules() []asg.Rule {
	return c.rulesFromRanges(iptools.PublicIPRanges())
}

func (c *Create) PrivateNetworksRules() []asg.Rule {
	return c.rulesFromRanges(iptools.PrivateIPRanges())
}

func (c *Create) rulesFromRanges(baseIPRanges []iptools.IPRange) []asg.Rule {
	var excludedIPRanges []iptools.IPRange

	for i := range c.ExcludedCIDRRanges {
		excludedIPRanges = append(excludedIPRanges, c.ExcludedCIDRRanges[i])
	}

	for i := range c.ExcludedSingleIPRanges {
		excludedIPRanges = append(excludedIPRanges, c.ExcludedSingleIPRanges[i])
	}

	for i := range c.ExcludedRanges {
		excludedIPRanges = append(excludedIPRanges, c.ExcludedRanges[i])
	}

	excludedIPRanges = append(excludedIPRanges, linkLocalIPRange)

	var rules []asg.Rule
	for i := range baseIPRanges {
		for _, newRange := range baseIPRanges[i].SliceRanges(excludedIPRanges) {
			rules = append(rules, asg.Rule{
				Destination: newRange.String(),
				Protocol:    protocolAll,
			})
		}
	}

	return rules
}
