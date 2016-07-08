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
	Config     flaghelpers.Path `long:"config" short:"c"`
	OutputPath string           `long:"output" short:"o"`
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

	if includedNetworksRules := cfg.IncludedNetworksRules(); len(includedNetworksRules) != 0 {
		if c.OutputPath == "" {
			return fmt.Errorf("--output is required when config contains included_networks")
		}

		networkRulesBytes, err := rulesBytes(includedNetworksRules)
		if err != nil {
			return err
		}

		err = writeFile(c.OutputPath, networkRulesBytes)
		if err != nil {
			return err
		}
	} else {
		publicRulesBytes, err := rulesBytes(cfg.PublicNetworksRules())
		if err != nil {
			return err
		}

		err = writeFile("public-networks.json", publicRulesBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		privateRulesBytes, err := rulesBytes(cfg.PrivateNetworksRules())
		if err != nil {
			return err
		}

		err = writeFile("private-networks.json", privateRulesBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	fmt.Fprintln(os.Stdout, "OK")

	return nil
}

func writeFile(filepath string, filebytes []byte) error {
	err := ioutil.WriteFile(filepath, filebytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to write private-networks.json: %s\n", err.Error()))
	}
	fmt.Printf("Wrote %s", filepath)
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
