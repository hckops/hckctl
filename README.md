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
  <a href="#development">Development</a>
</p>

<!--
TODO description/screenshot/video/gif

A novel Breach and Attack Simulation (BAS) engine with a declarative approach to launch manual and automated attacks, either against a sandbox lab or your infrastructure.
It leverages pre-defined and always up-to-date recipes of your everyday tools to probe and verify your security posture.
Designed to transparently run locally, remotely or integrated in pipelines and to analyze, aggregate and export reports.
-->

## Quick start

### Box

Create an [`alpine`](https://github.com/hckops/megalopolis/blob/main/box/base/alpine.yml) box to see how it works
```bash
# spawns a docker box locally on a shared network
hckctl box alpine

# deploys a box to your kubernetes cluster
hckctl box alpine --provider kube

# TODO add env credentials: alpine|changeme
# (mac|linux) tty
[open|xdg-open] http://localhost:7681
```

Spin-up a [`parrot`](https://github.com/hckops/megalopolis/blob/main/box/base/parrot.yml) box and access all port-forwarded ports locally to start hacking
```bash
# credentials: parrot|changeme
hckctl box parrot

# vnc
vncviewer localhost:5900

# (mac|linux) novnc
[open|xdg-open] http://localhost:6080
```

Attack your vulnerable [`dvwa`](https://github.com/hckops/megalopolis/blob/main/box/vulnerable/dvwa.yml) box or create your own
```bash
# TODO
hckctl box create dvwa
hckctl box start dvwa
hckctl box up dvwa # <<<

# TODO
hckctl box remove/delete dvwa
hckctl box stop dvwa
hckctl box down dvwa # <<<

# (mac|linux) web
[open|xdg-open] http://localhost:8080
```

*There is no difference between attacker or vulnerable boxes, if you can containerize it you can run it*

### Task

> TODO

```bash
# https://github.com/RustScan/RustScan/wiki/Installation-Guide#docker-whale
hckctl task rustscan ???

# TODO envs + args
```

<!--
### Lab

> Unleash the power of Kubernetes with GitOps to simulate whole infrastructures, for both red and blue teams

Easily start your remote htb-kali pwnbox connected to the [Hack The Box](https://www.hackthebox.com) VPN to sharpen your skills
```bash
# TODO create kube secret
# TODO htb-kali

# credentials: kali|changeme
hckctl lab htb-kali --provider argo
```

TODO
```bash
kube-goat
```

### Flow

> WIP

```bash
hckctl flow atomic-red-team T1485
hckctl flow scan 0.0.0.0
hckctl flow prowler
hckctl flow fuzz 0.0.0.0:8080/path
hckctl flow exploit/sql 0.0.0.0
hckctl flow tool/metasploit auxiliary/scanner/ssh/ssh_version
hckctl flow c2 ping
hckctl flow gen/pdf
hckctl flow campaign/phishing @example.com
hckctl flow api/virustotal/upload
hckctl flow scrape www.example.com
```
-->

### Template

Explore the public templates, and consider pinning a git `revision` to ensure reliability in a CI/CD pipeline
```bash
hckctl template list
```

Please, feel free to contribute to the companion [repository](https://github.com/hckops/megalopolis) and add more templates

## Setup

> TODO

```bash
# TODO latest
curl -sSL https://github.com/hckops/hckctl/releases/download/v0.1.0/hckctl_linux_x86_64.tar.gz | \
  tar -xzf - -C /usr/local/bin
```

Edit the config to override the defaults
```bash
hckctl config
```
<!--
If you are looking for a quick way to start with ArgoCD consider [kube-template](https://github.com/hckops/kube-template).
Just follow the readme, you'll be able to create and deploy a cluster on DigitalOcean using GitHub actions with literally a `git push`.
Once ready, update the `box.kube.configpath` config to use `clusters/do-template-kubeconfig.yaml`, that's all!
-->

## Development

* [just](https://github.com/casey/just)

```bash
# run
go run cmd/main.go

# build
just
./build/hckctl
```

TODO
* priority
    - cloud tunnel
    - finish cloud e.g. loader/event
    - server
        * refactor and uniform cmd/internal/pkg/config/logger etc.
        * update session/request
        * refactor providers
        * review charts
    - implement env
    - implement copy
    - mount vpn docker/kube
    - remove lab except for schema
    - delete internal folder
    - RELEASE
    - brew
    - add task
    - client timeouts
* general
    - update readme
        * remove comments
        * update setup
        * descriptions/screenshot/gif
    - add *guide*: all commands explained
    - add *example*: different uses cases e.g htb, etc
    - delete old branches (video)
    - disclaimer of responsibility
    - update internal cli diagram
    - convert TODOs left in GitHub issues
    - add GitHub org labels: feature/bug/question
    - review/delete GitHub project
    - highlight attacker and victim boxes to create specific scenario
    - add go reference badge
    - public `kali-core` image
    - PR to official doc to run
        * owasp/dvwa
        * https://github.com/vulhub/vulhub
        * https://houdini.secsi.io
* cli
    - autocomplete commands and values
        * e.g. `box connect <list of boxes>` with `ValidArgsFunction`
        * e.g. `box <list of box templates>` with `ValidArgsFunction`
        * see fix autocomplete
    - improve command validation e.g. docker `Args: cli.ExactArgs(1)`
    - filter/list box (list and delete) and template (list and validate) columns by provider + sorting
    - config add set command
    - add confirmation before
        * reset config
        * delete all
* template
    - add remote template
    - update directories to exclude in `resolvePath` e.g. charts
* box
    - review command: `create/remove` vs `start/stop` vs `up/down`
    - review command: `copy` vs `upload/download`
    - review tty resize
    - refactor common `BoxClient` methods (abstract)
    - mount `/dev/tun` for vpn
    - implement copy ???
    - kube: add env var
    - kube: add distroless support
    - kube: verify if close is needed or `return nil`
    - kube: `execBox` deployment always check/scale replica to 1 before exec (test with replica=0)
    - kube: `execBox` print envs
    - kube: verify `GetPodInfo` sidecar pod count
    - kube: verify `PodPortForward` callback vs channel
    - kube: update resources sizes
    - docker: `ContainerCreate` add env var
    - docker: `execBox` retrieve actual bound ports from container
    - docker: `execBox` print ports and envs
    - docker: `execBox` restart stopped/existed container
    - docker: `execBox` fix tunnelOnly (skip exec) and noTunnel (print warning ignored)
    - docker: COPY shared volume `XDG_DATA_HOME`
    - docker: support powershell `/usr/bin/pwsh` (attach with no tty and raw terminal) see `docker run --rm -it mcr.microsoft.com/powershell`
    - docker: add support for remote docker daemon with `DOCKER_HOST`
    - add podman provider
    - add context timeout
    - cloud: ssh key auth only + remove InsecureIgnoreHostKey
    - cloud: remove body from empty request `omitempty to remove "body":{}`
    - remove lab ??? are envs + copy enough?
    - list boxes in table with padding see `tabwriter` https://gosamples.dev/string-padding
    - template: rename `remote` to `git` and add `--remote` mutually exclusive flag to fetch via http
* task
    - TODO ???
* version
    - print if new version available
    - implement server `version` in json format docker/kube/cloud
* release
    - add brew https://goreleaser.com/customization/homebrew
    - test linux
    - test mac and mac1
    - test window vm
* plugins/bundles
    - man (plugin)
    - kube-inject (plugin) mount sidecar pod at runtime with debugging tools
    - pro (bundle) e.g. flow
