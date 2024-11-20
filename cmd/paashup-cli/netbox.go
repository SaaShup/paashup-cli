package main

import (
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
    "github.com/SaaShup/paashup-cli/internal/config"
)

type NetboxConfig struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

func useNetboxConfig(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return nil
	}
	var configpath string

	if os.Getenv("XDG_CONFIG_HOME") == "" {
		configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
	} else {
		configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
	}

	_ = ioutil.WriteFile(configpath+"current", []byte(c.Args().First()), 0644)
	println("Using " + c.Args().First() + " config")
	return nil
}

func setNetboxConfig(c *cli.Context) error {
	if c.Args().Len() != 3 {
		return cli.ShowCommandHelp(c, "set-config")
	}
	name := c.Args().Get(0)
	url := c.Args().Get(1)
	token := c.Args().Get(2)
	return config.SetConfig(name, url, token)
}

