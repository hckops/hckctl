<h1 align="center"><code>hckctl</code></h1>

<p align="center">
  <img width="160" src="docs/logo.svg" alt="logo">
</p>

<!--
The Cloud Native HaCKing Tool
-->

> WIP

## Development

```bash
# list boxes
go run main.go box list

# start local box
go run main.go box alpine --docker

# start remote box
go run main.go box parrot

# validate and print remote template
go run main.go template parrot | yq -o=json

# validate and print local template
go run main.go template -p ../megalopolis/boxes/official/alpine.yml

# build
just

# run
./build/hckctl
```
  