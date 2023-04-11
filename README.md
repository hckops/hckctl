<p align="center">
  <img width="160" src="docs/logo.svg" alt="logo">
</p>

<h1 align="center"><code>hckctl</code></h1>

<p align="center">
  <a href="https://github.com/hckops/hckctl/actions/workflows/ci.yaml">
    <img src="https://github.com/hckops/hckctl/actions/workflows/ci.yaml/badge.svg" alt="ci">
  </a>
</p>

<p align="center">
  <i>The Cloud Native HaCKing Tool</i><br>
  <a href="#quick-start">Quick start</a>&nbsp;&bull;
  <a href="#setup">Setup</a>&nbsp;&bull;
  <a href="#guide">Guide</a>&nbsp;&bull;
  <a href="#development">Development</a>
</p>

<!--
A novel BAS tool with a declarative approach to launch manual and simulated attacks either against self-contained labs or your infrastructure. It uses pre-defined always up-to-date recipes to probe and verify your security posture, designed to be integrated in automated pipelines and with the possibility to analyze, aggregate and export reports.
-->

> TODO description and screenshot/gif

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

### Template

```bash
# validates and prints remote template
hckctl template parrot | yq -o=json

# validates and prints local template
hckctl template -p ../megalopolis/boxes/official/alpine.yml
```

### Box

<!--
**Boxes** are ready-to-go docker images designed for security enthusiasts that want to spend more time hacking and need both an attacker and a vulnerable environment that is constantly updated, quick to start and just work

Main features:
* unified local and remote experience - run the same environments locally or in a remote cluster
* open source and publicly maintained - you want to know what you are running!
  - see [templates](https://github.com/hckops/megalopolis/tree/main/boxes)
  - see [docker images](https://github.com/hckops/megalopolis/tree/main/docker)
* constantly updated
  - see scheduled [action](https://github.com/hckops/megalopolis/blob/main/.github/workflows/docker-ci.yml)
* all declared ports are exposed and forwarded by default
* resources are automatically deleted once you close a box
* *the cloud provider is not publicly available at this time*
-->

```bash
# lists boxes
hckctl box list

# starts a docker box (default)
hckctl box alpine
hckctl box alpine --provider docker

# starts a kubernetes box
hckctl box alpine --provider kube

# starts a remote box (over ssh tunnel)
hckctl box alpine --provider cloud
```

<!--
### Lab

> **Labs** are user-defined hacking environments

Main features:
* override defaults e.g. credentials, environment variables, etc.
* attach volumes
* connect multiple boxes

> WIP coming soon
-->

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
  # branch, tag or sha of https://github.com/hckops/megalopolis
  revision: main
  # docker|kube|cloud
  provider: docker
  kube:
    namespace: labs
    # absolute path, default "~/.kube/config"
    configPath: ""
    resources:
      memory: 512Mi
      cpu: 500m
  cloud:
    host: 0.0.0.0
    port: 2222
    username: ""
    token: ""
log:
  # debug|info|warning|error
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
* box: load from local path?
* box: fix validation + vulnerable path
* box: support distroless and different shell
* box: list from megalopolis (hardcoded)
* box: add detached mode + reconnect to existing + tunnel only
* box: test with podman
* box: add timeout
* box: refactor box/template shared cmd
* box: cloud ssh key auth only + remove InsecureIgnoreHostKey
* box: docker/kube `cp` + `XDG_DATA_HOME`
* box: move logs from `/tmp` to `XDG_STATE_HOME`
* schema: convert to valid CRD?
* schema: stable version + move in megalopolis?
* man plugin
* config: add set/reset cmd
* add cmd version
* release: verify archives + add brew
* `pkg/client` replace callback with channels
* `pkg/client` review: docker/kube methods
* cmd
  ```bash
  # client and server: docker|kube|cloud
  hckctl version
  
  # current config
  hckctl config
  
  # --provider docker|kube|cloud
  # open
  hckctl box <TEMPLATE_NAME> [--revision <REVISION>]
  hckctl box --path <TEMPLATE_PATH>
  
  # returns BOX_ID
  hckctl box create <TEMPLATE_NAME> [--revision <REVISION>]
  hckctl box exec <BOX_ID>
  hckctl box tunnel <BOX_ID>
  hckctl box delete <BOX_ID>
  # new
  hckctl box cp <PATH_FROM> <PATH_TO>
  
  # all boxes: docker/kube
  hckctl box list
  
  hckctl template <TEMPLATE_NAME> [--revision <REVISION>]
  hckctl template --path <TEMPLATE_PATH>
  hckctl template list [box]
  ```
