package main

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/urfave/cli/v2"
)

func listHosts(c *cli.Context) error {
    url := "hosts/"

    var result HostList
    resultCall, err := netboxCall(c, url, "GET", nil)

    if err != nil {
        log.Fatal(err)
    }

    if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    for _, rec := range result.Results {
        fmt.Println(rec.Name)
    }

    return nil
}

func searchHost(c *cli.Context) (HostComplete, error) {

    if c.String("host") == "" {
        return HostComplete{}, fmt.Errorf("Host not found")
    }

    url := fmt.Sprintf("hosts/?name=%s", c.String("host"))
    var result HostList
    resultCall, err := netboxCall(c, url, "GET", nil)

    if err != nil {
        log.Fatal(err)
    }


    if err := json.Unmarshal(resultCall, &result); err != nil {  // Parse []byte to the go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }

    if result.Count == 1 {
        return result.Results[0], nil
    } else {
        return HostComplete{}, fmt.Errorf("Host not found")
    }

}

