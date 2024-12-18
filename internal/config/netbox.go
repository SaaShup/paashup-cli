package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
    "fmt"
)

type NetboxConfig struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

func SetConfig(name, url, token string) error {
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
                netboxConfig[i].URL = url
                netboxConfig[i].Token = token
                netboxConfig[i].Name = name
			}
		}
		netboxConfig = append(netboxConfig, NetboxConfig{Name: name, URL: url, Token: token})
	}

	file, _ := json.MarshalIndent(netboxConfig, "", " ")
	_ = ioutil.WriteFile(configpath+"netbox.json", file, 0644)
	_ = ioutil.WriteFile(configpath+"current", []byte(name), 0644)
    fmt.Println("Config set to " + name)
	return nil
}

func ReadConfig() (NetboxConfig, error){
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
    return NetboxConfig{}, nil
}
