package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/asg-creator/asg"
	"github.com/cloudfoundry-incubator/asg-creator/commands/internal/flaghelpers"
	"github.com/cloudfoundry-incubator/asg-creator/config"
)

type CreateCommand struct {
	Config flaghelpers.Path `long:"config" short:"c"`
}

func (c *CreateCommand) Execute(args []string) error {
	cfg := config.Create{}

	if c.Config != "" {
		var err error
		cfg, err = config.LoadCreateConfig(string(c.Config))
		if err != nil {
			return err
		}
	}

	publicRulesBytes, err := rulesBytes(cfg.PublicNetworksRules())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("public-networks.json", publicRulesBytes, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to write public-networks.json: %s\n", err.Error()))
	} else {
		fmt.Fprintln(os.Stdout, "Wrote public-networks.json")
	}

	privateRulesBytes, err := rulesBytes(cfg.PrivateNetworksRules())
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

func rulesBytes(rules []asg.Rule) ([]byte, error) {
	bs, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = json.Indent(&b, bs, "", "\t")
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
