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
  <i>The declarative Breach and Attack Simulation engine</i><br>
  <a href="#quick-start">Quick start</a>&nbsp;&bull;
  <a href="#setup">Setup</a>&nbsp;&bull;
  <a href="#development">Development</a>
</p>

Launch manual and automated attacks with pre-defined and always up-to-date templates of your favourite tools.
Designed to transparently run locally, remotely or integrated in pipelines, `hckctl` is free and open-source, no vendor lock-in, extensible and built using native Docker and Kubernetes api.
Create a custom vulnerable target or connect to your CTF platform without wasting anymore time on boring installations, environment setup or network configurations.

Leverage the managed cloud platform to access it from anywhere, orchestrate complex scenarios and analyze, aggregate and export your results.

## Quick start

### Box

Spin-up a [`box`](https://github.com/hckops/megalopolis/tree/main/box) and access all port-forwarded ports locally
```bash
# spawns a temporary docker (default) box locally
hckctl box alpine

# deploys an ephemeral box to your kubernetes cluster
hckctl box arch --provider kube

# creates a managed box
hckctl box parrot --provider cloud
```

### Task

Run a [`task`](https://github.com/hckops/megalopolis/tree/main/task) using pre-defined commands
```bash
# use the "default" arguments
hckctl task rustscan --input address=127.0.0.1
# equivalent of
hckctl task rustscan --command default --input address=127.0.0.1

# use the "full" arguments
hckctl task nmap --command full --input address=127.0.0.1 --input port=80

# use custom arguments
hckctl task rustscan --inline -- -a 127.0.0.1
```

Hack The Box example, start the `Lame` box
```bash
# TODO add vpn config

# run tasks
hckctl task nmap --network-vpn htb --command full --input address=10.10.10.3
hckctl task rustscan --network-vpn htb --inline -- -a 10.10.10.3 --ulimit 5000
hckctl task nuclei --network-vpn htb --input target=10.10.10.3
```

### Lab (preview)

Access your favourite platform ([HTB](https://www.hackthebox.com), [TryHackMe](https://tryhackme.com), [Vulnlab](https://www.vulnlab.com) etc.) from a personalized [`lab`](https://github.com/hckops/megalopolis/tree/main/lab)
```bash
# connects to a vpn, exposes public ports, mount dumps etc.
hckctl lab ctf-linux
```

### Flow (TODO)

Launch multiple tasks in parallel and combine the results
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

### Machine

> create and access AWS EC2, Azure Virtual Machines, DigitalOcean Droplet, QEMU etc.

### Man

> combine tldr and cheat

### Plugin

> add custom commands

### TUI

> see lazydocker

### Prompt

> chatgpt style https://github.com/snwfdhmp/awesome-gpt-prompt-engineering

### Template

Explore public templates. Pin a git `revision` to ensure reliability in a CI/CD pipeline
```bash
hckctl template list
```

Please, feel free to contribute to the companion [repository](https://github.com/hckops/megalopolis) and add more templates.

### Config

Edit default configurations
```bash
# vim ${HOME}/.config/hck/config.yml
# prints current configs
hckctl config

# resets default configs
hckctl config --reset
```

## Setup

Download the latest binaries
```bash
# TODO latest
HCKCTL_VERSION=???

curl -sSL https://github.com/hckops/hckctl/releases/download/${HCKCTL_VERSION}/hckctl_linux_x86_64.tar.gz | \
  tar -xzf - -C /usr/local/bin
```

## Contribute

> TODO example of how to point to a specific pr/revision in a forked repo

## Development

* [just](https://github.com/casey/just)

```bash
# run
go run cmd/main.go

# build
just
./build/hckctl
```

<!--
TODO
* priority
    - add task
    - update task schema and validation
    - play htb: linux/win
    - implement lab for docker (with kind) and kube
    - add flow with vpn
    - cloud size comparison
    - verify kube/cloud distroless support
    - verify kube/cloud no shell support
    - RELEASE
    - config migration between versions
* general
    - brew release
    - review client timeouts
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
    - flaky tests
        * kubernetes_test.go:TestNewResources
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
    - update directories to exclude in `resolvePath` e.g. charts
* box
    - review tty resize
    - implement copy ???
    - kube: add distroless support
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
    - flaky issue zerolog `could not write event: write /home/ubuntu/.local/state/hck/hckctl-ubuntu.log: file already closed`
* lab 
    - in `create` add override e.g. `--input alias=parrot --input password=changeme --input vpn=htb-eu`
    - verify optional merge/overrides
    - in `connect` merge/expand BoxEnv actual BoxEnv e.g. generated password
    - compose/template/infra
        * https://github.com/SpecterOps/BloodHound/blob/main/examples/docker-compose/docker-compose.yml
        * https://kompose.io
        * https://github.com/vulhub/vulhub
        * https://github.com/madhuakula/kubernetes-goat.git
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

-->
