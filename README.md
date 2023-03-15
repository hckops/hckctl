<p align="center">
  <img width="160" src="docs/logo.svg" alt="logo">
</p>

<h1 align="center"><code>hckctl</code></h1>

<h3 align="center">The Cloud Native HaCKing Tool</h3>

<div align="center">
  <a href="https://github.com/hckops/hckctl/actions/workflows/ci.yaml">
    <img src="https://github.com/hckops/hckctl/actions/workflows/ci.yaml/badge.svg" alt="ci">
  </a>
</div>
<br>

> TODO screenshots

## Quick start

Create an `alpine` box to see how it works
```bash
# spawns a docker box locally
hckctl box alpine

# deploys a box to your kubernetes cluster
hckctl box alpine --provider kube
```

Spin-up a `parrot` box to start hacking!
```bash
# credentials: parrot|password
hckctl box parrot

# vnc
vncviewer localhost:5900

# (mac|linux) novnc
[open|xdg-open] http://localhost:6080

# (mac|linux) tty
[open|xdg-open] http://localhost:7681
```

## Setup

> TODO

```bash
# TODO latest
curl -sSL https://github.com/hckops/hckctl/releases/download/v0.1.0/hckctl_linux_x86_64.tar.gz | \
  tar -xzf - -C /usr/local/bin
```

## Guide

### Box

> **Boxes** are ready-to-go docker images designed for security enthusiasts that want to spend more time hacking and need both an attacker and a vulnerable environment that is constantly updated, quick to start and just work

Main features:
* unified local and remote experience - run the same environments locally or in a remote cluster
* open source and publicly maintained - you want to know what you are running!
  - see [templates](https://github.com/hckops/megalopolis/tree/main/boxes)
  - see [docker images](https://github.com/hckops/megalopolis/tree/main/docker)
* constantly updated
  - see scheduled [action](https://github.com/hckops/megalopolis/blob/main/.github/workflows/docker-ci.yml)
* all declared ports are exposed and forwarded by default
* resources are automatically deleted once you close the boxes
* the *cloud* provider over ssh tunnel is not publicly available at this time

```bash
# lists boxes
hckctl box list

# starts a docker box (default)
hckctl box alpine
hckctl box alpine --provider docker

# starts a kubernetes box
hckctl box alpine --provider kube

# starts a remote box
hckctl box alpine --provider cloud
```

### Lab

> **Labs** are user-defined hacking environments for your specific needs

Main features:
* override defaults e.g. credentials, environment variables, etc.
* attach volumes
* connect multiple boxes

> WIP coming soon

### Template

```bash
# validates and prints remote template
hckctl template parrot | yq -o=json

# validates and prints local template
hckctl template -p ../megalopolis/boxes/official/alpine.yml
```

### Config

```bash
# prints current config
hckctl config

# edits config file
vim ~/.config/hck/config.yml
```

Default
```yaml
kind: config/v1
box:
  revision: main
  provider: docker
  kube:
    namespace: labs
    configPath: ~/.kube/config
    resources:
      memory: 512Mi
      cpu: 500m
  cloud:
    host: 0.0.0.0
    port: 2222
    username: ""
    token: ""
log:
  level: info
  filePath: /tmp/hckctl-ubuntu.log
```

## Development

* [just](https://github.com/casey/just)

```bash
# run
go run main.go

# build
just

# debug
./build/hckctl <CMD> --log-level debug

# logs
tail -f /tmp/hckctl-*.log
```

TODO
* box: add detached mode + reconnect to existing + tunnel only
* box: list from megalopolis (hardcoded)
* box: test with podman
* finalize schema (move in megalopolis)
* add cmd version
* release: verify archives + add brew
