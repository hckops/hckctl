<h1 align="center"><code>hckctl</code></h1>

<div align="center">
  <a href="https://github.com/hckops/hckctl/actions/workflows/ci.yaml">
    <img src="https://github.com/hckops/hckctl/actions/workflows/ci.yaml/badge.svg" alt="ci">
  </a>
</div>
<br>

<p align="center">
  <img width="160" src="docs/logo.svg" alt="logo">
</p>

The Cloud Native HaCKing Tool

## Setup

```bash
# TODO latest
# TODO verify archives + add brew
curl -sSL https://github.com/hckops/hckctl/releases/download/v0.1.0/hckctl_linux_x86_64.tar.gz | \
  tar -xzf - -C /usr/local/bin
```

## Quick start

Create an `alpine` box to see how it works
```bash
# spawns a docker box locally
hckctl box alpine

# deploys a box in your kubernetes cluster
hckctl box alpine --provider kube
```

> TODO screenshots

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

## Commands

### Box

Boxes are community driven and publicly maintained docker images, designed for security enthusiasts that want to spend more time hacking and need an environment that is constantly updated, quick to start and just work.

Main features:
* public, you want to know what you are running!
  - see templates: TODO
  - see docker images: TODO
* constantly unpdated
  - see scheduled actions TODO
* unified local and remote experience
* all declared ports are exposed and forwarded by default
* resources are automatically deleted once you close the box

```bash
# lists boxes
hckctl box list

# starts a docker box (default)
hckctl box alpine
hckctl box alpine --provider docker

# starts a kubernetes box
hckctl box alpine --provider kube
```

> TODO discord

The *cloud* version is in private beta only, request more information in the discord channel!
```bash
# starts a remote box
hckctl box alpine --provider cloud
```

### Lab

Labs are user-defined customized hacking environments for your specific needs

Main features:
* override defaults e.g. credentials, environment variables, etc.
* attach volumes
* multiple connected boxes

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
# edits config file
vim ~/.config/hck/config.yml

# prints current config
hckctl config
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
* add tests
* finalize schema (move in megalopolis)
* add cmd version
