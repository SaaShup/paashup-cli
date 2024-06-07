package main

import (
    "github.com/urfave/cli/v2"
    "net/http"
    "strings"
    "fmt"
    "io/ioutil"
    "bytes"
    "log"
)

func netboxCall(c *cli.Context, endpoint string, method string, jsonStr []byte) ([]byte, error) {
    netboxUrl := strings.TrimRight(c.String("netbox-url"), "/")
    client := &http.Client{}

    url := fmt.Sprintf("%s/api/plugins/docker/%s", netboxUrl, endpoint)
    req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBuffer(jsonStr)))

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

   // ...
    return ioutil.ReadAll(res.Body)
}
