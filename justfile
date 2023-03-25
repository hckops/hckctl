BUILD_PATH := "./build"
BIN_NAME := "hckctl"

GO_BUILD_ENV := "CGO_ENABLED=0"
GO_FILES := "./..."

default: build

format:
  go fmt {{GO_FILES}}

vet:
  go vet {{GO_FILES}}

test:
  go test {{GO_FILES}} -cover

build: format vet test
  rm -frv {{BUILD_PATH}}
  {{GO_BUILD_ENV}} go build -o {{BUILD_PATH}}/{{BIN_NAME}} main.go
