<h1 align="center"><code>hckctl</code></h1>

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
```

Box command
```bash
# list boxes
go run main.go box list

# start a docker box
go run main.go box alpine --docker

# start a kubernetes box
go run main.go box alpine --kube

# TODO start a remote box
go run main.go box parrot
```

Template command
```bash
# validate and print remote template
go run main.go template parrot | yq -o=json

# validate and print local template
go run main.go template -p ../megalopolis/boxes/official/alpine.yml
```
