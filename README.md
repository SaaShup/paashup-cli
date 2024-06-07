# paashup-cli

`paashup-cli` is a command-line tool to manage your paashup environment. It provides various commands to interact with containers, hosts, and images.

## Installation

To install `paashup-cli`, ensure you have Go installed and run:

```bash
go get github.com/yourusername/paashup-cli
```

## Usage

`paashup-cli` provides several commands and flags to interact with your paashup environment.

### Global Flags

- `--host, -H`: Host to connect to (can also be set using `HOST` environment variable).
- `--netbox-url, -N`: Netbox URL (required, can also be set using `NETBOX_URL` environment variable).
- `--netbox-token, -T`: Netbox Token (required, can also be set using `NETBOX_TOKEN` environment variable).

### Commands

#### Container

Manage containers with the following subcommands:

- `list`: List all containers. To list containers on a specific host, use the `--host` flag.

  ```bash
  paashup-cli container list [--host HOST]
  ```

- `logs`: Get logs of a container.

  ```bash
  paashup-cli container logs [CONTAINER_NAME]
  ```

- `start`: Start a container.

  ```bash
  paashup-cli container start [CONTAINER_NAME] [--nowait, -w]
  ```

- `stop`: Stop a container.

  ```bash
  paashup-cli container stop [CONTAINER_NAME] [--nowait, -w]
  ```

- `exec`: Execute a command in a container.

  ```bash
  paashup-cli container exec [CONTAINER_NAME] [COMMAND]
  ```

#### Host

Manage hosts with the following subcommands:

- `list`: List all hosts.

  ```bash
  paashup-cli host list
  ```

#### Image

Manage images with the following subcommands:

- `list`: List all images.

  ```bash
  paashup-cli image list
  ```

## Examples

### List All Containers

```bash
paashup-cli container list
```

### List Containers on a Specific Host

```bash
paashup-cli container list --host example.com
```

### Get Logs of a Container

```bash
paashup-cli container logs my_container
```

### Start a Container

```bash
paashup-cli container start my_container --nowait
```

### Stop a Container

```bash
paashup-cli container stop my_container --nowait
```

### Execute a Command in a Container

```bash
paashup-cli container exec my_container ls -la
```

### List All Hosts

```bash
paashup-cli host list
```

### List All Images

```bash
paashup-cli image list
```

## Environment Variables

You can set the following environment variables instead of using flags:

- `HOST`: Host to connect to.
- `NETBOX_URL`: Netbox URL.
- `NETBOX_TOKEN`: Netbox Token.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the LICENSE file for details.# paashup-cli
