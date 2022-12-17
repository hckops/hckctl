<h1 align="center"><code>hckctl</code></h1>

<p align="center">
  <img width="160" src="docs/logo.svg" alt="logo">
</p>

The Cloud Native HaCKing Tool

> WIP

## Development

* [just](https://github.com/casey/just)

```bash
# init project (first time)
go mod init github.com/hckops/hckctl

# install|update dependencies
go mod tidy

# build
just

# run
go run main.go
go run main.go box open --local
./build/hckctl box open -l
```
