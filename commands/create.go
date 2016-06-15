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
	Config flaghelpers.Path `long:"config" short:"c"`
}

func (c *CreateCommand) Execute(args []string) error {
	configFile := c.Config

	cfg := config.Create{}

	if configFile != "" {
		var err error
		cfg, err = config.LoadCreateConfig(string(configFile))
		if err != nil {
			return err
		}
	}

	publicRulesBytes, err := json.Marshal(cfg.PublicNetworksRules())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("public-networks.json", publicRulesBytes, os.ModePerm)

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to write public-networks.json: %s\n", err.Error()))
	} else {
		fmt.Fprintln(os.Stdout, "Wrote public-networks.json")
	}

	privateRulesBytes, err := json.Marshal(cfg.PrivateNetworksRules())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("private-networks.json", privateRulesBytes, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to write private-networks.json: %s\n", err.Error()))
	} else {
		fmt.Fprintln(os.Stdout, "Wrote private-networks.json")
	}

	fmt.Fprintln(os.Stdout, "OK")

	return nil
}
