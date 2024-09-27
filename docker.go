package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
    "github.com/SaaShup/paashup-sdk/pkg/docker"
    "github.com/SaaShup/paashup-sdk/pkg/netbox"
	"log"
    "time"
)

type listRegistryStruct struct {
    Id   int    `json:"id" yaml:"id"`
    Name string `json:"name" yaml:"name"`
    Url  string `json:"url" yaml:"url"`
    Username string `json:"username" yaml:"username"`
    Password string `json:"password" yaml:"password"`
    ImageCount int `json:"image_count" yaml:"image_count"`
}

func listRegistries(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token
    
	hostname := c.Args().First()
    var result docker.RegistryListStruct
	if hostname != "" {
		host, err := docker.HostSearchByName(c.Args().First())
		if err != nil {
			fmt.Println("Host not found")
			return nil
		}
        result, err = docker.RegistryListByHost(host.Id)
        if err != nil {
            fmt.Println(err)
            fmt.Println("Error getting Registry")
            return nil
        }
	} else {
        var err error
		result, err = docker.RegistryList()
        if err != nil {
            fmt.Println(err)
            fmt.Println("Error getting Registry")
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
			tbl.AddRow(rec.Id, rec.Name, rec.Host.Name, rec.Username, rec.Password, rec.Url, fmt.Sprintf("%d I", len(rec.Images)))
		}
		tbl.Print()
	} else {
		var listRegistry []listRegistryStruct
		for _, rec := range result.Results {
            listRegistry = append(listRegistry, listRegistryStruct{Id: rec.Id, Name: rec.Name, Username: rec.Username, Password: rec.Password, Url: rec.Url, ImageCount: len(rec.Images)})
		}
		switch c.String("format") {
		case "json-pretty":
			resp, err := json.MarshalIndent(listRegistry, "", "    ")
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "json":
			resp, err := json.Marshal(listRegistry)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "yaml":
			resp, err := yaml.Marshal(listRegistry)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		default:
			resp, err := json.Marshal(listRegistry)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		}
	}
	return nil
   
}


type listVolumeStruct struct {
    Id          int    `json:"id" yaml:"id"`
    Name        string `json:"name" yaml:"name"`
    Host        string `json:"host" yaml:"host"`
    Driver      string `json:"driver" yaml:"driver"`
    MountCount  int     `json:"mount_count" yaml:"mount_count"`
}

func listVolumes(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token
    
	hostname := c.Args().First()
    var result docker.VolumeListStruct
	if hostname != "" {
		host, err := docker.HostSearchByName(c.Args().First())
		if err != nil {
			fmt.Println("Host not found")
			return nil
		}
        result, err = docker.VolumeListByHost(host.Id)
        if err != nil {
            fmt.Println("Error getting Volumes")
            return nil
        }
	} else {
        var err error
		result, err = docker.VolumeList()
        if err != nil {
            fmt.Println("Error getting Volumes")
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
			tbl.AddRow(rec.Id, rec.Name, rec.Host.Name, rec.Driver, fmt.Sprintf("%d M", len(rec.Mounts)))
		}
		tbl.Print()
	} else {
		var listVolumes []listVolumeStruct
		for _, rec := range result.Results {
            listVolumes = append(listVolumes, listVolumeStruct{Id: rec.Id, Name: rec.Name, Host: rec.Host.Name, Driver: rec.Driver, MountCount: len(rec.Mounts)})
		}
		switch c.String("format") {
		case "json-pretty":
			resp, err := json.MarshalIndent(listVolumes, "", "    ")
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "json":
			resp, err := json.Marshal(listVolumes)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "yaml":
			resp, err := yaml.Marshal(listVolumes)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		default:
			resp, err := json.Marshal(listVolumes)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		}
	}
	return nil
   
}


type listImageStruct struct {
    Id          int    `json:"id" yaml:"id"`
    Name        string `json:"name" yaml:"name"`
    Host        string `json:"host" yaml:"host"`
    Size        int `json:"size" yaml:"size"`
    Digest      string `json:"digest" yaml:"digest"`
    ImageID     string `json:"imageID" yaml:"imageID"`
    Registry    string `json:"registry" yaml:"registry"`
    ContainerCount int `json:"container_count" yaml:"container_count"`
}

func listImages(c *cli.Context) error {
    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token
    
	hostname := c.Args().First()
    var result docker.ImageListStruct
	if hostname != "" {
		host, err := docker.HostSearchByName(c.Args().First())
		if err != nil {
			fmt.Println("Host not found")
			return nil
		}
        result, err = docker.ImageListByHost(host.Id)
        if err != nil {
            fmt.Println("Error getting Images")
            return nil
        }
	} else {
        var err error
		result, err = docker.ImageList()
        if err != nil {
            fmt.Println("Error getting Images")
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
			if rec.Size > 0 {
				tbl.AddRow(rec.Id, rec.Name, rec.Host.Name, rec.Registry.Name, rec.Image.Name, rec.Size, rec.ImageID, rec.Digest, fmt.Sprintf("%d P", len(rec.Containers)))
			} else {
				tbl.AddRow(color.RedString("%d", rec.Id), color.RedString("%s", rec.Name), color.RedString("%s", rec.Host.Name), color.RedString("%s", rec.Registry.Name),
					color.RedString("%s", rec.Image.Name), color.RedString("%s", rec.Size), color.RedString("%d P", len(rec.ImageID)),
					color.RedString("%d M", len(rec.Digest)), color.RedString("%d B", len(rec.Containers)))
			}
		}
		tbl.Print()
	} else {
		var listImages []listImageStruct
		for _, rec := range result.Results {
            listImages = append(listImages, listImageStruct{Id: rec.Id, Name: rec.Name, Host: rec.Host.Name, Registry: rec.Registry.Name, Size: rec.Size, ImageID: rec.ImageID, Digest: rec.Digest, ContainerCount: len(rec.Containers)})
		}
		switch c.String("format") {
		case "json-pretty":
			resp, err := json.MarshalIndent(listImages, "", "    ")
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "json":
			resp, err := json.Marshal(listImages)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		case "yaml":
			resp, err := yaml.Marshal(listImages)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		default:
			resp, err := json.Marshal(listImages)
			if err == nil {
				fmt.Printf("%s\n", resp)
			}
		}
	}
	return nil

   
}

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
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

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
    if c.Args().Len() != 3 {
        fmt.Println("Please provide a host, a container name and a command")
        cli.ShowAppHelpAndExit(c, 1)
    }

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

func operationContainer(containerName string, hostName string, operation string, wait bool) error{
    config, _ := readConfig(nil)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

    host, err := docker.HostSearchByName(hostName)
    if err != nil {
        return fmt.Errorf("Host not found")
    }

    container, err := docker.ContainerSearchByName(host, containerName)
    if err != nil {
        return fmt.Errorf("Container not found")
    }

    var operationContainer docker.Container
    switch operation {
    case "start":
        operationContainer, err = docker.ContainerStart(container)
    case "stop":
        operationContainer, err = docker.ContainerStop(container)
    case "restart":
        operationContainer, err = docker.ContainerRestart(container)
    case "kill":
        operationContainer, err = docker.ContainerKill(container)
    case "recreate":
        operationContainer, err = docker.ContainerRecreate(container)
    }

    if err != nil {
        fmt.Println(err)
        return nil
    }
    if operationContainer.Operation == operation {
        if !wait {
		    i := 0
		    for i < 20 {
                resultCall, err := docker.ContainerInspect(operationContainer.Id)
			    if err != nil {
				    log.Fatal(err)
			    }
                switch operation {
                case "start":
                    if resultCall.State == "running" {
                        break
                    }
                case "stop":
                    if resultCall.State == "exited" {
                        break
                    }
                case "restart":
                    if resultCall.State == "running" {
                        break
                    }
                case "kill":
                    if resultCall.State == "exited" {
                        break
                    }
                case "recreate":
                    if resultCall.State == "running" {
                        break
                    }
                }
                
			    time.Sleep(1 * time.Second)
			    i++
		    }
		    return fmt.Errorf("Timeout")
        }

    } else {
        fmt.Println("Operation not executed")
    }
    return nil
}

func stopContainer(c *cli.Context) error {
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

    err := operationContainer(c.Args().Get(c.Args().Len()-1), c.Args().First(), "stop", c.Bool("wait"))
    if err != nil {
        fmt.Println(err)
    }

	fmt.Println("Container %s stopped\n", c.Args().Get(c.Args().Len()-1))
	return nil
}

func startContainer(c *cli.Context) error {
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

    err := operationContainer(c.Args().Get(c.Args().Len()-1), c.Args().First(), "stop", c.Bool("wait"))
    if err != nil {
        fmt.Println(err)
    }

	fmt.Printf("Container %s started\n", c.Args().Get(c.Args().Len()-1))
	return nil
}

func restartContainer(c *cli.Context) error {
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

    err := operationContainer(c.Args().Get(c.Args().Len()-1), c.Args().First(), "restart", c.Bool("wait"))
    if err != nil {
        fmt.Println(err)
    }

	fmt.Printf("Container %s restarted\n", c.Args().Get(c.Args().Len()-1))
	return nil
}

func killContainer(c *cli.Context) error {
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

    err := operationContainer(c.Args().Get(c.Args().Len()-1), c.Args().First(), "kill", c.Bool("wait"))
    if err != nil {
        fmt.Println(err)
    }

	fmt.Printf("Container %s killed\n", c.Args().Get(c.Args().Len()-1))
	return nil
}

func recreateContainer(c *cli.Context) error {
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

    err := operationContainer(c.Args().Get(c.Args().Len()-1), c.Args().First(), "recreate", c.Bool("wait"))
    if err != nil {
        fmt.Println(err)
    }

	fmt.Printf("Container %s recreated\n", c.Args().Get(c.Args().Len()-1))
	return nil
}



func getLogs(c *cli.Context) error {
    if c.Args().Len() != 2 {
        fmt.Println("Please provide a host and a container name")
        cli.ShowAppHelpAndExit(c, 1)
    }

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
