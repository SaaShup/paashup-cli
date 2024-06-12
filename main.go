package main

import (
    "log"
    "os"
    "github.com/urfave/cli/v2"
)

func main() {
    dockerSubcommands := []*cli.Command{
        {
            Name: "docker",
            Subcommands: []*cli.Command{
                {
                    Name: "container",
                    Subcommands: []*cli.Command{
                        {
                            Name: "ps",
                            Usage: "List all containers\nExample: paashup-cli docker container ps [-a] [hostname]\n",
                            Action: psContainers,
                            Flags: []cli.Flag{
                                &cli.BoolFlag{
                                    Name:    "all",
                                    Aliases: []string{"a"},
                                    Usage:   "Show all containers",
                                },
                            },
                        },
                        {
                            Name:  "list",
                            Usage: "List All containers, to list only on specific host proide the hostname\n",
                            Action: listContainers,
                        },
                        {
                            Name: "logs",
                            Usage: "Get logs of a container\nExample: paashup-cli docker container logs <hostname> <containername>\n",
                            Action: getLogs,
                        },
                        {
                            Name: "start",
                            Usage: "Start a container\nExample: paashup-cli docker container start [--nowait] <hostname ><containername>\n",
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
                            Usage: "Stop a container\nExample: paashup-cli docker container stop [--nowait] <hostname> <containername>\n",
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
                            Name: "inspect",
                            Usage: "inspect a container\nExample: paashup-cli docker container inspect <hostname> <containername>\n",
                            Action: inspectContainer,
                        },
                        {
                            Name: "exec",
                            Usage: "Execute a command in a container\nExample: paashup-cli docker container exec <hostname> <containername> '<command>'\n",
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
        },
    }
    app := &cli.App{
        Name:  "paashup-cli",
        Version: version,
        Usage: "Manage your paashup",
        Flags: []cli.Flag{
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
            &cli.StringFlag{
                Name:    "format",
                DefaultText: "json",
                Aliases: []string{"f"},
                Usage:   "choose between yaml, json Default: json",
            },
        },
        Commands: []*cli.Command{
            dockerSubcommands[0],
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
