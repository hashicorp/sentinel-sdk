GOTOOLS = \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/vektra/mockery/cmd/mockery \
	gotest.tools/gotestsum

test:
	go test ./...

test-circle: tools
	mkdir -p test-results/
	gotestsum --junitfile test-results/results.xml

tools:
	go install $(GOTOOLS)

generate: tools
	go generate ./...

.PHONY: test test-circle tools generate
