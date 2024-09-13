package main

import (
    "github.com/supabase-community/supabase-go"
    "github.com/supabase-community/gotrue-go/types"
    "github.com/urfave/cli/v2"
    "fmt"
    "log"
    "encoding/json"
    "io/ioutil"
    "os"
    "net/http"
)

func platformCreateAccount(c *cli.Context) error {

    if c.Args().Len() != 2 {
        fmt.Println("Please provide an email and password")
        cli.ShowAppHelpAndExit(c, 1)
    }

    // Create a new Supabase client
    client, err := supabase.NewClient(PLATFORM_URL, PLATFORM_PUB_KEY, nil)

    if err != nil {
        log.Fatal("Could not connect to platform!")  
    }
    // Create a new accout
    
    signuUpRequest := types.SignupRequest{
        Email: c.Args().First(),
        Password: c.Args().Get(c.Args().Len()-1),
    }
    _, err = client.Auth.Signup(signuUpRequest)

    if err != nil {
        fmt.Println(err)
        log.Fatal("Could not create account!")  
    }

    fmt.Printf("Account created successfully! Your user is %s\n", c.Args().First())
    return nil
}

func platformLogout(c *cli.Context) error {
    var configpath string

    if os.Getenv("XDG_CONFIG_HOME") == "" {
        configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
    } else {
        configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
    }

    if _, err := os.Stat(configpath + "platform.token"); os.IsNotExist(err) {
        log.Fatal("You are not logged in!")
    }

    os.Remove(configpath + "platform.token")

    fmt.Println("Logged out successfully!")
    return nil
}

func platformReadLogin(c *cli.Context) (string, error) {
    var configpath string

    if os.Getenv("XDG_CONFIG_HOME") == "" {
        configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
    } else {
        configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
    }

    data, err := ioutil.ReadFile(configpath + "platform.token")

    if err != nil {
        return "", err
    }

    return string(data), nil
}

type PlatformList struct {
    Data []struct {
        Id int `json:"id"`
        Created_at string `json:"created_at"`
        Name string `json:"name"`
        User_id string `json:"user_id"`
    } `json:"data"` 
}

func platformList(c *cli.Context) error {
    token, err := platformReadLogin(c)
    
    if err != nil {
        log.Fatal("You are not logged in!")
    }
    client := &http.Client{} 
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/functions/v1/paashup-list", PLATFORM_URL), nil)

	if err != nil {
		log.Fatal("could not create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
    var data PlatformList
    body, _ := ioutil.ReadAll(res.Body)
  	if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to the go struct pointer
        return err
	}
    
    for _, project := range data.Data {
        fmt.Printf("https://%s.paashup.cloud\n", project.Name)
    }
    return nil
}

func platformLogin(c *cli.Context) error {
    client, _ := supabase.NewClient(PLATFORM_URL, PLATFORM_PUB_KEY, nil)

     if c.Args().Len() != 2 {
        fmt.Println("Please provide a username and password")
        cli.ShowAppHelpAndExit(c, 1)
    }
    
    data, err := client.SignInWithEmailPassword(c.Args().First(),  c.Args().Get(c.Args().Len()-1))

    if err != nil {
        log.Fatal("Could not login to the platform!")
    }

    var configpath string

    if os.Getenv("XDG_CONFIG_HOME") == "" {
		configpath = os.Getenv("HOME") + "/.config/paashup-cli/"
	} else {
		configpath = os.Getenv("XDG_CONFIG_HOME") + "/paashup-cli/"
	}

	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		os.MkdirAll(configpath, 0755)
	}
	if _, err := os.Stat(configpath + "platform.token"); os.IsNotExist(err) {
		os.Create(configpath + "platform.token")
	}

	_ = ioutil.WriteFile(configpath+"platform.token", []byte(data.AccessToken), 0644)

    fmt.Println("Logged in successfully!")

    return nil
}

func platformInit(c *cli.Context) error {
    token, err := platformReadLogin(c)
    
    if err != nil {
        log.Fatal("You are not logged in!")
    }
    client := &http.Client{} 
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/functions/v1/paashup-init", PLATFORM_URL), nil)

	if err != nil {
		log.Fatal("could not create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
    var data PlatformList
    body, _ := ioutil.ReadAll(res.Body)
  	if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to the go struct pointer
        return err
	}
    
    for _, project := range data.Data {
        fmt.Printf("https://%s.paashup.cloud\n", project.Name)
    }
    return nil

}
