// Package framework contains a high-level framework for implementing
// Sentinel imports with Go.
//
// The direct sdk.Import interface is a low-level interface that is
// tediuos, clunky, and difficult to implement correctly. The
// interface is this way to assist in the performance of imports
// while executing Sentinel policies. This package provides a
// high-level API that eases import implementation while still
// supporting the performance-sensitive interface underneath.
//
// Imports are generally activated in this framework by serving the
// plugin with the root namespace embedded in Import:
//
//     package main
//
//     import (
//         "github.com/hashicorp/sentinel-sdk"
//         "github.com/hashicorp/sentinel-sdk/rpc"
//     )
//
//     func main() {
//         rpc.Serve(&rpc.ServeOpts{
//             ImportFunc: func() sdk.Import {
//                 return &framework.Import{Root: &root{}}
//             },
//         })
//     }
//
// The plugin framework is based around the concept of namespaces.
// Root is the entrypoint namespace and must be implemented as a
// minimum. From there, nested access may be delegated to other
// Namespace implementations.
//
// Namespaces outside of the root must at least implement the
// Namespace interface. All namespaces, including the root, may
// implement the optional Call or Map interfaces, to support function
// calls or selective memoization calls, respectively.
//
// Root namespaces are generally global, that is, for the lifetime of
// the execution of Sentinel, one single import Root namespace state
// will be shared by all policies that need to be executed. Take care
// when storing state in the Root namespace. If you require state
// in the Root namespace that must be unique across policy
// executions, implement the NamespaceCreator interface.
//
// The Root namespace (or the NamespaceCreator interface, which
// embeds Root) may optionally implement the New interface, which
// allows for the construction of namespaces via the handling of
// arbitrary object data. New is ignored for namespaces past the
// root.
//
// Non-primitive import return data is normally memoized, including
// for namespaces. This prevents expensive calls over the plugin RPC.
// Memoization can be controlled by a couple of methods:
//
// * Implementing the Map interface allows for the explicit return of
// a map of values, sidestepping struct memoization. Normally, this
// is combined with the MapFromKeys function which will call Get for
// each defined key and add the return values to the map. Note that
// multi-key import calls always bypass memoization - so if foo.bar
// is a namespace that implements Map but foo.bar.baz is looked up in
// a single expression, it does not matter if baz is excluded from
// Map.
//
// * Struct memoization is implicit otherwise. Only exported fields
// are acted on - fields are lower and snake cased where applicable.
// To control this behavior, you can use the "sentinel" struct tag.
// sentinel:"NAME" will alter the field to have the name indicated by
// NAME, while an empty string will exclude the field.
//
// Additionally, there are a couple of nuances that the plugin author
// should be cognizant of:
//
// * nil values within slices, maps, and structs are converted to
// nulls in the return object.
//
// * Returning a nil from a Get call is undefined, not null.
//
// The author can alter this behavior explicitly by assigning or
// returning the sdk.Null and sdk.Undefined values.
package framework
