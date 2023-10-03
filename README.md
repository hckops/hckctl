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
<br>

Launch manual and automated attacks with pre-defined and always up-to-date templates of your favourite tools.

Designed to transparently run locally, remotely or integrated in pipelines, `hckctl` is free and open-source, no vendor lock-in, extensible and built using native providers api.

Create a custom vulnerable target (box with a specific CVE or whole infrastructures) or connect to your CTF platform ([HTB](https://www.hackthebox.com), [TryHackMe](https://tryhackme.com), [Vulnlab](https://www.vulnlab.com), etc.) without wasting anymore time on boring installations, environment setup or network configurations.

Leverage the cloud platform or request a dedicated managed cluster to:
* orchestrate complex scenarios
* monitor and observe your security posture
* analyze, aggregate and export results via api
* trigger actions based on your requirements

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

### Lab (preview)

> TODO video

Access your target from a personalized [`lab`](https://github.com/hckops/megalopolis/tree/main/lab)
```bash
# connects to a vpn, exposes public ports, mount dumps etc.
hckctl lab ctf-linux
```

#### HTB example

> TODO

### Task

Run a [`task`](https://github.com/hckops/megalopolis/tree/main/task) using pre-defined commands
```bash
# use the "default" arguments
hckctl task rustscan --input address=127.0.0.1
# equivalent of
hckctl task rustscan --command default --input address=127.0.0.1

# use the "full" preset arguments
hckctl task nmap --command full --input address=127.0.0.1 --input port=80

# invoke it with custom arguments
hckctl task rustscan --inline -- -a 127.0.0.1
```

#### HTB example

Prerequisites
* start the retired [Lame](https://app.hackthebox.com/machines/Lame) machine in your account
* edit your network vpn config
    ```bash
    vim ${HOME}/.config/hck/config.yml
    # update path
    network:
      vpn:
      - name: htb
        path: /home/demo/ctf/openvpn/htb_demo_eu_vip_28.ovpn
    ```

Run your tasks against the machine
```bash
# scan with nmap
hckctl task nmap --network-vpn htb --command full --input address=10.10.10.3

# scan with rustscan
hckctl task rustscan --network-vpn htb --inline -- -a 10.10.10.3 --ulimit 5000

# scan with nuclei
hckctl task nuclei --network-vpn htb --input target=10.10.10.3

# TODO ffuf
```

See [output](./docs/task-htb-example.txt) example

### Flow (preview)

Launch multiple tasks in parallel, collect and combine the results
```bash
hckctl flow scan www.example.com
hckctl flow fuzz 127.0.0.1:8080
hckctl flow sql 127.0.0.1:3306
hckctl flow atomic-red-team 127.0.0.1 T1485
hckctl flow c2 ping
hckctl flow campaign/phishing @example.com
```

### Template

Explore all available templates. Pin a git `revision` to ensure reliability in automated pipelines
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

## Roadmap

* `machine`: create and access VMs e.g. DigitalOcean Droplet, AWS EC2, Azure Virtual Machines, QEMU etc.
* `tui`: similar to lazydocker and k9s
* `plugin`: add custom cli commands in any language
  - `man`: combine tldr and cheat with task commands
  - `prompt`: chatgpt prompt style

## Setup

Download the latest binaries
```bash
# TODO latest
HCKCTL_VERSION=???

curl -sSL https://github.com/hckops/hckctl/releases/download/${HCKCTL_VERSION}/hckctl_linux_x86_64.tar.gz | \
  tar -xzf - -C /usr/local/bin
```

## Development

* [just](https://github.com/casey/just)

```bash
# run
go run internal/main.go

# build
just
./build/hckctl

tail -F ${HOME}/.local/state/hck/log/hckctl-*.log
```

## Contribute

> TODO example of how to point to a specific pr/revision in a forked repo

<!--

2) task share dir
3) box share dir
4) task kube/cloud
5) refactor box/lab network docker/kube/cloud
6) flow example

TODO
* priority
    - add box/lab --network-vpn support
    - add task providers: kube and cloud
    - add task shareDir volume or copy dir e.g. ffuf + seclists
    - save tasks logs to file in logDir + add header with commands/parameters
    - TODO update lab and task cli example/description
    - play htb: linux/win
    - add flow example
    - add context client timeout e.g. vpn or target not available
    - verify kube/cloud distroless support
    - verify kube/cloud no-shell support
    - RELEASE
* general
    - brew release
    - review client timeouts
    - verify config migration between versions
    - update readme
        * remove comments
        * update setup
        * descriptions/screenshot/gif
    - delete old branches (video)
    - disclaimer of responsibility
    - update internal cli diagram
    - convert TODOs left in GitHub issues
    - add GitHub org labels: feature/bug/question
    - review/delete GitHub project
    - add go reference badge
    - public `preview/kali-core` image
    - create PR to external official doc to run
        * owasp/dvwa
        * https://github.com/vulhub/vulhub
        * https://houdini.secsi.io
    - flaky tests (?)
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
    - add filters and review output e.g. table
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
* task
    - `history` command to list old tasks i.e. names of log files e.g. <TIMESTAMP>-task-<NAME>-<RANDOM>
    - rename output log file with timestamp?
    - prepend file output with task/command yaml?
    - add command to remove all logs
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
* prompt
    - https://github.com/snwfdhmp/awesome-gpt-prompt-engineering

-->
