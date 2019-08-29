GOTOOLS = \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/vektra/mockery/cmd/mockery \
	gotest.tools/gotestsum

SENTINEL_VERSION = 0.10.4

test:
	go test ./...

generate: tools
	go generate ./...

modules:
	go mod download

# FIXME: Remove the "testing" filter after Sentinel 0.11.0. These
# tests are currently broken due to the discrepancy in protocol
# version.
test-circle:
	mkdir -p test-results/sentinel-sdk
	gotestsum --format=short-verbose --junitfile test-results/sentinel-sdk/results.xml $(shell go list ./... | grep -v testing)

tools:
	go install $(GOTOOLS)

/usr/bin/sentinel:
	gpg --import .circleci/hashicorp.gpg && \
	cd /tmp && \
	curl -O https://releases.hashicorp.com/sentinel/${SENTINEL_VERSION}/sentinel_${SENTINEL_VERSION}_linux_amd64.zip && \
	curl -O https://releases.hashicorp.com/sentinel/${SENTINEL_VERSION}/sentinel_${SENTINEL_VERSION}_SHA256SUMS && \
	curl -O https://releases.hashicorp.com/sentinel/${SENTINEL_VERSION}/sentinel_${SENTINEL_VERSION}_SHA256SUMS.sig && \
	gpg --verify sentinel_${SENTINEL_VERSION}_SHA256SUMS.sig && \
	shasum --check --ignore-missing sentinel_${SENTINEL_VERSION}_SHA256SUMS && \
	cd /usr/bin && \
	sudo unzip /tmp/sentinel_${SENTINEL_VERSION}_linux_amd64.zip

.PHONY: test generate modules test-circle tools
