<p align="center">
  <img width="160" src="docs/logo.svg" alt="logo">
</p>

<h1 align="center"><code>hckctl</code></h1>

<p align="center">
  <a href="https://github.com/hckops/hckctl/actions/workflows/ci.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/hckops/hckctl/ci.yml?label=ci&style=flat-square" alt="ci">
  </a>
  <a href="https://github.com/hckops/hckctl/actions/workflows/release.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/hckops/hckctl/release.yml?label=release&style=flat-square" alt="release">
  </a>
  <a href="https://pkg.go.dev/github.com/hckops/hckctl">
    <img src="https://pkg.go.dev/badge/github.com/hckops/hckctl.svg" alt="go-reference">
  </a>
  <a href="https://discord.gg/4PYPV9qP27">
    <img src="https://img.shields.io/badge/discord-join?label=join&logo=discord&style=flat-square&color=5865F2&logoColor=white" alt="discord">
  </a>
</p>

<p align="center">
  <i>The declarative Breach and Attack Simulation toolkit: one tool to rule 'em all</i><br>
  <a href="#quick-start">Quick start</a>&nbsp;&bull;
  <a href="#provider">Provider</a>&nbsp;&bull;
  <a href="#setup">Setup</a>&nbsp;&bull;
  <a href="#development">Development</a>&nbsp;&bull;
  <a href="#roadmap">Roadmap</a>&nbsp;&bull;
  <a href="#contribute">Contribute</a>
</p>
<br>

Launch manual and automated attacks with pre-defined and always up-to-date templates of your favourite tools.

Attack your vulnerable target infrastructure or connect to your training platform ([HTB](https://www.hackthebox.com), [TryHackMe](https://tryhackme.com), [Vulnlab](https://www.vulnlab.com), etc.) without wasting anymore time on boring installations, environment setup and network configurations.

Package, distribute and run known exploits to find weaknesses on authorized targets in a declarative way.

Designed to transparently run locally, remotely or integrated in pipelines and with guaranteed stability and backward compatibility over time.
`hckctl` is free, open source and community driven, no vendor lock-in, extensible and built using native providers api.

Leverage the cloud platform or request a dedicated managed environment to:
* orchestrate complex attack scenarios
* constantly probe and monitor your security posture
* analyze, aggregate and export results via api
* trigger instant actions based on observed events and patterns

## Quick start

### Box

Spin-up a [`box`](https://github.com/hckops/megalopolis/tree/main/box) and access all port-forwarded ports locally
```bash
# spawns a temporary docker box locally
hckctl box alpine
#[box-alpine-<RANDOM>][tty] tunnel (remote) 7681 -> (local) 7681
#[box-alpine-<RANDOM>] TTYD_USERNAME=root
#[box-alpine-<RANDOM>] TTYD_PASSWORD=alpine

# deploys a detached box to a kubernetes cluster
hckctl box start arch --provider kube
# tunnels tty port only
hckctl box open box-arch-<RANDOM> --no-exec

# creates a pwnbox box connected to your hack the box account
hckctl box preview/parrot-sec --network-vpn htb
# connects to vnc
vncviewer localhost:5900

# starts a background box to attack locally
hckctl box start vulnerable/owasp-juice-shop
```

### Lab (preview)

> TODO video

Access your target from a managed [`lab`](https://github.com/hckops/megalopolis/tree/main/lab) to
* tunnel multiple vpn connections through a highly available ssh proxy
* expose public endpoints with custom domains
* mount and keep in sync `dumps` e.g. git, s3
* load secrets from a vault
* save/restore workdir snapshots
* deploy private templates and infrastructures e.g. [Kompose](https://kompose.io), [Helm](https://helm.sh)
```bash
# starts demo lab (cloud only)
hckctl lab ctf-linux
```

### Task

Run a single-stage [`task`](https://github.com/hckops/megalopolis/tree/main/task) using pre-defined commands
```bash
# shows the "help" command
hckctl task nuclei --command help

# uses the "default" preset command and arguments
hckctl task rustscan
# equivalent of
hckctl task rustscan --input address=127.0.0.1
hckctl task scanner/rustscan --command default --input address=127.0.0.1

# runs the "full" preset command against the retired "Lame" machine (with docker)
# see https://app.hackthebox.com/machines/Lame
hckctl task nmap --network-vpn htb --command full --input address=10.10.10.3 
# equivalent of (with kube)
hckctl task nmap --network-vpn htb --provider kube --inline -- nmap 10.10.10.3 -sC -sV

# downloads common wordlists
git clone --depth 1 https://github.com/danielmiessler/SecLists.git \
  ${HOME}/.local/state/hck/share/wordlists/SecLists
# fuzzing loading a local template against the retired "Knife" machine (with kube)
# see https://app.hackthebox.com/machines/Knife
hckctl task \
  --local ../megalopolis/task/fuzzer/gobuster.yml \
  --network-vpn htb \
  --provider kube \
  --input address=10.10.10.242 \
  --input wordlist=wordlists/SecLists/Discovery/Web-Content/Apache.fuzz.txt

# monitors the logs
tail -F ${HOME}/.local/state/hck/task/log/task-*
```

### Flow (preview)

Run multistage tasks in parallel, collect and output the combined results
```bash
hckctl flow scan www.example.com
hckctl flow fuzz 127.0.0.1:8080
hckctl flow sql 127.0.0.1:3306
hckctl flow atomic-red-team 127.0.0.1 T1485
hckctl flow c2 ping
hckctl flow phishing @example.com
```

### Template

Explore all available templates or write your own and validate it locally
```bash
# lists all templates
hckctl template list

# validates all templates
hckctl template validate "../megalopolis/**/*.{yml,yaml}"
```

Inspired by [GitOps](https://www.gitops.tech), the whole project is centered around git as source of truth, schema validation and versioning. Pin a `revision` (branch, tag, or sha) if you need to ensure long term stability
```bash
# uses template "megalopolis/task/scanner/trivy" @ commit hash "12e7599"
hckctl task trivy --revision v0.1.0
```

### Config

Override the default configurations
```bash
# prints path and current configs
hckctl config

# resets default configs
hckctl config --reset
```

How to configure vpn networks
```bash
# edits config file
vim ${HOME}/.config/hck/config.yml

# example
network:
  vpn:
  - name: htb
    path: /home/demo/ctf/openvpn/htb_demo_eu_vip_28.ovpn
  - name: thm
    path: /home/demo/ctf/openvpn/thm_demo_us_regular_3.ovpn
```

## Provider

### Docker

Follow the official [instructions](https://docs.docker.com/engine/install) to install Docker Engine. The fastest way to get started is with the [convenience script](https://get.docker.com)
```bash
# downloads and runs script
curl -fsSL https://get.docker.com -o get-docker.sh
./sudo sh get-docker.sh
```

[lazydocker](https://github.com/jesseduffield/lazydocker) is the recommended tool to watch and monitor containers

### Kubernetes

#### Remote

If you are looking for a simple and cheap way to get started with a *remote* cluster use [kube-template](https://github.com/hckops/kube-template) on [DigitalOcean](https://www.digitalocean.com/products/kubernetes)
```bash
provider:
  kube:
    configPath: "/PATH/TO/kube-template/clusters/do-template-kubeconfig.yaml"
```

#### Local

Use [minikube](https://minikube.sigs.k8s.io), [kind](https://kind.sigs.k8s.io) or [k3s](https://k3s.io) to setup a local cluster
```bash
provider:
  kube:
    # absolute path, empty by default uses "${HOME}/.kube/config"
    configPath: ""
    namespace: hckops
```

#### Troubleshooting

Useful dev tools, see [`hckops/kube-base`](https://github.com/hckops/actions/blob/main/docker/Dockerfile.base)
```bash
# starts tmp container
docker run --rm --name hck-tmp-local --network host -it \
  -v ${HOME}/.kube/config:/root/.kube/config hckops/kube-base

# watches pods
kubectl klock -n hckops pods
```

Depending on your local environment, you might need to override IPv6 config in the *local* cluster to use the `--network-vpn` flag. Set also `--embed-certs` if you need to use the dev tools
```bash
# starts local cluster
minikube start --embed-certs \
  --extra-config="kubelet.allowed-unsafe-sysctls=net.ipv6.conf.all.disable_ipv6"

# runs with temporary privileges to connect to a vpn
env HCK_CONFIG_NETWORK.PRIVILEGED=true hckctl box alpine --provider kube --network-vpn htb
# equivalent of
network:
  # default is false, override for local clusters
  privileged: true
```

### Cloud

Access to the platform is limited and in ***private preview***. If you are interested, please leave a comment or a :thumbsup: to this [issue](https://github.com/hckops/hckctl/issues/104) and we'll reach out with more details
```bash
provider:
  cloud:
    host: <ADDRESS>
    port: 2222
    username: <USERNAME>
    token: <TOKEN>
```

### Podman (coming soon)

Follow the official [instructions](https://podman.io/docs/installation) to install Podman

## Setup

Download the latest binaries
```bash
HCKCTL_VERSION=0.12.0

# install or update
curl -sSL https://github.com/hckops/hckctl/releases/latest/download/hckctl-${HCKCTL_VERSION}-linux-x86_64.tar.gz | \
  sudo tar -xzf - -C /usr/local/bin

# verify
hckctl version

# uninstall
sudo rm /usr/local/bin/hckctl
```

## Development

* [just](https://github.com/casey/just)

```bash
# run
go run internal/main.go

# debug
go run internal/main.go task test/debug --provider kube --inline -- tree /hck/share

# build
just
./build/hckctl

# logs
tail -F ${HOME}/.local/state/hck/log/hckctl-*.log

# publish (without "v" prefix)
just publish <MAJOR.MINOR.PATCH>
```

## Roadmap

* `machine` create and access VMs e.g. DigitalOcean Droplet, AWS EC2, Azure Virtual Machines, QEMU etc.
* `tui` similar to lazydocker and k9s together
* `network` support WireGuard, Tor, ProxyChains, etc.
* `plugin` add custom cli commands in any language
  - `man` combine tldr and cheat with task commands
  - `htb` and `thm` api to start/stop/list machines and submit flags
  - `prompt` chatgpt prompt style

## Contribute

Create your custom template and test it locally
```bash
# loads local template
hckctl box --local ../megalopolis/box/preview/powershell.yml
```

Please, feel free to contribute to the companion [repository](https://github.com/hckops/megalopolis) and add more community templates to the catalog.
Credit should go to all the authors and maintainers for their open source tools, without them this project wouldn't exist!

<!--

box remote kube: after killing vnc/portforward
E1020 19:55:12.436966  149063 portforward.go:381] error copying from remote stream to local connection: readfrom tcp4 127.0.0.1:5900->127.0.0.1:54768: write tcp4 127.0.0.1:5900->127.0.0.1:54768: write: broken pipe

>>> remove TryHackMe demo from readme

* fix cloud
* update platform prs
* verify network connectivity between boxes/tasks i.e. kube.svc
* add task cloud (kube provider)
* use public PKG
* lab inputs
* convert TODOs left in GitHub issues
* add GitHub org labels: feature/bug/question

* test all catalog
* discord + social links (?)
* verify binaries
* test on mac-m1 and win (docker images)
* review all command cli example/description

TODO demo
* solve the machine and add how to after docker https://github.com/juice-shop/juice-shop#docker-container
* auto-exploitation box
* metasploit plugin
* windows examples

>>> lab + kompose https://github.com/kubernetes/kompose
composeRef e.g. https://github.com/digininja/DVWA/blob/master/compose.yml

TODO
* priority
    - debug `htb-postman`
    - add flow example
    - verify cloud no-shell support
    - play htb: linux/win
    - RELEASE example https://github.com/boz/kail#homebrew
    - docker release and gh-action
    - add copyTo/copyFrom box/task
* general
    - strict schema validation
    - add disclaimer of responsibility to readme?
    - public discord server (review channels visibility)
    - brew release
    - review context/http/client timeouts e.g. vpn or target not available
    - verify config migration between versions
    - add readme lab video/gif https://asciinema.org
    - delete old branches (video)
    - update internal cli diagram
    - review/delete GitHub project
    - add go reference badge
    - public `preview/kali-core` image
    - create PR to external official doc to run
        * owasp/dvwa
        * https://github.com/vulhub/vulhub
        * https://houdini.secsi.io
    - flaky tests (?)
        * kubernetes_test.go:TestNewResources
    - rename `template` to catalog? or alias?
    - cmd aliases e.g. start/up/create
* cli
    - autocomplete commands and values
        * e.g. `box connect <list of boxes>` with `ValidArgsFunction`
        * e.g. `box <list of box templates>` with `ValidArgsFunction`
        * see fix autocomplete
    - config add set command
    - add confirmation before
        * reset config
        * delete all
* template
    - add `--remote` mutually exclusive flag
    - keep up-to-date directories to exclude in `resolvePath` e.g. charts
    - add filters and review output e.g. table
* box
    - print/event shared directory, same as envs, ports etc.
    - review tty resize
    - expose copy from/to ???
    - kube: verify if `close()` is needed or `return nil`
    - kube: `execBox` deployment always check/scale replica to 1 before exec (test with replica=0)
    - kube: update resources sizes + comparison
    - docker: COPY shared volume `XDG_DATA_HOME`
    - docker: support powershell `/usr/bin/pwsh` (attach with no tty and raw terminal) see `docker run --rm -it mcr.microsoft.com/powershell`
    - docker: add support for remote docker daemon with `DOCKER_HOST`
    - add podman provider
    - add context timeout
    - cloud: ssh key auth only + remove InsecureIgnoreHostKey
    - cloud: remove body from empty request `omitempty to remove "body":{}`
    - list boxes in table with padding see `tabwriter` https://gosamples.dev/string-padding
    - filter/list box (list and delete) and template (list and validate) columns by provider + sorting
    - flaky issue `zerolog: could not write event: write /home/<REDACTED>/.local/state/hck/log/hckctl-<REDACTED>.log: file already closed`
* lab 
    - `command` cli description and example
    - in `create` add override e.g. `--input alias=parrot --input password=changeme --input vpn=htb-eu`
    - inputs should look for HCK_LAB_??? env var override if --input is not present before using default
    - verify optional merge/overrides
    - in `connect` merge/expand BoxEnv actual BoxEnv e.g. generated password
    - compose/template/infra
        * https://github.com/SpecterOps/BloodHound/blob/main/examples/docker-compose/docker-compose.yml
        * https://kompose.io
        * https://github.com/vulhub/vulhub
        * https://github.com/madhuakula/kubernetes-goat.git
* task
    - inputs should look for HCK_TASK_??? env var override if --input is not present before using default
    - review TaskV1 schema i.e. `pages`, `license`, command `description`
    - `history` command to list old tasks i.e. names of log files e.g. <TIMESTAMP>-task-<NAME>-<RANDOM>
    - for debug purposes prepend file output with interpolated task (yaml) or command parameters
    - add command to remove all logs
    - skip output file for `help` and `version`
    - limit default kube resources
    - add `--background` to omit stdout and ignore interrupt handler i.e. only file output
* version
    - print if new version available
    - implement server and providers `version` in json format docker/kube/cloud
* release
    - add brew https://goreleaser.com/customization/homebrew
    - test linux
    - test mac and mac1
    - test window vm
    - verify release workflow should depend on ci workflow
* prompt
    - https://github.com/snwfdhmp/awesome-gpt-prompt-engineering
* megalopolis
    - (docker) https://github.com/edoardottt/scilla

-->
