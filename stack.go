package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"github.com/goccy/go-yaml"
    "io/ioutil"
)

type stackCompose struct {
    Version  string `json:"version" yaml:"version"`
    Services map[string]struct {
        Registry string `json:"registry" yaml:"registry"`
        Image string `json:"image" yaml:"image"`
        Ports []string `json:"ports" yaml:"ports"`
        Environment map[string]string `json:"environment" yaml:"environment"`
        Hostname string `json:"hostname" yaml:"hostname"`
        ContainerName string `json:"container_name" yaml:"container_name"`
        Volumes []string `json:"volumes" yaml:"volumes"`
        Networks []string `json:"networks" yaml:"networks"`
    } `json:"services" yaml:"services"`
    Networks map[string]struct {
        Driver string `json:"driver" yaml:"driver"`
    } `json:"networks" yaml:"networks"`
    Volumes map[string]struct {
        Driver string `json:"driver" yaml:"driver"`
    } `json:"volumes" yaml:"volumes"`
    Registry []struct {
        Url string `json:"url" yaml:"url"`
        Username string `json:"username" yaml:"username"`
        Password string `json:"password" yaml:"password"`
        Name string `json:"name" yaml:"name"`
    } `json:"registry" yaml:"registry"`
}

func stackDeployRun(compose stackCompose) error {
    fmt.Println("Deploying stack...")

    r, _ := json.MarshalIndent(compose, "", "    ")
    fmt.Printf("%s\n", r)
    return nil
}

func stackDeploy(c *cli.Context) error { 
    var compose stackCompose
    var err error

    if c.Args().Len() == 0 {
        log.Fatal("No compose file specified")
    }

    composeFile := c.Args().Get(0)

    file, _ := ioutil.ReadFile(composeFile)
    if c.String("format") == "json" {
        err = json.Unmarshal([]byte(file), &compose)
    } else {
        err = yaml.Unmarshal([]byte(file), &compose)
    }

    if err != nil {
        log.Fatal(err)
    }

    err = stackDeployRun(compose)
    if err != nil {
        log.Fatal(err)
    }

    return nil
}
