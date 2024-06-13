package main

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/goccy/go-yaml"
    "strings"
    "time"
    "github.com/fatih/color"
    "github.com/rodaine/table"
    "github.com/urfave/cli/v2"
    "github.com/mattn/go-runewidth"
)

func searchContainer(c *cli.Context, h HostComplete, containerName string) (Container, error) {
    for _, container := range h.Containers {
        if container.Name == containerName {
            return container, nil
        }
    }
    return Container{}, fmt.Errorf("Container not found")
}

type operationType struct {
    Operation string `json:"operation"`
}

type listContainerStruct struct {
    Id int `json:"id" yaml:"id"`
    Name string `json:"name" yaml:"name"`
    Host string `json:"host" yaml:"host"`
    Image string `json:"image" yaml:"image"`
    State string `json:"state" yaml:"state"`
    PortsCount int `json:"ports_count" yaml:"ports_count"`
    MountsCount int `json:"mounts_count" yaml:"mounts_count"`
    BindsCount int `json:"binds_count" yaml:"binds_count"`
    NetworksCount int `json:"networks_count" yaml:"networks_count"`
    EnvCount int `json:"env_count" yaml:"env_count"`
    LabelsCount int `json:"labels_count" yaml:"labels_count"`
}

func calcWidhtColorRed(s string) int {
    return runewidth.StringWidth(strings.Replace(strings.Replace(s, "\x1b[31m", "", 1), "\x1b[0m", "", 1))
}

func psContainers(c *cli.Context) error {
    hostname := c.Args().First()
    var url string
    if hostname != "" {
        host, err := searchHost(c.Args().First(), c)
        if err != nil {
            fmt.Println("Host not found")
            return nil
        }
        url = fmt.Sprintf("containers/?host_id=%d", host.Id)
    } else {
        url = fmt.Sprintf("containers/")
    }

    resultCall, err := netboxCall(c, url, "GET", nil)

    if err != nil {
        log.Fatal(err)
    }

    var result ContainerList

    if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
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
        switch c.String("format"){
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

func operationContainer(c *cli.Context, operation string) (Container, error) {

    host, err := searchHost(c.Args().First(), c)
    if err != nil {
        return Container{}, fmt.Errorf("Host not found")
    }

    container, err := searchContainer(c, host, c.Args().Get(c.Args().Len()-1))
    if err != nil {
        return Container{}, fmt.Errorf("Container not found")
    }

    url := fmt.Sprintf("containers/%d/", container.Id)
    operationS := &operationType{Operation: operation}
    jsonStr, _ := json.Marshal(operationS)

    resultCall, err := netboxCall(c, url, "PATCH", jsonStr)

    if err != nil {
        log.Fatal(err)
    }

    var result Container

    if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    if result.Operation == operation && !c.Bool("nowait") {
        i := 0
        for i < 20 {
            var resultFor Container
            resultCall, err = netboxCall(c, url, "GET", nil)
            if err != nil {
                log.Fatal(err)
            }
            if err := json.Unmarshal(resultCall, &resultFor); err != nil {  // Parse []byte to the go struct pointer
                fmt.Println("Can not unmarshal JSON")
            }
            switch operation {
                case "start":
                    if resultFor.State == "running" {
                        return result, nil
                    }
                case "stop":
                    if resultFor.State == "exited" {
                        return result, nil
                    }
            }
            time.Sleep(1 * time.Second)
            i++
        }
        return result, fmt.Errorf("Timeout")

    } else if result.Operation == operation {
        return result, nil
    } else {
        return result, fmt.Errorf("Operation not executed")
    }
    
}

type command struct {
    Cmd []string `json:"cmd"`
}

func inspectContainer(c *cli.Context) error {
    
        host, err := searchHost(c.Args().First(), c)
        if err != nil {
            fmt.Println("Host not found")
            return nil
        }
    
        container, err := searchContainer(c, host, c.Args().Get(c.Args().Len()-1))
        if err != nil {
            fmt.Println("Container not found")
            return nil
        }
    
        url := fmt.Sprintf("containers/%d/", container.Id)
        resultCall, err := netboxCall(c, url, "GET", nil)
    
        if err != nil {
            log.Fatal(err)
        }
        
        var result ContainerComplete
        if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
            fmt.Println("Can not unmarshal JSON")
        }
        switch c.String("format"){
            case "json-pretty":
                resp, err := json.MarshalIndent(result, "", "    ")
                if err == nil {
                    fmt.Printf("%s\n", resp)
                }
            case "json":
                resp, err := json.Marshal(result)
                if err == nil {
                    fmt.Printf("%s\n", resp)
                }
            case "yaml":
                resp, err := yaml.Marshal(result)
                if err == nil {
                    fmt.Printf("%s\n", resp)
                }
            default:
                resp, err := json.Marshal(result)
                if err == nil {
                    fmt.Printf("%s\n", resp)
                }
        }
    
        return nil
    
}

func execContainer(c *cli.Context) error {

    host, err := searchHost(c.Args().First(), c)
    if err != nil {
        fmt.Println("Host not found")
        return nil
    }

    container, err := searchContainer(c, host, c.Args().Get(c.Args().Len()-2))
    if err != nil {
        fmt.Println("Container not found")
        return nil
    }

    url := fmt.Sprintf("containers/%d/exec/", container.Id)
    command := &command{Cmd: strings.Fields(c.Args().Get(c.Args().Len()-1))}
    jsonStr, _ := json.Marshal(command)


    resultCall, err := netboxCall(c, url, "POST", jsonStr)

    if err != nil {
        log.Fatal(err)
    }

    var result Exec

    if err := json.Unmarshal(resultCall, &result); err != nil { 
        fmt.Println("Can not unmarshal JSON")
    }

    fmt.Println(result.Stdout)

    return nil

}

func stopContainer(c *cli.Context) error {
    container, err := operationContainer(c, "stop")
    if err != nil {
        fmt.Println(err)
        return nil
    }
    fmt.Println("Container " + container.Name + " stopped")
    return nil
}

func startContainer(c *cli.Context) error {
    container, err := operationContainer(c, "start")
    if err != nil {
        fmt.Println(err)
        return nil
    }
    fmt.Println("Container " + container.Name + " started")
    return nil
}

func getLogs(c *cli.Context) error {

    host, err := searchHost(c.Args().First(), c)
    if err != nil {
        fmt.Println("Host not found")
        return nil
    }

    container, err := searchContainer(c, host, c.Args().Get(c.Args().Len()-1))
    if err != nil {
        fmt.Println("Container not found")
        return nil
    }

    url := fmt.Sprintf("containers/%d/logs/", container.Id)
    resultCall, err := netboxCall(c, url, "GET", nil)

    if err != nil {
        log.Fatal(err)
    }


    fmt.Printf("%s\n", resultCall)

    return nil
}
