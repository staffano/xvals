package xvals

import (
	"fmt"
	"strconv"
)

// ctxt contains the executable wide view of the xvals.
var ctxt = []XvalProvider{}

// objectStore is the default object store.
var objectStore = NewObjectStore()

// HasValue returns true if the value exist in the context
func HasValue(key string) bool {
	_, e := Value(key)
	return e == nil
}

// Value returns the value as a string. An error is returned
// if the value didn't exist.
func Value(key string) (string, error) {
	for _, v := range ctxt {
		if r, e := v.Value(key); e == nil {
			return r, nil
		}
	}
	return "", fmt.Errorf("key not found %s", key)
}

// BoolValue is a convenience function to fetch and parse
// a value as a boolean
func BoolValue(key string) (val bool, err error) {
	var sv string
	if sv, err = Value(key); err != nil {
		return false, err
	}
	return strconv.ParseBool(sv)
}

// IntValue is a convenience function to fetch and parse
// a value as an int
func IntValue(key string) (val int, err error) {
	var sv string
	if sv, err = Value(key); err != nil {
		return 0, err
	}
	return strconv.Atoi(sv)
}

// Dump returns a merged set of all values available.
func Dump() map[string]string {
	res := make(map[string]string)
	// loop backwards, so that the values of the more prioritized
	// providers are used.
	for i := len(ctxt) - 1; i >= 0; i-- {
		for k, v := range ctxt[i].Dump() {
			res[k] = v
		}
	}
	return res
}

// Store operations

// Objects returns the objects known to the store.
func Objects() map[string]Object {
	return objectStore.Objects()
}

// Get an object based on type and name
func GetObject(typ, name string) (Object, error) {
	return objectStore.Get(typ, name)
}

// NewObject creates a new object with the name and type and adds it to
// default store.
func NewObject(typ, name string) (Object, error) {
	return objectStore.New(typ, name)
}

// ReloadObjects reloads objects based on the current external values
func ReloadObjects() {
	objectStore.Reload(Dump())
}

// WithObject adds support for a specific object type.
func WithObject(descr Descriptor) {
	objectStore.AddDescriptor(descr)
}
