# Sentinel Import SDK

This repository contains the [Sentinel](https://www.hashicorp.com/sentinel)
import SDK. This SDK allows developers to extend Sentinel to source external
information for use in their policies.

Sentinel imports can be written in any language, but the recommended
language is [Go](https://golang.org/). We provide a high-level framework
to make writing imports in Go extremely easy. For other languages, imports
can be written by implementing the [protocol](https://github.com/hashicorp/sentinel-sdk/blob/master/proto/import.proto) over gRPC.

To get started writing a Sentinel import, we recommend reading the
[extending Sentinel](https://docs.hashicorp.com/sentinel/extending/dev) guide.

## License

```
Copyright 2019 HashiCorp

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
```
