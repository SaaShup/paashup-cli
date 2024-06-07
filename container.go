package main

import (
    "fmt"
    "log"
    "encoding/json"
    "strings"
    "github.com/urfave/cli/v2"
)

func searchContainer(c *cli.Context, h HostComplete, containerName string) (Container, error) {
    for _, container := range h.Containers {
        if container.Name == containerName {
            return container, nil
        }
    }
    return Container{}, fmt.Errorf("Container not found")
}

func operationContainer(c *cli.Context, operation string) (Container, error) {

    host, err := searchHost(c)
    if err != nil {
        return Container{}, fmt.Errorf("Host not found")
    }

    container, err := searchContainer(c, host, c.Args().Get(c.Args().Len()-1))
    if err != nil {
        fmt.Println("Container not found")
        return Container{}, fmt.Errorf("Container not found")
    }

    url := fmt.Sprintf("containers/%d/", container.Id)
    var jsonStr = []byte(fmt.Sprintf(`{"operation":"%s"}`, operation))

    resultCall, err := netboxCall(c, url, "POST", jsonStr)

    if err != nil {
        log.Fatal(err)
    }

    var result Container

    if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    return result, nil
}

type command struct {
    Cmd []string `json:"cmd"`
}

func execContainer(c *cli.Context) error {

    host, err := searchHost(c)
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
        fmt.Println("Container not found")
        return nil
    }
    if container.Operation == "stop" {
        fmt.Println("Container " + container.Name + " stopped")
    } else {
        fmt.Println("Container " + container.Name + " not stopped")
    }
    return nil
}

func startContainer(c *cli.Context) error {
    container, err := operationContainer(c, "start")
    if err != nil {
        fmt.Println("Container not found")
        return nil
    }
    if container.Operation == "start" {
        fmt.Println("Container " + container.Name + " started")
    } else {
        fmt.Println("Container " + container.Name + " not started")
    }
    return nil
}

func getLogs(c *cli.Context) error {

    host, err := searchHost(c)
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

func listContainers(c *cli.Context) error {

    var url string
    if c.String("host") == "" {
        url = fmt.Sprintf("containers/")
    } else {
        var host, err = searchHost(c)
        if err != nil {
            fmt.Println("Host not found")
            return nil
        }
        url = fmt.Sprintf("containers/?host_id=%d", host.Id)
    }
    resultCall, err := netboxCall(c, url, "GET", nil)

    if err != nil {
        log.Fatal(err)
    }

    var result ContainerList

    if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    for _, rec := range result.Results {
        fmt.Println(rec.Name)
    }

    return nil
}

