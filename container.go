package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "bytes"
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
    netboxUrl := strings.TrimRight(c.String("netbox-url"), "/")
    client := &http.Client{}

    host, err := searchHost(c)
    if err != nil {
        return Container{}, fmt.Errorf("Host not found")
    }

    container, err := searchContainer(c, host, c.Args().Get(c.Args().Len()-1))
    if err != nil {
        fmt.Println("Container not found")
        return Container{}, fmt.Errorf("Container not found")
    }

    var url = fmt.Sprintf("%s/api/plugins/docker/containers/%d/", netboxUrl, container.Id)
    var jsonStr = []byte(fmt.Sprintf(`{"operation":"%s"}`, operation))
    req, err := http.NewRequest("PATCH", url, ioutil.NopCloser(bytes.NewBuffer(jsonStr)))

    if err != nil {
        log.Fatal(err)
    }
    req.ContentLength = int64(len(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.String("netbox-token")))
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    defer res.Body.Close()

    var result Container

    b, err := ioutil.ReadAll(res.Body)

    if err := json.Unmarshal(b, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    return result, nil
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
    netboxUrl := strings.TrimRight(c.String("netbox-url"), "/")
    client := &http.Client{}

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

    var url = fmt.Sprintf("%s/api/plugins/docker/containers/%d/logs/", netboxUrl, container.Id)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.String("netbox-token")))
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    defer res.Body.Close()

    b, err := ioutil.ReadAll(res.Body)

    fmt.Printf("%s\n", b)

    return nil
}

func listContainers(c *cli.Context) error {
    netboxUrl := strings.TrimRight(c.String("netbox-url"), "/")
    client := &http.Client{}

    var url string
    if c.String("host") == "" {
        url = fmt.Sprintf("%s/api/plugins/docker/containers/", netboxUrl)
    } else {
        var host, err = searchHost(c)
        if err != nil {
            fmt.Println("Host not found")
            return nil
        }
        url = fmt.Sprintf("%s/api/plugins/docker/containers/?host_id=%d", netboxUrl, host.Id)
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.String("netbox-token")))
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    defer res.Body.Close()
    var result ContainerList

    b, err := ioutil.ReadAll(res.Body)

    if err := json.Unmarshal(b, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    for _, rec := range result.Results {
        fmt.Println(rec.Name)
    }

    return nil
}

