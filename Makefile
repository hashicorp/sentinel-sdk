GOTOOLS = \
	github.com/golang/protobuf/protoc-gen-go@latest \
	github.com/vektra/mockery/v2@latest \
	gotest.tools/gotestsum@latest

SENTINEL_VERSION = 0.18.0
SENTINEL_BIN_PATH := $(shell go env GOPATH)/bin

test: tools
	gotestsum --format=short-verbose $(TEST) $(TESTARGS)

generate: tools
	go generate ./...

modules:
	go mod download && go mod verify

test-ci:
	mkdir -p test-results/sentinel-sdk
	gotestsum --format=short-verbose --junitfile test-results/sentinel-sdk/results.xml

tools:
	@echo $(GOTOOLS) | xargs -t -n1 go install
	go mod tidy

.PHONY: test generate modules test-ci tools
