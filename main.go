package main

import (
    "log"
    "os"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "paashup-cli",
        Usage: "Manage your paashup",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "host",
                Aliases: []string{"H"},
                Usage:   "Host to connect to",
                EnvVars: []string{"HOST"},
            },
            &cli.StringFlag{
                Name:    "netbox-url",
                Aliases: []string{"N"},
                Usage:   "Netbox URL",
                Required: true,
                EnvVars: []string{"NETBOX_URL"},
            },
            &cli.StringFlag{
                Name:    "netbox-token",
                Aliases: []string{"T"},
                Usage:   "Netbox Token",
                Required: true,
                EnvVars: []string{"NETBOX_TOKEN"},
            },
        },
        Commands: []*cli.Command{
            {
                Name: "container",
                Subcommands: []*cli.Command{
                    {
                        Name:  "list",
                        Usage: "List All containers, to list only on specific host use --host",
                        Action: listContainers,
                    },
                    {
                        Name: "logs",
                        Usage: "Get logs of a container",
                        Action: getLogs,
                    },
                    {
                        Name: "start",
                        Usage: "Start a container",
                        Action: startContainer,
                        Flags: []cli.Flag{
                            &cli.BoolFlag{
                                Name:    "nowait",
                                Aliases: []string{"w"},
                                Usage:   "Wait for container to start",
                            },
                        },
                    },
                    {
                        Name: "stop",
                        Usage: "Stop a container",
                        Action: stopContainer,
                        Flags: []cli.Flag{
                            &cli.BoolFlag{
                                Name:    "nowait",
                                Aliases: []string{"w"},
                                Usage:   "Wait for container to start",
                            },
                        },

                    },
                    {
                        Name: "exec",
                        Usage: "Execute a command in a container",
                        Action: execContainer,
                    },
                },
            },
            {
                Name: "host",
                Subcommands: []*cli.Command{
                    {
                        Name:  "list",
                        Usage: "List All Host",
                        Action: listHosts,
                    },
                },
            },
            {
                Name: "image",
                Subcommands: []*cli.Command{
                    {
                        Name:  "list",
                        Usage: "List all images",
                    },
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
