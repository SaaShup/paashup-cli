package main

import (
    "github.com/supabase-community/supabase-go"
    "github.com/supabase-community/gotrue-go/types"
    "github.com/urfave/cli/v2"
    "fmt"
    "log"
    "encoding/json"
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

func platformInit(c *cli.Context) error {

    if c.Args().Len() != 2 {
        fmt.Println("Please provide a username and password")
        cli.ShowAppHelpAndExit(c, 1)
    }

    client, err := supabase.NewClient(PLATFORM_URL, PLATFORM_PUB_KEY, nil)

    if err != nil {
        log.Fatal("Could not connect to platform!")
    }

    client.SignInWithEmailPassword(c.Args().First(),  c.Args().Get(c.Args().Len()-1))
    
    data, _, err := client.From("netbox").Insert(map[string]interface{}{}, true, "", "", "exact").Execute()

    if err != nil {
        log.Fatal("Could not init project")
    }

    var response []struct {
        Id int `json:"id"`
        Created_at string `json:"created_at"`
        Name string `json:"name"`
        User_id string `json:"user_id"`
    }

    if err := json.Unmarshal(data, &response); err != nil { // Parse []byte to the go struct pointer
        log.Fatal(err)
    }

    fmt.Printf("Your paashup url is https://%s.paashup.com\n", response[0].Name)

    return nil
}
