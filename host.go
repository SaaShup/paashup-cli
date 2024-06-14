package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/rodaine/table"
)

type listHostStruct struct {
    Id   int    `json:"id" yaml:"id"`
    Name string `json:"name" yaml:"name"`
    State string `json:"state" yaml:"state"`
    ImageCount int `json:"image_count" yaml:"image_count"`
    ContainerCount int `json:"container_count" yaml:"container_count"`
    VolumeCount int `json:"volume_count" yaml:"volume_count"`
    NetworkCount int `json:"network_count" yaml:"network_count"`
    AgentVersion string `json:"agent_version" yaml:"agent_version"`
}

func listHosts(c *cli.Context) error {
	url := "hosts/"

	var result HostList
	resultCall, err := netboxCall(c, url, "GET", nil)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(resultCall, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if c.String("format") == "" {
		table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
			return ""
		}
		table.DefaultWidthFunc = calcWidhtColorRed
		var tbl table.Table
		// Does not allow to disable header.. :(
		// should we migrate to another table library?
		tbl = table.New("", "", "", "", "", "", "", "").WithPadding(2)
		for _, rec := range result.Results {
			if rec.State == "running" {
				tbl.AddRow(rec.Id, rec.Name, rec.State, fmt.Sprintf("%d I", len(rec.Images)), fmt.Sprintf("%d C", len(rec.Containers)),
                    fmt.Sprintf("%d V", len(rec.Volumes)), fmt.Sprintf("%d N", len(rec.Networks)), rec.Agent_version)
			} else {
				tbl.AddRow(color.RedString("%d", rec.Id), color.RedString("%s", rec.Name), color.RedString("%s", rec.State),
                    color.RedString("%d I", len(rec.Images)), color.RedString("%d C", len(rec.Containers)), color.RedString("%d V", len(rec.Volumes)), 
                    color.RedString("%d N", len(rec.Networks)), color.RedString("%s", rec.Agent_version))
			}
		}
		tbl.Print()
	} else {
		var listHost []listHostStruct
		for _, rec := range result.Results {
            listHost = append(listHost, listHostStruct{rec.Id, rec.Name, rec.State, len(rec.Images), len(rec.Containers), len(rec.Volumes), len(rec.Networks), rec.Agent_version})
		}
		switch c.String("format") {
		case "json-pretty":
			resp, err := json.MarshalIndent(listHost, "", "    ")
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "json":
			resp, err := json.Marshal(listHost)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "yaml":
			resp, err := yaml.Marshal(listHost)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		default:
			resp, err := json.Marshal(listHost)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		}
	}
	return nil


}

func searchHost(host string, c *cli.Context) (HostComplete, error) {

	if host == "" {
		return HostComplete{}, fmt.Errorf("Host not found")
	}

	url := fmt.Sprintf("hosts/?name=%s", host)
	var result HostList
	resultCall, err := netboxCall(c, url, "GET", nil)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(resultCall, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if result.Count == 1 {
		return result.Results[0], nil
	} else {
		return HostComplete{}, fmt.Errorf("Host not found")
	}

}

func inspectHost(c *cli.Context) error {
    host := c.Args().First()
    if host == "" {
        cli.ShowAppHelpAndExit(c, 1)
    }

    hostData, err := searchHost(host, c)
    if err != nil {
        log.Fatal(err)
    }
    switch c.String("format") {
    case "json-pretty":
        resp, err := json.MarshalIndent(hostData, "", "    ")
        if err == nil {
            fmt.Printf("%s\n", resp)
        }
    case "json":
        resp, err := json.Marshal(hostData)
        if err == nil {
            fmt.Printf("%s\n", resp)
        }
    case "yaml":
        resp, err := yaml.Marshal(hostData)
        if err == nil {
            fmt.Printf("%s\n", resp)
        }
    default:
        resp, err := json.Marshal(hostData)
        if err == nil {
            fmt.Printf("%s\n", resp)
        }
    }
    return nil
}
