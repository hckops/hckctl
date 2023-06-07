BUILD_PATH := "./build"
BIN_NAME := "hckctl"

GO_BUILD_ENV := "CGO_ENABLED=0"
GO_FILES := "./..."

default: (build BUILD_PATH)

install:
  go mod tidy
  go mod vendor

format:
  go fmt {{GO_FILES}}

vet:
  go vet {{GO_FILES}}

test:
  go test {{GO_FILES}} -cover

build output $VERSION_COMMIT="$(git rev-parse HEAD)" $VERSION_TIMESTAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)": install format vet test
  rm -frv {{output}}
  {{GO_BUILD_ENV}} go build \
    -ldflags="-X github.com/hckops/hckctl/pkg/command.commit={{VERSION_COMMIT}} -X github.com/hckops/hckctl/pkg/command.timestamp={{VERSION_TIMESTAMP}}" \
    -o {{output}}/{{BIN_NAME}} cmd/main.go
