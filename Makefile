GOTOOLS = \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/vektra/mockery/cmd/mockery

test:
	go test ./...

tools:
	go install $(GOTOOLS)

generate: tools
	go generate ./...

.PHONY: test tools generate
