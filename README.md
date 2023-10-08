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
  <i>The declarative Breach and Attack Simulation tool</i><br>
  <a href="#quick-start">Quick start</a>&nbsp;&bull;
  <a href="#setup">Setup</a>&nbsp;&bull;
  <a href="#provider">Provider</a>&nbsp;&bull;
  <a href="#development">Development</a>&nbsp;&bull;
  <a href="#roadmap">Roadmap</a>&nbsp;&bull;
  <a href="#contribute">Contribute</a>
</p>
<br>

Launch manual and automated attacks with pre-defined and always up-to-date templates of your favourite tools.

Designed to transparently run locally, remotely or integrated in pipelines. `hckctl` is free and open-source, no vendor lock-in, extensible and built using native providers api.

Attack your vulnerable target infrastructure or connect to your training platform ([HTB](https://www.hackthebox.com), [TryHackMe](https://tryhackme.com), [Vulnlab](https://www.vulnlab.com), etc.) without wasting anymore time on boring installations, environment setup and network configurations.

Leverage the cloud platform or request a dedicated managed cluster to:
* orchestrate complex attack scenarios
* constantly probe and monitor your security posture
* analyze, aggregate and export results via api
* trigger instant actions based on observed events and patterns

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

#### HTB demo

Prerequisites
* start the retired [Postman](https://app.hackthebox.com/machines/Postman) machine in your account
* edit your vpn network config
    ```bash
    vim ${HOME}/.config/hck/config.yml

    network:
      vpn:
      - name: htb
        # update with your openvpn config path
        path: /home/demo/ctf/openvpn/htb_demo_eu_vip_28.ovpn
    ```

Start your *pwnbox* and solve the challenges
```bash
# pulls a preview box (first time might take a while)
hckctl box preview/parrot-sec --network-vpn htb
```

Start an auto-exploitation box
```bash
# TODO review

# exploits the machine and spawns a reverse shell
hckctl box --network-vpn htb --local ../megalopolis/box/ctf/htb-postman.yml
```

### Lab (preview)

> TODO video

Access your target from a managed [`lab`](https://github.com/hckops/megalopolis/tree/main/lab)
```bash
# connects to a vpn, exposes public ports, mount dumps etc.
hckctl lab ctf-linux
```

### Task

Run a [`task`](https://github.com/hckops/megalopolis/tree/main/task) using pre-defined commands
```bash
# default commands
hckctl task gobuster --command help
hckctl task fuzzer/ffuf --command version

# use the "default" preset arguments
hckctl task rustscan
# equivalent of
hckctl task rustscan --input address=127.0.0.1
hckctl task scanner/rustscan --command default --input address=127.0.0.1

# use the "full" preset arguments
hckctl task nmap --command full --input address=127.0.0.1 --input port=80

# invoke it with custom arguments
hckctl task nuclei --inline -- -u https://example.com

# monitor the logs
tail -F ${HOME}/.local/state/hck/task/log/task-*
```

#### TryHackMe demo

> TODO

<!--
Prerequisites
* start the retired [Lame](https://app.hackthebox.com/machines/Lame) and [Knife](https://app.hackthebox.com/machines/Knife) machines in your account
* edit your vpn network config (see box example above)

Run tasks against the vulnerable machine
```bash
# scan with nmap
hckctl task nmap --network-vpn htb --command full --input address=10.10.10.3

# scan with rustscan
hckctl task rustscan --network-vpn htb --inline -- -a 10.10.10.3 --ulimit 5000

# scan with nuclei
hckctl task nuclei --network-vpn htb --input address=10.10.10.3
```
See [output](./docs/task-htb-example.txt) example

Use the shared directory to mount local paths
```bash
# download your wordlists
mkdir -p ${HOME}/.local/state/hck/share/wordlists
git clone --depth 1 https://github.com/danielmiessler/SecLists.git \
  ${HOME}/.local/state/hck/share/wordlists/SecLists

# fuzzing with ffuf
hckctl task ffuf --network-vpn htb --input address=10.10.10.242

# fuzzing with gobuster
hckctl task \
  --local ../megalopolis/task/fuzzer/gobuster.yml \
  --network-vpn htb \
  --input address=10.10.10.242 \
  --input wordlist=wordlists/SecLists/Discovery/Web-Content/Apache.fuzz.txt
```
-->

### Flow (preview)

Launch multiple tasks in parallel, collect and combine the results
```bash
hckctl flow scan www.example.com
hckctl flow fuzz 127.0.0.1:8080
hckctl flow sql 127.0.0.1:3306
hckctl flow atomic-red-team 127.0.0.1 T1485
hckctl flow c2 ping
hckctl flow phishing @example.com
```

### Template

Explore all available templates. Pin a git `revision` to ensure reliability in automated pipelines
```bash
hckctl template list

# TODO how to validate local template
```

Please, feel free to contribute to the companion [repository](https://github.com/hckops/megalopolis) and add more templates.

### Config

Edit the default configurations
```bash
# prints path and current configs
hckctl config

# unix path
vim ${HOME}/.config/hck/config.yml

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

## Provider

### Docker

Follow the official [instructions](https://docs.docker.com/engine/install) to install Docker Engine. The fastest way to get started is with the [convenience script](https://get.docker.com)
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
./sudo sh get-docker.sh
```

### Kubernetes

Use [minikube](https://minikube.sigs.k8s.io) or [kind](https://kind.sigs.k8s.io) to setup a local cluster
```bash
provider:
  kube:
    # by default uses "~/.kube/config"
    configPath: ""
    namespace: hckops
```
If you are looking for a simple way to get started with a remote cluster consider [kube-template](https://github.com/hckops/kube-template)
```bash
provider:
  kube:
    configPath: "~/PATH/TO/kube-template/clusters/do-template-kubeconfig.yaml"
```

### Cloud

Access to the platform is in ***preview*** and limited. If you are interested, please leave a comment or a :thumbsup: to this [issue](https://github.com/hckops/hckctl/issues/104) and we'll reach out with more details
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

## Roadmap

* `machine` create and access VMs e.g. DigitalOcean Droplet, AWS EC2, Azure Virtual Machines, QEMU etc.
* `tui` similar to lazydocker and k9s
* `network` support Tor and ProxyChains
* `plugin` add custom cli commands in any language
  - `man` combine tldr and cheat with task commands
  - `htb` and `thm` api to start/stop/list machines and submit flags
  - `prompt` chatgpt prompt style

## Contribute

> TODO example of how to point to a specific pr/revision in a forked repo

<!--

* rename `template` to catalog? or alias?
* cmd aliases e.g. start/up/create
* replace task/htb example with thm
* verify/support kube config relative path
* update cloud

TODO
* priority
    - add task providers: kube and cloud
    - add box kube --network-vpn
    - debug `htb-postman`
    - lab inputs
    - add flow example
    - verify kube/cloud distroless support
    - verify kube/cloud no-shell support
    - play htb: linux/win
    - RELEASE
* general
    - add disclaimer of responsibility to readme?
    - public discord server (review channels visibility)
    - brew release
    - review context/http/client timeouts e.g. vpn or target not available
    - verify config migration between versions
    - add readme lab video/gif
    - delete old branches (video)
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
    - keep up-to-date directories to exclude in `resolvePath` e.g. charts
    - add filters and review output e.g. table
* box
    - print/event shared directory, same as envs, ports etc.
    - review tty resize
    - expose copy from/to ???
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
    - rename output log file with timestamp?
    - for debug purposes prepend file output with interpolated task (yaml) or command parameters
    - add command to remove all logs
    - skip output file for `help` and `version`
    - add argument `--volume` to restrict shared directories/files
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

-->
