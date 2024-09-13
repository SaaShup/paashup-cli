# paashup-cli

`paashup-cli` is a command-line tool for managing your paashup environment, including Docker containers, hosts, images, and Netbox configurations. This tool provides various commands to list, inspect, start, stop, and execute operations on Docker containers, as well as managing Docker hosts, images, and Netbox configurations.

## Build
using docker:
```bash
docker run -it --rm -v ./:/go golang:1.22 go build -buildvcs=false .
```

## Installation

To install `paashup-cli`, ensure you have Go installed and run:

```bash
go get github.com/saashup/paashup-cli
```

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

`paashup-cli` provides multiple commands categorized under platforms, Docker resources, Netbox configurations, and stack deployments.

### Global Flags

- `--format, -f`: Choose between `yaml`, `json`, or `json-pretty` formats. Default is `json`.

### Commands Overview

The CLI is structured with several main commands, each with its own set of subcommands.

#### 1. Platform Commands

Manage platforms within your paashup environment.

- **Create an Account**

  ```bash
  paashup-cli platform account create <username> <password>
  ```

  Creates a new account. Example:

  ```bash
  paashup-cli platform account create user1 password123
  ```

- **Login in Platform**
  ```bash
  paashup-cli platform login <username> <password>
  ```

- **Initialize a Platform**

  ```bash
  paashup-cli platform init
  ```

  Initializes a platform. Example:

  ```bash
  paashup-cli platform init
  ```

- **List Platforms**

  ```bash
  paashup-cli platform ls
  ```
#### 2. Netbox Commands

Manage Netbox configurations for your paashup environment.

- **Set Netbox Configuration**

  ```bash
  paashup-cli netbox set-config NAME NETBOX_URL NETBOX_TOKEN
  ```

  Sets a Netbox configuration. Example:

  ```bash
  paashup-cli netbox set-config myconfig http://netbox.example.com token123
  ```

- **Use Netbox Configuration**

  ```bash
  paashup-cli netbox use NAME
  ```

  Selects a Netbox configuration by name. Example:

  ```bash
  paashup-cli netbox use myconfig
  ```

#### 3. Stack Commands

Manage and deploy stacks.

- **Deploy a Stack**

  ```bash
  paashup-cli stack deploy YAMLFILE
  ```

  Deploys a stack from a YAML file. Example:

  ```bash
  paashup-cli stack deploy stack.yaml
  ```

#### 4. Docker Commands

Manage Docker resources including containers, hosts, images, volumes, and registries.

- **Docker Container Commands**

  - **List Containers**

    ```bash
    paashup-cli docker container ls [hostname]
    ```

    Lists all containers. Example:

    ```bash
    paashup-cli docker container ls myhost
    ```

  - **Get Logs**

    ```bash
    paashup-cli docker container logs <hostname> <containername>
    ```

    Fetches logs for a specific container. Example:

    ```bash
    paashup-cli docker container logs myhost mycontainer
    ```

  - **Start a Container**

    ```bash
    paashup-cli docker container start [--nowait] <hostname> <containername>
    ```

    Starts a specific container. Example:

    ```bash
    paashup-cli docker container start --nowait myhost mycontainer
    ```

  - **Stop a Container**

    ```bash
    paashup-cli docker container stop [--nowait] <hostname> <containername>
    ```

    Stops a specific container. Example:

    ```bash
    paashup-cli docker container stop --nowait myhost mycontainer
    ```

  - **Inspect a Container**

    ```bash
    paashup-cli docker container inspect <hostname> <containername>
    ```

    Inspects a specific container. Example:

    ```bash
    paashup-cli docker container inspect myhost mycontainer
    ```

  - **Execute Command in a Container**

    ```bash
    paashup-cli docker container exec <hostname> <containername> '<command>'
    ```

    Executes a command within a specific container. Example:

    ```bash
    paashup-cli docker container exec myhost mycontainer 'ls -la'
    ```

- **Docker Host Commands**

  - **List Hosts**

    ```bash
    paashup-cli docker host ls
    ```

    Lists all hosts. Example:

    ```bash
    paashup-cli docker host ls
    ```

  - **Inspect a Host**

    ```bash
    paashup-cli docker host inspect <hostname>
    ```

    Inspects a specific host. Example:

    ```bash
    paashup-cli docker host inspect myhost
    ```

- **Docker Image Commands**

  - **List Images**

    ```bash
    paashup-cli docker image ls
    ```

    Lists all images. Example:

    ```bash
    paashup-cli docker image ls
    ```

- **Docker Volume Commands**

  - **List Volumes**

    ```bash
    paashup-cli docker volume ls
    ```

    Lists all volumes. Example:

    ```bash
    paashup-cli docker volume ls
    ```

- **Docker Registry Commands**

  - **List Registries**

    ```bash
    paashup-cli docker registry ls
    ```

    Lists all registries. Example:

    ```bash
    paashup-cli docker registry ls
    ```

## Example Workflows

### Set Up a New Environment

1. **Set Netbox Configuration**

   ```bash
   paashup-cli netbox set-config myconfig http://netbox.example.com token123
   ```

2. **Initialize a Platform**

   ```bash
   paashup-cli platform init
   ```

3. **Deploy a Stack**

   ```bash
   paashup-cli stack deploy stack.yaml
   ```

### Manage Docker Containers

1. **List All Containers**

   ```bash
   paashup-cli docker container ls
   ```

2. **Start a Container**

   ```bash
   paashup-cli docker container start --nowait myhost mycontainer
   ```

3. **Get Logs from a Container**

   ```bash
   paashup-cli docker container logs myhost mycontainer
   ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
