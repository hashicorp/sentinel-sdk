# Sentinel Import SDK

[![CircleCI](https://circleci.com/gh/hashicorp/sentinel-sdk/tree/master.svg?style=svg)](https://circleci.com/gh/hashicorp/sentinel-sdk/tree/master)
[![GoDoc](https://godoc.org/github.com/hashicorp/sentinel-sdk?status.svg)](https://godoc.org/github.com/hashicorp/sentinel-sdk)

This repository contains the [Sentinel](https://www.hashicorp.com/sentinel)
import SDK. This SDK allows developers to extend Sentinel to source external
information for use in their policies.

Sentinel imports can be written in any language, but the recommended language is
[Go](https://golang.org/). We provide a high-level framework to make writing
imports in Go extremely easy. For other languages, imports can be written by
implementing the
[protocol](https://github.com/hashicorp/sentinel-sdk/blob/master/proto/import.proto)
over gRPC.

To get started writing a Sentinel import, we recommend reading the [extending
Sentinel](https://docs.hashicorp.com/sentinel/extending/dev) guide.

You can also view the import API via
[GoDoc](https://godoc.org/github.com/hashicorp/sentinel-sdk).

## SDK Compatibility Matrix

Sentinel's plugin protocol is, at this time, _not_ backwards compatible.  This
means that a specific version of the Sentinel runtime is always coupled to a
specific version of the plugin protocol, and SDK. The following table can help
you determine which version of the SDK is necessary to work with which versions
of Sentinel.

Sentinel Version|Plugin Protocol Version|SDK Version
-|-|-
**Current (Up to v0.10.4)**|**1**|**Up to v0.1.1**
Planned for v0.11.0|2|Since v0.2.0

## Development Info

The following tools are required to work with the Sentinel SDK:

* [The Sentinel runtime](https://docs.hashicorp.com/sentinel/downloads), usually
  at the most recent version. There are rare exceptions to this, such as when
  the protocol is in active development. Refer to the [SDK Compatibility
  Matrix](#sdk-compatibility-matrix) to locate the correct version of the SDK to
  work with the most current version of the runtime.
* [Google's Protocol
  Buffers](https://developers.google.com/protocol-buffers/docs/downloads).

After both of these are installed, you can use the following `make` commands:

* `make test` will run tests on the SDK. You can use the `TEST` and `TESTARGS`
  variables to control the packages and test arguments, respectively.
* `make tools` will install any necessary Go tools.
* `make generate` will generate any auto-generated code. Currently this includes
  the protocol, mockery files, and the code for the plugin testing toolkit.

The `modules`, `test-circle`, and `/usr/bin/sentinel` targets are only used in
Circle and are not necessary for interactive development.

## Help and Discussion

For issues specific to the SDK, please use the GitHub issue tracker (the
[Issues](https://github.com/hashicorp/sentinel-sdk/issues) tab).

For general Sentinel support and discussion, please use the [Sentinel Community
Forum](https://discuss.hashicorp.com/c/sentinel).
