// Package sdk contains the low-level interfaces and API for creating Sentinel
// plugins. A Sentinel plugin can provide data dynamically to Sentinel policies.
//
// For plugin authors, the subfolder "framework" contains a high-level
// framework for easily implementing imports in Go.
package sdk

import (
	"time"
)

type undefined struct{}
type null struct{}

var (
	// Undefined is a constant value that represents the undefined object in
	// Sentinel. By making this the return value, it'll be converted to
	// Undefined.
	Undefined = &undefined{}

	// Null is a constant value that represents the null object in Sentinel.
	// By making this a return value, it will convert to null.
	Null = &null{}
)

//go:generate rm -f mock_Import.go mock_Import_Closer.go
//go:generate mockery -inpkg -note "Generated code. DO NOT MODIFY." -name=Import
//go:generate cp mock_Import_Closer.go.src mock_Import_Closer.go

// Import is an importable package.
//
// Imports are a namespace that may contain objects and functions. The
// root level has no value nor can it be called. For example `import "a'`
// allows access to fields within "a" such as `a.b` but doesn't allow
// referencing or calling it directly with `a` alone. This is an important
// difference between imports and external objects (which themselves express
// a value).
type Import interface {
	// Configure is called to configure the plugin before it is accessed.
	// This must be called before any call to Get().
	Configure(map[string]interface{}) error

	// Get is called when an import field is accessed or called as a function.
	//
	// Get may request more than one value at a time, represented by multiple
	// GetReq values. The result GetResult should contain the matching
	// KeyId for the requests.
	//
	// The result value is not a map keyed on KeyId to allow flexibility
	// in the future of potentially allowing pre-fetched data. This has no
	// effect currently.
	Get(reqs []*GetReq) ([]*GetResult, error)
}

// GetReq are the arguments given to Get for an Import.
type GetReq struct {
	// ExecId is a unique ID representing the particular execution for this
	// request. This can be used to maintain state locally.
	//
	// ExecDeadline is a hint of when the execution will be over and the
	// state can be thrown away. The time given here will always be in UTC
	// time. Note that this is susceptible to clock shifts, but Go is planning
	// to make the time APIs monotonic by default (see proposal 12914). After
	// that this will be resolved.
	ExecId       uint64
	ExecDeadline time.Time

	// Keys is the list of keys being requested. For example for "a.b.c"
	// where "a" is the import, Keys would be ["b", "c"].
	Keys []GetKey

	// KeyId is a unique ID for this key. This should match exactly the
	// GetResult KeyId so that the result for this can be found quickly.
	KeyId uint64

	// Context, if supplied, is an arbitrary object intended to
	// represent the data from an existing namespace. If the import
	// supports the framework.New interface, the contents are passed to
	// it, with any resulting namespace being what Get operates on.
	//
	// The Get call operates on the root of the import if this is set
	// to nil. If this is set and the import does not implement
	// framework.New, an error is returned.
	Context map[string]interface{}
}

// GetKey is an individual key in the larger possible selector of the
// specific import call, along with any supplied arguments for the
// specific key.
type GetKey struct {
	// The key for this part of the request.
	Key string

	// The list of arguments for a call expression. This is "nil" if
	// this key is not a call. This may be length zero (but non-nil) if
	// this is a call with no arguments.
	Args []interface{}
}

// Call returns true if this request is a call expression.
func (g *GetKey) Call() bool {
	return g.Args != nil
}

// GetKeys returns a list of the string keys in the GetReq, without
// the arguments.
func (g *GetReq) GetKeys() []string {
	s := make([]string, len(g.Keys))
	for i, k := range g.Keys {
		s[i] = k.Key
	}

	return s
}

// GetResult is the result structure for a Get request.
type GetResult struct {
	KeyId    uint64                 // KeyId matching GetReq.KeyId, or zero.
	Keys     []string               // Keys structure from GetReq.Keys, or new key set.
	Value    interface{}            // Value compatible with lang/object.ToObject
	Context  map[string]interface{} // Updated Context if it was sent
	Callable bool                   // true if returned Value is callable
}

// GetResultList is a wrapper around a slice of GetResult structures
// to provide helpers.
type GetResultList []*GetResult

// KeyId gets the result with the given key ID, or nil if its not found.
func (r GetResultList) KeyId(id uint64) *GetResult {
	for _, v := range r {
		if v.KeyId == id {
			return v
		}
	}

	return nil
}
