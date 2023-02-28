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

<!--
The Cloud Native HaCKing Tool
-->

> WIP

## Development

* [just](https://github.com/casey/just)

```bash
# build
just

# run
./build/hckctl

# debug
./build/hckctl <CMD> --log-level debug

# logs
tail -f /tmp/hckctl-*.log
```

Box command
```bash
# lists boxes
go run main.go box list

# starts a docker box (default)
go run main.go box alpine
go run main.go box alpine --provider docker

# starts a kubernetes box
go run main.go box alpine --provider kube

# TODO starts a remote box
go run main.go box parrot --provider cloud
```

Template command
```bash
# validates and prints remote template
go run main.go template parrot | yq -o=json

# validates and prints local template
go run main.go template -p ../megalopolis/boxes/official/alpine.yml
```

Config command
```bash
# edits config file
vim ~/.config/hck/config.yml

# prints current config
go run main.go config
```
