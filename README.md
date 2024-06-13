# paashup-cli

`paashup-cli` is a command-line tool to manage your paashup environment. It provides various commands to interact with containers, hosts, and images.

## Installation

To install `paashup-cli`, ensure you have Go installed and run:

```bash
go get github.com/saashup/paashup-cli
```

# paashup-cli

`paashup-cli` is a command-line tool for managing your paashup environment, including Docker containers, hosts, images, and Netbox configurations. This tool provides various commands to list, inspect, start, stop, and execute operations on Docker containers, as well as managing Docker hosts, images, and Netbox configurations.

## Setup
First you need to define the netbox source:
```sh 
paashup-cli netbox set-config name http://netbox.local mytoken
```

For autocompletion, download the completion script (autocompletion.bash) and then source it like this:
```sh 
PROG=paashup-cli source autocompletion.bash

```

## Usage

### Global Flags

- `--format, -f`: Choose between `yaml` and `json` (Default: `json`)

### Commands

#### Netbox

The `netbox` command allows you to manage Netbox configurations. It has several subcommands:

- **set-config**: Set Netbox configuration.
  ```sh
  paashup-cli netbox set-config NAME NETBOX_URL NETBOX_TOKEN
  ```

- **use**: Use a specific Netbox configuration.
  ```sh
  paashup-cli netbox use NAME
  ```

#### Docker

The `docker` command allows you to manage Docker containers, hosts, and images. It has several subcommands:

##### Container

- **list**: List all containers, optionally for a specific host.
  ```sh
  paashup-cli docker container ls [hostname]
  ```
  - Aliases: `ps`, `ls`

- **logs**: Get logs of a container.
  ```sh
  paashup-cli docker container logs <hostname> <containername>
  ```

- **start**: Start a container. Optionally, wait for the container to start.
  ```sh
  paashup-cli docker container start [--nowait] <hostname> <containername>
  ```
  - `--nowait, -w`: Do not wait for the container to start.

- **stop**: Stop a container. Optionally, wait for the container to stop.
  ```sh
  paashup-cli docker container stop [--nowait] <hostname> <containername>
  ```
  - `--nowait, -w`: Do not wait for the container to stop.

- **inspect**: Inspect a container.
  ```sh
  paashup-cli docker container inspect <hostname> <containername>
  ```

- **exec**: Execute a command in a container.
  ```sh
  paashup-cli docker container exec <hostname> <containername> '<command>'
  ```

##### Host

- **list**: List all Docker hosts.
  ```sh
  paashup-cli docker host list
  ```

##### Image

- **list**: List all Docker images.
  ```sh
  paashup-cli docker image list
  ```

## Examples

Here are a few examples of how to use `paashup-cli`:

- Set a Netbox configuration:
  ```sh
  paashup-cli netbox set-config myconfig http://netbox.local mytoken
  ```

- Use a specific Netbox configuration:
  ```sh
  paashup-cli netbox use myconfig
  ```

- List all containers:
  ```sh
  paashup-cli docker container ls
  ```

- Get logs of a container:
  ```sh
  paashup-cli docker container logs myhostname mycontainer
  ```

- Start a container and wait for it to start:
  ```sh
  paashup-cli docker container start myhostname mycontainer
  ```

- Stop a container without waiting:
  ```sh
  paashup-cli docker container stop --nowait myhostname mycontainer
  ```

- Inspect a container:
  ```sh
  paashup-cli docker container inspect myhostname mycontainer
  ```

- Execute a command in a container:
  ```sh
  paashup-cli docker container exec myhostname mycontainer 'ls -la'
  ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

`paashup-cli` is released under the BSD 3-Clause License. See `LICENSE` for more information.

---

This README provides a general overview of `paashup-cli` and how to use it. For more detailed documentation and examples, please refer to the official documentation.
