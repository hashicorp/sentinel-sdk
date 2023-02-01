// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"reflect"
)

// various convenience types for reflect calls
var (
	mapTyp            = reflect.TypeOf((*Map)(nil)).Elem()
	interfaceTyp      = reflect.TypeOf((*interface{})(nil)).Elem()
	sliceInterfaceTyp = reflect.TypeOf([]interface{}{})
)

// Reflect takes a value and uses reflection to traverse the value, finding
// any further namespaces that need to be converted to types that can be
// sent across the plugin barrier.
//
// Currently, this means flattening them all to maps. In the future, we intend
// to support "thunks" to allow efficiently transferring this data without
// having to flatten it all.
func (m *Plugin) reflect(value interface{}) (interface{}, error) {
	v, err := m.reflectValue(reflect.ValueOf(value))
	if err != nil {
		return nil, err
	}

	if !v.IsValid() {
		return nil, nil
	}

	return v.Interface(), nil
}

func (m *Plugin) reflectValue(v reflect.Value) (reflect.Value, error) {
	// If the value isn't valid, return right away
	if !v.IsValid() {
		return v, nil
	}

	// Unwrap the interface wrappers
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	// Determine if we have a nil pointer. This will turn a typed nil
	// into a plain nil so that we can turn it into an undefined value
	// properly.
	ptr := v
	for ptr.Kind() == reflect.Ptr {
		ptr = ptr.Elem()
	}
	if !ptr.IsValid() {
		return ptr, nil
	}

	// If the value implements Map, then we call that and use the map
	// value as the actual thing to look at.
	if v.Type().Implements(mapTyp) {
		m, err := v.Interface().(Map).Map()
		if err != nil {
			return v, err
		}

		v = reflect.ValueOf(m)
	}

	switch v.Kind() {
	case reflect.Map:
		return m.reflectMap(v)

	case reflect.Slice:
		return m.reflectSlice(v)

	default:
		return v, nil
	}
}

func (m *Plugin) reflectMap(mv reflect.Value) (reflect.Value, error) {
	// Create a new map for this. This avoids conflicts and panics on shared
	// data, and ensures we aren't altering data in the original namespace.
	// map[string]interface{} is always used, regardless of the actual type of
	// the map sent along. This prevents type panics.
	//
	// Do a quick check to see if we have a zero-value map (nil map). If we do,
	// return that.
	if mv.IsZero() {
		return mv, nil
	}

	// Otherwise make a map and proceed with copy.
	//
	// Preserve key type from the original map.
	result := reflect.MakeMapWithSize(reflect.MapOf(mv.Type().Key(), interfaceTyp), mv.Len())
	for _, k := range mv.MapKeys() {
		v, err := m.reflectValue(mv.MapIndex(k))
		if err != nil {
			return mv, err
		}

		// If the value isn't valid, we set the value of the map to
		// the zero value for the proper type.
		if !v.IsValid() {
			v = reflect.Zero(mv.Type().Elem())
		}

		result.SetMapIndex(k, v)
	}

	return result, nil
}

func (m *Plugin) reflectSlice(v reflect.Value) (reflect.Value, error) {
	// Create a new slice for this. This avoids conflicts and panics on
	// shared data, and ensures that we aren't altering data in the
	// original namespace. []interface{} is always used, regardless of
	// the actual type of the value sent along. This prevents type
	// panics.
	//
	// Do a quick check to see if we have a zero-value map (nil map). If we do,
	// return that.
	if v.IsZero() {
		return v, nil
	}

	// Otherwise make a slice and proceed with copy.
	result := reflect.MakeSlice(sliceInterfaceTyp, v.Len(), v.Cap())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		newElem, err := m.reflectValue(elem)
		if err != nil {
			return v, err
		}

		// If the value isn't valid, we set the value of the element to
		// the zero value for the proper type.
		if !newElem.IsValid() {
			newElem = reflect.Zero(v.Type().Elem())
		}

		result.Index(i).Set(newElem)
	}

	return result, nil
}
