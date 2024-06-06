package main

type Host struct {
    Id int `json:"id"`
    Url string `json:"url"`
    Display string `json:"display"`
    Name string `json:"name"`
    State string `json:"state"`
    Agent_version string `json:"agent_version"`
    Docker_api_version string `json:"docker_api_version"`
    Endpoint string `json:"endpoint"`
}

type Image struct {
    Id int `json:"id"` 
    Url string `json:"url"`
    Display string `json:"display"`
    Name string `json:"name"`
    Version string `json:"version"`
    Size int `json:"size"`
    ImageID string `json:"ImageID"`
    Digest string `json:"Digest"`
}

type Volume struct {
    Id int `json:"id"`
    Url string `json:"url"`
    Display string `json:"display"`
    Name string `json:"name"`
    Driver string `json:"driver"`
}

type Network struct {
    Id int `json:"id"`
    Url string `json:"url"`
    Display string `json:"display"`
    Name string `json:"name"`
    Driver string `json:"driver"`
    NetworkID string `json:"NetworkID"`
    State string `json:"state"`
}

type Registry struct {
    Id int `json:"id"`
    Url string `json:"url"`
    Display string `json:"display"`
    Name string `json:"name"`
    ServerAddress string `json:"serveraddress"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email string `json:"email"`
}

type Container struct {
    Id int `json:"id"`
    Url string `json:"url"`
    Display string `json:"display"`
    Name string `json:"name"`
    ContainerID string `json:"ContainerID"`
    State string `json:"state"`
    Status string `json:"status"`
    Restart_policy string `json:"restart_policy"`
    Operation string `json:"operation"`
    Hostname string `json:"hostname"`
}

type ContainerList struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []struct {
        Container
        Ports []struct {
            Public_port int `json:"public_port"`
            Private_port int `json:"private_port"`
            Type string `json:"type"`
        }
        Env []struct {
            Var_name string `json:"var_name"`
            Value string `json:"value"`
        }
        Labels []struct {
            Key string `json:"key"`
            Value string `json:"value"`
        }
        Mounts []struct {
            Source string `json:"source"`
            Read_only bool `json:"read_only"`
            Volume Volume `json:"volume"`
        }
        Binds []struct {
            Host_path string `json:"host_path"`
            Container_path string `json:"container_path"`
            Read_only bool `json:"read_only"`
        }
        Network_settings []struct {
            Network Network `json:"network"`
        }
        Created string `json:"created"`
        Custom_fields interface{} `json:"custom_fields"`

        Last_updated string `json:"last_updated"`
        Tags []string `json:"tags"`

        Host Host `json:"host"`
        Image Image `json:"image"`
    } `json:"results"` 
}

type HostComplete struct {
    Token struct {
        Id int `json:"id"`
        Url string `json:"url"`
        Display string `json:"display"`
        Key string `json:"key"`
        Write_enabled bool `json:"write_enabled"`
    }
    Netbox_base_url string `json:"netbox_base_url"`
    Custom_fields interface{} `json:"custom_fields"`
    Last_updated string `json:"last_updated"`
    Tags []string `json:"tags"`
    Images []Image `json:"images"`
    Volumes []Volume `json:"volumes"`
    Networks []Network `json:"networks"`
    Containers []Container `json:"containers"`
    Registries []Registry `json:"registries"`
    Host
}

type HostList struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []HostComplete `json:"results"`
}
