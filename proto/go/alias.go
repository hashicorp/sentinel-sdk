package proto

// Empty is an alias of the EmptyResp type. Empty was renamed to EmptyResp in
// the protocol buffer file to avoid a name conflict in the global "proto"
// package. This alias may go away based on changing the protocol buffers
// package name and is to allow for backwards compatability.
type Empty = EmptyResp
