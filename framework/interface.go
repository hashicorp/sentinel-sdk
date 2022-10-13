package framework

//go:generate rm -f mock_*.go
//go:generate mockery --inpackage --note "Generated code. DO NOT MODIFY." --name=Root --testonly
//go:generate mockery --inpackage --note "Generated code. DO NOT MODIFY." --name=Namespace --testonly
// go:generate mockery --inpackage --note "Generated code. DO NOT MODIFY." --name=NamespaceCreator --testonly

// Root is the plugin root. For any plugin, there is only a single root.
// For example, if you're implementing a plugin named "time", then the "time"
// identifier itself represents the plugin root.
//
// The root of an plugin is configurable and is able to return the actual
// interfaces uses for value retrieval. The root itself can never contain
// a value, be callable, return all mappings, etc.
//
// A single root implementation and instance may be shared by many policy
// executions if their configurations match.
type Root interface {
	// Configure is called to configure this plugin with the operator
	// supplied configuration for this plugin.
	Configure(map[string]interface{}) error

	// Root must further implement one of two interfaces: NamespaceCreator
	// or Namespace itself. See the documentation for each for when you'd
	// want to implement one or the other. If neither is implemented,
	// an error will be returned immediately upon configuration.
}

// NamespaceCreator is an interface only used in conjunction with the
// Root interface. It allows the Root implementation to create a unique
// Namespace implementation for each policy execution.
//
// This is useful for plugins that maintain state per policy execution.
// For example for the "time" package, it may be useful for the state to
// be the current time so that all access returns a singular view of time
// for a policy execution.
//
// If your plugin doesn't require per-execution state, Root should
// implement Namespace directly instead.
type NamespaceCreator interface {
	Root

	// Namespace is called to return the root namespace for accessing keys.
	//
	// This will be called once for each policy execution. If data and access
	// is shared by all policy executions (such as static data), then you
	// can return a singleton value.
	//
	// If each policy execution should maintain its own state, then this
	// should return a new value.
	Namespace() Namespace
}

// New is an interface indicating that the namespace supports object
// construction via the handling of arbitrary object data. New is
// only supported on root namespaces, so either created through
// Root or NamespaceCreator.
//
// The format of the object and the kinds of namespaces returned by
// the constructor are up to the plugin author.
type New interface {
	Namespace

	// New is called to construct new namespaces based on arbitrary
	// receiver data.
	//
	// The format of the object and the kinds of namespaces returned by
	// the constructor are up to the plugin author.
	//
	// Namespaces returned by this function must implement
	// framework.Map, or else errors will be returned on
	// post-processing of the receiver.
	//
	// New should return an error if there are issues instantiating the
	// namespace. This includes if the namespace cannot be determined
	// from the receiver data. Returning nil from this function will
	// return undefined to the caller.
	New(map[string]interface{}) (Namespace, error)
}

// Namespace represents a namespace of attributes that can be requested
// by key. For example in "time.pst.hour, time.pst.minute", "time.pst" would
// be a namespace.
//
// Namespaces are either represented or returned by the Root implementation.
// Root is the top-level implementation for a plugin. See Plugin and Root
// for more details.
//
// A Namespace on its own doesn't allow accessing the full mapping of
// keys and values. Map may be optionally implemented to support this.
// Following the example in the first paragraph of this documentation,
// "time.pst" itself wouldn't be allowed for a Namespace on its own. If
// the implementation also implements Map, then "time.pst" would return
// a complete mapping.
type Namespace interface {
	// Get requests the value for a specific key. This must return a value
	// convertable by lang/object.ToObject or another Interface value.
	//
	// If the value doesn't exist, nil should be returned. This will turn
	// into "undefined" eventually in the Sentinel policy. If you want to
	// return an explicit "null" value, please return object.Null directly.
	//
	// If an Interface implementation is returned, this is treated like
	// a namespace. For example, "time.pst" may return an Interface since
	// the value itself expects further keys such as ".hour".
	Get(string) (interface{}, error)
}

// Map is a Namespace that supports returning the entire map of data.
// For example, if "time.pst" implemented this, then the writer of a policy
// may request "time.pst" and get the entire value back as a map.
type Map interface {
	Namespace

	// Map returns the entire map for this value. The return value
	// must only contain values convertable by lang/object.ToObject. It
	// cannot contain functions or other framework interface implementations.
	Map() (map[string]interface{}, error)
}

// Call is a Namespace that supports call expressions. For example, "time.now()"
// would invoke the Func function for "now".
type Call interface {
	Namespace

	// Func returns a function to call for the given string. The function
	// must take some number of arguments and return (interface{}, error).
	// The argument types may be Go types and the framework will handle
	// conversion and validation automatically.
	//
	// The returned function may also return only interface{}. In this case,
	// it is assumed an error scenario is impossible. Any other number of
	// return values will result in an error.
	//
	// This should return nil if the key doesn't support being called.
	Func(string) interface{}
}
