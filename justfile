BUILD_PATH := "./build"
BIN_NAME := "hckctl"

default: build

format:
	go fmt ./...

test:
	go test ./...

build: format test
	rm -frv {{BUILD_PATH}}
	go build -o {{BUILD_PATH}}/{{BIN_NAME}} main.go
