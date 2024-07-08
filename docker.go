package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
    "github.com/SaaShup/paashup-sdk/docker"
    "github.com/SaaShup/paashup-sdk/netbox"
	"log"
    "time"
)

type listContainerStruct struct {
	Id            int    `json:"id" yaml:"id"`
	Name          string `json:"name" yaml:"name"`
	Host          string `json:"host" yaml:"host"`
	Image         string `json:"image" yaml:"image"`
	State         string `json:"state" yaml:"state"`
	PortsCount    int    `json:"ports_count" yaml:"ports_count"`
	MountsCount   int    `json:"mounts_count" yaml:"mounts_count"`
	BindsCount    int    `json:"binds_count" yaml:"binds_count"`
	NetworksCount int    `json:"networks_count" yaml:"networks_count"`
	EnvCount      int    `json:"env_count" yaml:"env_count"`
	LabelsCount   int    `json:"labels_count" yaml:"labels_count"`
}


func psContainers(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token
    
	hostname := c.Args().First()
    var result docker.ContainerListStruct
	if hostname != "" {
		host, err := docker.HostSearchByName(c.Args().First())
		if err != nil {
			fmt.Println("Host not found")
			return nil
		}
        result, err = docker.ContainerListByHost(host.Id)
        if err != nil {
            fmt.Println("Error getting containers")
            return nil
        }
	} else {
        var err error
		result, err = docker.ContainerList()
        if err != nil {
            fmt.Println("Error getting containers")
            return nil
        }
	}

	if c.String("format") == "" {
		table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
			return ""
		}
		table.DefaultWidthFunc = calcWidhtColorRed
		var tbl table.Table
		// Does not allow to disable header.. :(
		// should we migrate to another table library?
		tbl = table.New("", "", "", "", "", "", "", "", "", "", "").WithPadding(2)
		for _, rec := range result.Results {
			if rec.State == "running" {
				tbl.AddRow(rec.Id, rec.Name, rec.Host.Name, rec.Image.Name, rec.State, fmt.Sprintf("%d P", len(rec.Ports)),
					fmt.Sprintf("%d M", len(rec.Mounts)), fmt.Sprintf("%d B", len(rec.Binds)), fmt.Sprintf("%d N", len(rec.Network_settings)),
					fmt.Sprintf("%d E", len(rec.Env)), fmt.Sprintf("%d L", len(rec.Labels)))
			} else {
				tbl.AddRow(color.RedString("%d", rec.Id), color.RedString("%s", rec.Name), color.RedString("%s", rec.Host.Name),
					color.RedString("%s", rec.Image.Name), color.RedString("%s", rec.State), color.RedString("%d P", len(rec.Ports)),
					color.RedString("%d M", len(rec.Mounts)), color.RedString("%d B", len(rec.Binds)), color.RedString("%d N", len(rec.Network_settings)),
					color.RedString("%d E", len(rec.Env)), color.RedString("%d L", len(rec.Labels)))
			}
		}
		tbl.Print()
	} else {
		var listContainer []listContainerStruct
		for _, rec := range result.Results {
			listContainer = append(listContainer, listContainerStruct{Id: rec.Id, Name: rec.Name, Host: rec.Host.Name, Image: rec.Image.Name, State: rec.State, PortsCount: len(rec.Ports), MountsCount: len(rec.Mounts), BindsCount: len(rec.Binds), NetworksCount: len(rec.Network_settings), EnvCount: len(rec.Env), LabelsCount: len(rec.Labels)})
		}
		switch c.String("format") {
		case "json-pretty":
			resp, err := json.MarshalIndent(listContainer, "", "    ")
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "json":
			resp, err := json.Marshal(listContainer)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "yaml":
			resp, err := yaml.Marshal(listContainer)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		default:
			resp, err := json.Marshal(listContainer)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		}
	}
	return nil

}

func inspectContainer(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

	host, err := docker.HostSearchByName(c.Args().First())
	if err != nil {
		fmt.Println("Host not found")
		return nil
	}

	container, err := docker.ContainerSearchByName(host, c.Args().Get(c.Args().Len()-1))
	if err != nil {
		fmt.Println("Container not found")
		return nil
	}

	resultCall, err := docker.ContainerInspect(container.Id)

	if err != nil {
		log.Fatal(err)
	}

	switch c.String("format") {
	case "json-pretty":
		resp, err := json.MarshalIndent(resultCall, "", "    ")
		if err == nil {
			fmt.Printf("%s\n", resp)
		}
	case "json":
		resp, err := json.Marshal(resultCall)
		if err == nil {
			fmt.Printf("%s\n", resp)
		}
	case "yaml":
		resp, err := yaml.Marshal(resultCall)
		if err == nil {
			fmt.Printf("%s\n", resp)
		}
	default:
		resp, err := json.Marshal(resultCall)
		if err == nil {
			fmt.Printf("%s\n", resp)
		}
	}

	return nil

}

func execContainer(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

	host, err := docker.HostSearchByName(c.Args().First())
	if err != nil {
		fmt.Println("Host not found")
		return nil
	}

	container, err := docker.ContainerSearchByName(host, c.Args().Get(c.Args().Len()-2))
	if err != nil {
		fmt.Println("Container not found")
		return nil
	}

	resultCall, err := docker.ContainerExec(container.Id, c.Args().Get(c.Args().Len() - 1))

	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(resultCall)

	return nil

}

func stopContainer(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

    host, err := docker.HostSearchByName(c.Args().First())
    if err != nil {
        fmt.Println("Host not found")
        return nil
    }

    container, err := docker.ContainerSearchByName(host, c.Args().Get(c.Args().Len()-1))
    if err != nil {
        fmt.Println("Container not found")
        return nil
    }

	operationContainer, err := docker.ContainerStop(container)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if operationContainer.Operation == "stop" {
        if !c.Bool("wait") {
		    i := 0
		    for i < 20 {
                resultCall, err := docker.ContainerInspect(operationContainer.Id)
			    if err != nil {
				    log.Fatal(err)
			    }
			    if resultCall.State == "exited" {
                    break
                }
			    time.Sleep(1 * time.Second)
			    i++
		    }
		    return fmt.Errorf("Timeout")
        }
	
	} else {
		fmt.Println("Operation not executed")
	}

	fmt.Println("Container " + operationContainer.Name + " stopped")
	return nil
}

func startContainer(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

    host, err := docker.HostSearchByName(c.Args().First())
    if err != nil {
        fmt.Println("Host not found")
        return nil
    }

    container, err := docker.ContainerSearchByName(host, c.Args().Get(c.Args().Len()-1))
    if err != nil {
        fmt.Println("Container not found")
        return nil
    }

	operationContainer, err := docker.ContainerStart(container)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if operationContainer.Operation == "start" {
        if !c.Bool("wait") {
		    i := 0
		    for i < 20 {
                resultCall, err := docker.ContainerInspect(operationContainer.Id)
			    if err != nil {
				    log.Fatal(err)
			    }
			    if resultCall.State == "running" {
                    break
                }
			    time.Sleep(1 * time.Second)
			    i++
		    }
		    return fmt.Errorf("Timeout")
        }
	
	} else {
		fmt.Println("Operation not executed")
	}

	fmt.Println("Container " + operationContainer.Name + " started")
	return nil
}

func getLogs(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

	host, err := docker.HostSearchByName(c.Args().First())
	if err != nil {
		fmt.Println("Host not found")
		return nil
	}

	container, err := docker.ContainerSearchByName(host, c.Args().Get(c.Args().Len()-1))
	if err != nil {
		fmt.Println("Container not found")
		return nil
	}

	resultCall, err := docker.ContainerLogs(container.Id)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", resultCall)

	return nil
}

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
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

	result, err := docker.HostList()

	if err != nil {
		log.Fatal(err)
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

func inspectHost(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

    host := c.Args().First()
    if host == "" {
        cli.ShowAppHelpAndExit(c, 1)
    }

    hostData, err := docker.HostSearchByName(host)
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
