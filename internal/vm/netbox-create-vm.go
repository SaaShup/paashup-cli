package vm

import (
    "net/http"
    "io/ioutil"
    "bytes"
    "log"
    "fmt"
    "strings"
    "github.com/SaaShup/paashup-sdk/pkg/netbox"
    "encoding/json"
    "github.com/urfave/cli/v2"
    "github.com/SaaShup/paashup-cli/internal/config"
)

type NetboxVmResponse struct {
    Name string `json:"name"`
    Status struct {
        Value string `json:"value"`
    } `json:"status"`
    PrimaryIp struct {
        Address string `json:"address"`
    } `json:"primary_ip"`
}

func FindVm(c *cli.Context, name string) (NetboxVmResponse, error) {
    config, err := config.ReadConfig() 
	netboxUrl := strings.TrimRight(config.URL, "/")

    type FindVM struct {
        Count int `json:"count"`
        Results []NetboxVmResponse `json:"results"`
    }
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/virtualization/virtual-machines/?name=%s", netboxUrl, name), nil)

	if err != nil {
		return NetboxVmResponse{}, err
	}
	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", netbox.NETBOX_TOKEN))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

    result, err := ioutil.ReadAll(res.Body)
	var resp FindVM
    if err := json.Unmarshal(result, &resp); err != nil {
        return NetboxVmResponse{}, err
    }
    if len(resp.Results) == 0 {
        return NetboxVmResponse{}, fmt.Errorf("VM %s not found", name)
    }

	return resp.Results[0], nil

}

func CreateVm(c *cli.Context, name string) (NetboxVmResponse, error) {
    type ClusterVM struct {
        Name string `json:"name"`
    }

    type Vm struct {
        Name string `json:"name"`
        Status string `json:"status"`
        Cluster ClusterVM `json:"cluster"`
    }

    jsonStr, _ := json.Marshal(Vm{Name: name, Status: "planned", Cluster: ClusterVM{ Name: "saashup" }})
    config, err := config.ReadConfig()
	netboxUrl := strings.TrimRight(config.URL, "/")
	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/virtualization/virtual-machines/", netboxUrl), ioutil.NopCloser(bytes.NewBuffer(jsonStr)))

	if err != nil {
		return NetboxVmResponse{}, err
	}

	req.ContentLength = int64(len(jsonStr))
	req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", netbox.NETBOX_TOKEN))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

    result, err := ioutil.ReadAll(res.Body)
	var resp NetboxVmResponse
    if err := json.Unmarshal(result, &resp); err != nil {
        return resp, err
    }
	return resp, nil
}
