package main

import (
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "os"
    "encoding/json"
    "net/http"
    "strings"
    "fmt"
    "bytes"
    "log"

)

type NetboxConfig struct {
    Name  string `json:"name"`
    URL   string `json:"url"`
    Token string `json:"token"`
}

func setConfig(name, url, token string) error {
    var configpath string
    var netboxConfig []NetboxConfig
    if os.Getenv("XDG_CONFIG_HOME") == "" {
        configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
    } else {
        configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
    }

    if _, err := os.Stat(configpath); os.IsNotExist(err) {
        os.MkdirAll(configpath, 0755)
    }
    if _, err := os.Stat(configpath + "netbox.json"); os.IsNotExist(err) {
        os.Create(configpath + "netbox.json")
        netboxConfig = append(netboxConfig, NetboxConfig{Name: name, URL: url, Token: token})
    } else {
        file, _ := ioutil.ReadFile(configpath + "netbox.json")
        json.Unmarshal(file, &netboxConfig)
        for i, config := range netboxConfig {
            if config.Name == name {
                netboxConfig[i].Name = name
                netboxConfig[i].URL = url
                netboxConfig[i].Token = token
            }
        }
        netboxConfig = append(netboxConfig, NetboxConfig{Name: name, URL: url, Token: token})
    }

    file, _ := json.MarshalIndent(netboxConfig, "", " ")
    _ = ioutil.WriteFile(configpath+"netbox.json", file, 0644)
    _ = ioutil.WriteFile(configpath+"current", []byte(name), 0644)
    return nil
}

func useNetboxConfig(c *cli.Context) error {
    if c.NArg() != 1 {
        return nil
    }
    var configpath string

    if os.Getenv("XDG_CONFIG_HOME") == "" {
        configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
    } else {
        configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
    }

    _ = ioutil.WriteFile(configpath+"current", []byte(c.Args().First()), 0644)
    println("Using " + c.Args().First() + " config")
    return nil
}

func readConfig(c *cli.Context) (NetboxConfig, error) {
    var configpath string
    if os.Getenv("XDG_CONFIG_HOME") == "" {
        configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
    } else {
        configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
    }

    name, _ := ioutil.ReadFile(configpath + "current")
    file, _ := ioutil.ReadFile(configpath + "netbox.json")
    var netboxConfig []NetboxConfig
    json.Unmarshal(file, &netboxConfig)
    for _, config := range netboxConfig {
        if config.Name == string(name) {
            return config, nil
        }
    }
    cli.ShowAppHelpAndExit(c, 1)
    return NetboxConfig{}, nil
}

func setNetboxConfig(c *cli.Context) error {
    if c.NArg() != 3 {
        return cli.ShowCommandHelp(c, "set-config")
    }
    name := c.Args().Get(0)
    url := c.Args().Get(1)
    token := c.Args().Get(2)
    return setConfig(name, url, token)
}

func netboxCall(c *cli.Context, endpoint string, method string, jsonStr []byte) ([]byte, error) {
    netboxConfig, err := readConfig(c)
    if err != nil {
        log.Fatal(err)
        cli.ShowAppHelpAndExit(c, 1)
    }
    netboxUrl := strings.TrimRight(netboxConfig.URL, "/")
    client := &http.Client{}

    url := fmt.Sprintf("%s/api/plugins/docker/%s", netboxUrl, endpoint)
    req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBuffer(jsonStr)))

    if err != nil {
        log.Fatal(err)
    }

    req.ContentLength = int64(len(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Token %s", netboxConfig.Token))
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    defer res.Body.Close()

   // ...
    return ioutil.ReadAll(res.Body)
}
