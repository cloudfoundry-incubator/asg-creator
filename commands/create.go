package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/asg-creator/commands/internal/flaghelpers"
	"github.com/cloudfoundry-incubator/asg-creator/config"
)

type CreateCommand struct {
	Config flaghelpers.Path `long:"config" required:"true"`
}

func (c *CreateCommand) Execute(args []string) error {
	configFile := c.Config

	createConfig, err := config.LoadCreateConfig(string(configFile))
	if err != nil {
		return err
	}

	if createConfig.PublicNetworks {
		rules := createConfig.PublicNetworksRules()

		bs, err := json.Marshal(rules)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, "Wrote public-networks.json\n")

		ioutil.WriteFile("public-networks.json", bs, os.ModePerm)
	}

	if createConfig.PrivateNetworks {
		rules := createConfig.PrivateNetworksRules()

		bs, err := json.Marshal(rules)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, "Wrote private-networks.json\n")

		ioutil.WriteFile("private-networks.json", bs, os.ModePerm)
	}

	fmt.Fprintln(os.Stdout, "OK")

	return nil
}
