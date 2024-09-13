package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
    platformSubcommands := []*cli.Command{
        { 
            Name: "platform",
            Subcommands: []*cli.Command{
                {
                    Name: "account",
                    Subcommands: []*cli.Command{
                        {
                            Name: "create",
                            Usage: "Create an account\nExample: paashup-cli account create <email> <password>\n",
                            Action: platformCreateAccount,
                        },
                    },
                },
                {
                    Name: "init",
                    Usage: "Initialize a platform\nExample: paashup-cli platform init <email> <password>\n",
                    Action: platformInit,
                },
                {
                    Name: "login",
                    Usage: "Login inside the platform\nExample: paashup-cli platform login <email> <password>\n",
                    Action: platformLogin,
                },
                {
                    Name: "list",
                    Aliases: []string{"ls", "ps"},
                    Usage: "List all platforms\nExample: paashup-cli platform list\n",
                    Action: platformList,
                },
            },
        },
    }
	netboxSubcommands := []*cli.Command{
		{
			Name: "netbox",
			Subcommands: []*cli.Command{
				{
					Name:   "set-config",
					Usage:  "Set Netbox Config\nExample: paashup-cli netbox set-config NAME NETBOX_URL NETBOX_TOKEN\n",
					Action: setNetboxConfig,
				},
				{
					Name:   "use",
					Usage:  "Use Netbox Config\nExample: paashup-cli netbox use NAME\n",
					Action: useNetboxConfig,
				},
			},
		},
	}
    stackSubcommands := []*cli.Command{
        {
            Name: "stack",
            Subcommands: []*cli.Command{
                {
                    Name: "deploy",
                    Usage: "Deploy a stack\nExample: paashup-cli stack deploy YAMLFILE\n",
                    Action: stackDeploy,
    				Flags: []cli.Flag{
	    				&cli.BoolFlag{
		    				Name:    "nowait",
							Aliases: []string{"w"},
							Usage:   "Wait for container to start",
						},
					},
                },
            },
        },
    }
	dockerSubcommands := []*cli.Command{
		{
			Name: "docker",
			Subcommands: []*cli.Command{
				{
					Name: "container",
					Subcommands: []*cli.Command{
						{
							Name:    "list",
							Usage:   "List all containers\nExample: paashup-cli docker container ls [hostname]\n",
							Aliases: []string{"ps", "ls"},
							Action:  psContainers,
						},
						{
							Name:   "logs",
							Usage:  "Get logs of a container\nExample: paashup-cli docker container logs <hostname> <containername>\n",
							Action: getLogs,
						},
						{
							Name:   "start",
							Usage:  "Start a container\nExample: paashup-cli docker container start [--nowait] <hostname ><containername>\n",
							Action: startContainer,
							Flags: []cli.Flag{
								&cli.BoolFlag{
									Name:    "detach",
									Aliases: []string{"d"},
									Usage:   "Run job in background",
								},
							},
						},
						{
							Name:   "stop",
							Usage:  "Stop a container\nExample: paashup-cli docker container stop [--nowait] <hostname> <containername>\n",
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
							Name:   "inspect",
							Usage:  "inspect a container\nExample: paashup-cli docker container inspect <hostname> <containername>\n",
							Action: inspectContainer,
						},
						{
							Name:   "exec",
							Usage:  "Execute a command in a container\nExample: paashup-cli docker container exec <hostname> <containername> '<command>'\n",
							Action: execContainer,
						},
					},
				},
				{
					Name: "host",
					Subcommands: []*cli.Command{
						{
							Name:   "list",
							Usage:  "List All Host",
                            Aliases: []string{"ls", "ps"},
							Action: listHosts,
						},
                        {
                            Name: "inspect",
                            Usage: "Inspect a host",
                            Action: inspectHost,
                        },
					},
				},
				{
					Name: "image",
					Subcommands: []*cli.Command{
						{
							Name:  "list",
							Usage: "List all images",
                            Aliases: []string{"ls", "ps"},
                            Action: listImages,
						},
					},
				},
				{
					Name: "volume",
					Subcommands: []*cli.Command{
						{
							Name:  "list",
							Usage: "List all Volumes",
                            Aliases: []string{"ls", "ps"},
                            Action: listVolumes,
						},
					},
				},
				{
					Name: "registry",
					Subcommands: []*cli.Command{
						{
							Name:  "list",
							Usage: "List all Registries",
                            Aliases: []string{"ls", "ps"},
                            Action: listRegistries,
						},
					},
				},

            },
		},
	}
	app := &cli.App{
		Name:                 "paashup-cli",
		Version:              version,
		Usage:                "Manage your paashup\nTo setup please use first paashup-cli netbox set-config\n",
		EnableBashCompletion: true,
		Suggest:              true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "choose between yaml, json, json-pretty Default: json",
			},
		},
		Commands: []*cli.Command{
            platformSubcommands[0],
			dockerSubcommands[0],
			netboxSubcommands[0],
            stackSubcommands[0],
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
