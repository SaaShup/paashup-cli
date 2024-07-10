package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"github.com/goccy/go-yaml"
    "github.com/SaaShup/paashup-sdk/docker"
    "github.com/SaaShup/paashup-sdk/netbox"
    "io/ioutil"
    "time"
    "strings"
)

type stackCompose struct {
    Version  string `json:"version" yaml:"version"`
    Services map[string]struct {
        Registry string `json:"registry" yaml:"registry"`
        Image string `json:"image" yaml:"image"`
        Ports []docker.Port `json:"ports" yaml:"ports"`
        Labels []docker.Label `json:"labels" yaml:"labels"`
        Env []docker.Env `json:"env" yaml:"env"`
        Hostname string `json:"hostname" yaml:"hostname"`
        Volumes map[string]struct {
            Source string `json:"source" yaml:"source"`
        } `json:"volumes" yaml:"volumes"`
        Networks []string `json:"networks" yaml:"networks"`
        Restart_policy string `json:"restart_policy" yaml:"restart_policy"`
    } `json:"services" yaml:"services"`
    Networks map[string]struct {
        Hostname string `json:"hostname" yaml:"hostname"`
        Driver string `json:"driver" yaml:"driver"`
    } `json:"networks" yaml:"networks"`
    Volumes map[string]struct {
        Hostname string `json:"hostname" yaml:"hostname"`
        Driver string `json:"driver" yaml:"driver"`
    } `json:"volumes" yaml:"volumes"`
    Registry map[string]struct {
        ServerAddress string `json:"serveraddress" yaml:"serveraddress"`
        Username string `json:"username" yaml:"username"`
        Password string `json:"password" yaml:"password"`
        Hostname string `json:"hostname" yaml:"hostname"`
    } `json:"registry" yaml:"registry"`
}

func stackDeployRun(c *cli.Context, compose stackCompose) error {
    fmt.Println("Deploying stack...")

    config, _ := readConfig(c)
    netbox.NETBOX_URL = config.URL
    netbox.NETBOX_TOKEN = config.Token

    for name, registry := range compose.Registry {
        host, err := docker.HostSearchByName(registry.Hostname)
        if err != nil {
            log.Fatal(err)
        }
        resp, err := docker.RegistrySearchByName(name, host.Id)
        if err != nil {
            createRegistryStruct := docker.RegistryCreateStruct{Name: name, ServerAddress: registry.ServerAddress, Username: registry.Username, Password: registry.Password, Host: host.Id}
            createRegistry, err := docker.RegistryCreate(createRegistryStruct)
            if err != nil {
                log.Fatal(fmt.Sprintf("Failed to create registry %s", name))
            }
            fmt.Printf("Registry %s created\n", createRegistry.Name)
        } else {
            fmt.Printf("Registry %s found\n", resp.Name)
        }
    }

    for name, network := range compose.Networks {
        host, err := docker.HostSearchByName(network.Hostname)
        if err != nil {
            log.Fatal(err)
        }
        resp, err := docker.NetworkSearchByName(name, host.Id)
        if err != nil {
            createNetworkStruct := docker.NetworkCreateStruct{Name: name, Driver: network.Driver, Host: host.Id, State: "creating", NetworkId: "0"}
            createNetwork, err := docker.NetworkCreate(createNetworkStruct)
            if err != nil {
                log.Fatal(fmt.Sprintf("Failed to create network %s", name))
            }
            fmt.Printf("Network %s created\n", createNetwork.Name)
        } else {
            fmt.Printf("Network %s found\n", resp.Name)
        }
    }

    for name, volume := range compose.Volumes {
        host, err := docker.HostSearchByName(volume.Hostname)
        if err != nil {
            log.Fatal(err)
        }
        resp, err := docker.VolumeSearchByName(name, host.Id)
        if err != nil {
            createVolumeStruct := docker.VolumeCreateStruct{Name: name, Driver: volume.Driver, Host: host.Id}
            createVolume, err := docker.VolumeCreate(createVolumeStruct)
            if err != nil {
                log.Fatal(fmt.Sprintf("Failed to create volume %s", name))
            }
            fmt.Printf("Volume %s created\n", createVolume.Name)
        } else {
            fmt.Printf("Volume %s found\n", resp.Name)
        }
    }

    for name, config := range compose.Services {
        host, err := docker.HostSearchByName(config.Hostname)
        if err != nil {
            log.Fatal(err)
        }
        registry, err := docker.RegistrySearchByName(config.Registry, host.Id)

        if err != nil {
            log.Fatal(err,)
        }

        var imageCreated docker.ImageComplete
        s := strings.Split(config.Image, ":")
        imageVersion := s[len(s)-1]
        imageName := strings.Replace(config.Image, ":" + imageVersion, "", -1)

        image, err := docker.ImageSearchByNameAndVersion(imageName, imageVersion, host.Id)
        if err != nil {
            createImageStruct := docker.ImageCreateStruct{Name: imageName, Host: host.Id, Registry: registry.Id, ImageID: "0", Version: imageVersion}
            imageCreated, err = docker.ImageCreate(createImageStruct)
            if err != nil {
                cli.Exit(err, 1)
            }
            for i := 0; i < 30; i++ {
                res, err := docker.ImageInspect(imageCreated.Id)
                if err != nil {
                    log.Fatal(err)
                }
                if res.Size > 0 {
                    break
                }
                time.Sleep(15 * time.Second)
            }
            res, _ := docker.ImageInspect(imageCreated.Id)
            if res.Size == 0 {
                log.Fatal(fmt.Sprintf("Failed to create image %s", imageCreated.Name))
            }
            fmt.Printf("Image %s created\n", imageCreated.Name)
            image = imageCreated
        } else {
            fmt.Printf("Image %s found\n", image.Name)
        }

        createContainerStruct := docker.ContainerCreateStruct{Name: name, Image: image.Id, Host: host.Id, State: "none", Operation: "create"}
        if config.Restart_policy == "" {
            createContainerStruct.Restart_policy = "no"
        } else {
            createContainerStruct.Restart_policy = config.Restart_policy
        }
        if len(config.Labels) == 0 {
            createContainerStruct.Labels = []docker.Label{}
        } else {
            createContainerStruct.Labels = config.Labels
        }
        if len(config.Env) == 0 {
            createContainerStruct.Env = []docker.Env{}
        } else {
            createContainerStruct.Env = config.Env
        }
        if len(config.Ports) == 0 {
            createContainerStruct.Ports = []docker.Port{}
        } else {
            createContainerStruct.Ports = config.Ports
        }
        if len(config.Volumes) == 0 {
            createContainerStruct.Mounts = []docker.Mount{}
        } else {
            for name, volume := range config.Volumes {
                volumeSearch, err := docker.VolumeSearchByName(name, host.Id)
                if err != nil {
                    cli.Exit(err, 1)
                }
                createContainerStruct.Mounts = append(createContainerStruct.Mounts, docker.Mount{Source: volume.Source, Volume: docker.Volume{Name: volumeSearch.Name, Driver: volumeSearch.Driver}})
            }
        }
        if len(config.Networks) == 0 {
            createContainerStruct.NetworkSettings = []docker.Network{}
        } else {
            for _, name := range config.Networks {
                networkSearch, err := docker.NetworkSearchByName(name, host.Id)
                if err != nil {
                    cli.Exit(err, 1)
                }
                createContainerStruct.NetworkSettings = append(createContainerStruct.NetworkSettings, docker.Network{Name: networkSearch.Name, Driver: networkSearch.Driver})
            }
        }
        container, err := docker.ContainerCreate(createContainerStruct)
        if err != nil {
            fmt.Printf("Failed to create container %s\n", name)
            log.Fatal(err)
        }
        for i := 0; i < 30; i++ {
            res, err := docker.ContainerInspect(container.Id)
            if err != nil {
                log.Fatal(err)
            }
            if res.State == "created" && res.Operation == "none" {
                break
            }
            time.Sleep(15 * time.Second)
        }
        fmt.Printf("Container %s created\n", container.Name) 


        _, err = docker.ContainerStart(docker.Container{Id: container.Id})
        if err != nil {
            log.Fatal(err)
        }

        for i := 0; i < 30; i++ {
            res, err := docker.ContainerInspect(container.Id)
            if err != nil {
                log.Fatal(err)
            }
            if res.State == "running" {
                break
            }
            if res.State == "exited" {
                log.Fatal(fmt.Sprintf("Container %s exited", container.Name))
            }
            time.Sleep(15 * time.Second)
        }

        fmt.Printf("Container %s started\n", container.Name)
    }
    fmt.Println("Stack deployed")
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

    err = stackDeployRun(c, compose)
    if err != nil {
        log.Fatal(err)
    }

    return nil
}
