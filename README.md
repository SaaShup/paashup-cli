# paashup-cli

`paashup-cli` is a command-line tool to manage your paashup environment. It provides various commands to interact with containers, hosts, and images.

## Installation

To install `paashup-cli`, ensure you have Go installed and run:

```bash
go get github.com/saashup/paashup-cli
```

## Usage

### Global Flags

- `--netbox-url, -N`: Netbox URL (Required, can be set via `NETBOX_URL` environment variable)
- `--netbox-token, -T`: Netbox Token (Required, can be set via `NETBOX_TOKEN` environment variable)
- `--format, -f`: Choose between `yaml` and `json` (Default: `json`)

### Commands

#### Docker

The `docker` command allows you to manage Docker containers, hosts, and images. It has several subcommands:

##### Container

- **list**: List all containers, optionally for a specific host.
  ```sh
  paashup-cli docker container list
  ```

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

- List all containers:
  ```sh
  paashup-cli docker container list
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

`paashup-cli` is released under the BSD 3 License. See `LICENSE` for more information.

---

This README provides a general overview of `paashup-cli` and how to use it. For more detailed documentation and examples, please refer to the official documentation.
