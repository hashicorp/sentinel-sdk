package framework

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/sentinel-sdk"
	"github.com/hashicorp/sentinel-sdk/encoding"
)

var stringTyp = reflect.TypeOf("")

// Import implements sdk.Import. Configure and return this structure
// to simplify implementation of sdk.Import.
type Import struct {
	// Root is the implementation of the import that the user of the
	// framework should implement. It represents the minimum necessary
	// implementation for an import. See the docs for Root for more details.
	Root Root

	// namespaceMap keeps track of all the Namespaces for the various
	// executions. These are cleaned up based on the ExecDeadline.
	namespaceMap  map[uint64]Namespace
	namespaceLock sync.RWMutex
}

// plugin.Import impl.
func (m *Import) Configure(raw map[string]interface{}) error {
	// Verify the root implementation is a Namespace or NamespaceCreator.
	switch m.Root.(type) {
	case Namespace:
	case NamespaceCreator:
	default:
		return fmt.Errorf("invalid import implementation, please report a " +
			"bug to the developer of this import")
	}

	// Configure the object itself
	return m.Root.Configure(raw)
}

// plugin.Import impl.
func (m *Import) Get(reqs []*sdk.GetReq) ([]*sdk.GetResult, error) {
	resp := make([]*sdk.GetResult, len(reqs))
	for i, req := range reqs {
		// Get the namespace
		ns := m.namespace(req)

		// If Context is supplied and the root supports New, handle it.
		// We use the value of constructorOk later on to determine if the
		// return value will be callable as well.
		constructor, constructorOk := ns.(New)
		if req.Context != nil {
			if constructorOk {
				var err error
				ns, err = constructor.New(req.Context)
				if err != nil {
					return nil, fmt.Errorf("error instantiating namespace: %s", err)
				}

				if ns == nil {
					// No namespace was returned. This is technically
					// undefined, but we need to check to make sure there were
					// no function calls first, or else this is an error, not
					// undefined.
					for i, k := range req.Keys {
						if k.Call() {
							return nil, fmt.Errorf(
								"attempting to call function %q on undefined receiver",
								strings.Join(req.GetKeys()[:i+1], "."))
						}
					}

					// If this was just a get call, we can short-circuit the
					// result here, with undefined.
					resp[i] = &sdk.GetResult{
						KeyId: req.KeyId,
						Keys:  req.GetKeys(),
						Value: sdk.Undefined,
					}
					continue
				}
			} else {
				// Invalid implementation. This should not happen and is
				// indicative of something more than likely wrong with the
				// runtime. Nonetheless, this is not the import's problem
				// as the malformed data did not come from it.
				return nil, errors.New(
					"sdk.GetReq.Context present but import does not support framework.New")
			}
		}

		// For each key, perform a get
		var result interface{} = ns
		for i, k := range req.Keys {
			// If we have arguments at this level, perform a function call.
			if k.Call() {
				x, ok := result.(Call)
				if !ok {
					return nil, fmt.Errorf(
						"key %q doesn't support function calls",
						strings.Join(req.GetKeys()[:i+1], "."))
				}

				v, err := m.call(x.Func(k.Key), k.Args)
				if err != nil {
					return nil, fmt.Errorf(
						"error calling function %q: %s",
						strings.Join(req.GetKeys()[:i+1], "."), err)
				}

				result = v
				continue
			}

			switch x := result.(type) {
			// For namespaces, we get the next value in the chain
			case Namespace:
				v, err := x.Get(k.Key)
				if err != nil {
					return nil, fmt.Errorf(
						"error retrieving key %q: %s",
						strings.Join(req.GetKeys()[:i+1], "."), err)
				}

				result = v

			// For maps with string keys, get the value. If the value is
			// nil, return sdk.Null to ensure that we don't mess with how
			// reflection deals with "invalid" zero values in maps. See
			// Import.reflectMap for more details.
			case map[string]interface{}:
				var ok bool
				if result, ok = x[k.Key]; ok {
					if result == nil {
						result = sdk.Null
					}
				} else {
					result = sdk.Undefined
				}

			// Else...
			default:
				// If it is a map with reflection with a string key,
				// then access it.
				v := reflect.ValueOf(x)
				if v.Kind() == reflect.Map && v.Type().Key() == stringTyp {
					// If the value exists within the map, set it to the value
					if v = v.MapIndex(reflect.ValueOf(k.Key)); v.IsValid() {
						result = v.Interface()
						break
					}
				}

				// Finally, its undefined
				result = nil
			}

			if result == nil {
				break
			}
		}

		var err error
		result, err = m.resultReflect(result)
		if err != nil {
			return nil, fmt.Errorf(
				"error retrieving key %q: %s",
				strings.Join(req.GetKeys(), "."), err)
		}

		// Convert the result based on types
		if result == nil {
			result = sdk.Undefined
		}

		// Start building the actual result
		resp[i] = &sdk.GetResult{
			KeyId: req.KeyId,
			Keys:  req.GetKeys(),
			Value: result,
		}

		// If the root supported framework.New and we have a map, flag
		// the ability to call methods on the result.
		if _, ok := result.(map[string]interface{}); ok && constructorOk {
			resp[i].Callable = true
		}

		// If Context was supplied, get the receiver to be returned
		if req.Context != nil {
			respCtxRaw, err := m.resultReflect(ns)
			if err != nil {
				return nil, fmt.Errorf(
					"error marshaling receiver after retrieving key %q: %s",
					strings.Join(req.GetKeys(), "."), err)
			}

			respCtx, ok := respCtxRaw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf(
					"error marshaling receiver after retrieving key %q: receiver is no longer an object",
					strings.Join(req.GetKeys(), "."))
			}

			if respCtx == nil {
				return nil, fmt.Errorf(
					"error marshaling receiver after retrieving key %q: receiver is now nil",
					strings.Join(req.GetKeys(), "."))
			}

			resp[i].Context = respCtx
		}

		// End of request loop. Result should be fully built here and
		// request will proceed to the next request (if it exists in the
		// payload).
	}

	// Done processing all requests, return response.
	return resp, nil
}

func (m *Import) resultReflect(result interface{}) (interface{}, error) {
	// If we have a Map implementation, we return the whole thing.
	if m, ok := result.(Map); ok {
		var err error
		result, err = m.Map()
		if err != nil {
			return nil, err
		}
	}

	// We now need to do a bit of reflection to convert any dangling
	// namespace values into values that can be returned across the
	// plugin interface.
	result, err := m.reflect(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// namespace returns the namespace for the request.
func (m *Import) namespace(req *sdk.GetReq) Namespace {
	if global, ok := m.Root.(Namespace); ok {
		return global
	}

	// Look for it in the cache of executions
	m.namespaceLock.RLock()
	ns, ok := m.namespaceMap[req.ExecId]
	m.namespaceLock.RUnlock()
	if ok {
		return ns
	}

	nsFunc, ok := m.Root.(NamespaceCreator)
	if !ok {
		panic("Root must be NamespaceCreator if not Namespace")
	}

	// Not found, we have to create it
	m.namespaceLock.Lock()
	defer m.namespaceLock.Unlock()

	// If it was created while we didn't have the lock, return it
	ns, ok = m.namespaceMap[req.ExecId]
	if ok {
		return ns
	}

	// Init if we have to
	if m.namespaceMap == nil {
		m.namespaceMap = make(map[uint64]Namespace)
	}

	// Create it
	ns = nsFunc.Namespace()
	m.namespaceMap[req.ExecId] = ns

	// Create the expiration function
	time.AfterFunc(time.Until(req.ExecDeadline), func() {
		m.invalidateNamespace(req.ExecId)
	})

	return ns
}

func (m *Import) invalidateNamespace(id uint64) {
	m.namespaceLock.Lock()
	defer m.namespaceLock.Unlock()
	delete(m.namespaceMap, id)
}

// call performs the typed function call via reflection for f.
func (m *Import) call(f interface{}, args []interface{}) (interface{}, error) {
	// If a function call isn't supported for this key, then it is an error
	if f == nil {
		return nil, fmt.Errorf("function call unsupported")
	}

	// Reflect on the function and verify it is a function
	funcVal := reflect.ValueOf(f)
	if funcVal.Kind() != reflect.Func {
		return nil, fmt.Errorf(
			"internal error: import didn't return function for key")
	}
	funcType := funcVal.Type()

	// Verify argument count
	if len(args) != funcType.NumIn() {
		return nil, fmt.Errorf(
			"expected %d arguments, got %d",
			funcType.NumIn(), len(args))
	}

	// Go through the arguments and convert them to the proper type
	funcArgs := make([]reflect.Value, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		arg := args[i]
		argValue := reflect.ValueOf(arg)

		// If the raw argument cannot be assign to the expected arg
		// types then we attempt a conversion. This is slow because we
		// expect this to be rare.
		t := funcType.In(i)
		if !argValue.Type().AssignableTo(t) {
			v, err := encoding.GoToValue(arg)
			if err != nil {
				return nil, fmt.Errorf(
					"error converting argument to %s: %s",
					t, err)
			}

			arg, err = encoding.ValueToGo(v, t)
			if err != nil {
				return nil, fmt.Errorf(
					"error converting argument to %s: %s",
					t, err)
			}

			argValue = reflect.ValueOf(arg)
		}

		funcArgs[i] = argValue
	}

	// Call the function
	funcRets := funcVal.Call(funcArgs)

	// Build the return values
	var err error
	if len(funcRets) > 1 {
		if v := funcRets[1].Interface(); v != nil {
			err = v.(error)
		}
	}

	return funcRets[0].Interface(), err
}
