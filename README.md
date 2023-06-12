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

A novel BAS engine with a declarative approach to launch manual and simulated attacks, either against a sandbox lab or your infrastructure. It leverages pre-defined and always up-to-date recipes to probe and verify your security posture, designed to be integrated in automated pipelines and to analyze, aggregate and export reports.

> TODO screenshot/gif

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

> TODO

## Development

* [just](https://github.com/casey/just)

```bash
# run
go run cmd/main.go

# build
just
```

TODO
* general
    - autocomplete
* template
    - add offline mode source revision
    - update directories to exclude in `resolvePath` e.g. charts
* box
    - kube replace resources with size
    - refactor Exec and wait condition to detach without remove
    - verify support for remote docker daemon with `DOCKER_HOST`
    - review logs and errors output
    - support distroless
    - add podman provider
    - add context timeout
    - cloud ssh key auth only + remove InsecureIgnoreHostKey
    - docker/kube `cp` + `XDG_DATA_HOME`
* config
    - add set command
    - add confirmation before reset
* version
    - print server/cloud
    - print if new version
    - auto update
    - rename fields i.e. commit vs version or print both
* release
    - add brew https://goreleaser.com/customization/homebrew
    - test windows
* plugins
    - man
