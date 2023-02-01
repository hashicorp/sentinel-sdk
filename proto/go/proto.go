// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package proto contains the Go generated files for the protocol buffer files.
package proto

//go:generate protoc -I ../ ../plugin.proto --go_out=plugins=grpc:.
