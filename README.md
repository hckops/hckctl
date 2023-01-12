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
# box
go run main.go box parrot
go run main.go box alpine --docker

# template
go run main.go template parrot | yq -o=json
go run main.go template -p ../megalopolis/boxes/official/alpine.yml

# build
just

# run
./build/hckctl
```
  