package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/urfave/cli/v2"
)

func listHosts(c *cli.Context) error {
    client := &http.Client{}

    var url = c.String("netbox-url") + "/api/plugins/docker/hosts/"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Token " + c.String("netbox-token"))
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    defer res.Body.Close()
    var result HostList

    b, err := ioutil.ReadAll(res.Body)

    if err := json.Unmarshal(b, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    for _, rec := range result.Results {
        fmt.Println(rec.Name)
    }

    return nil
}

func searchHost(c *cli.Context) (HostComplete, error) {
    client := &http.Client{}

    if c.String("host") == "" {
        return HostComplete{}, fmt.Errorf("Host not found")
    }
    var url = c.String("netbox-url") + "/api/plugins/docker/hosts/?name=" + c.String("host")
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Token " + c.String("netbox-token"))
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    defer res.Body.Close()
    var result HostList

    b, err := ioutil.ReadAll(res.Body)

    if err := json.Unmarshal(b, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    if result.Count == 1 {
        return result.Results[0], nil
    } else {
        return HostComplete{}, fmt.Errorf("Host not found")
    }

}

