package xvals

import (
	"fmt"
	"strings"
)

// Object is the interface all supported objects must implement. It is described by
// a Descriptor, which is used to generate the object.
//
// A field is defined by a key in the format:
//
//		<TYPE>_<NAME>_<FIELD>=<VALUE>
//      <----    KEY    ---->
//
// An Object is made up of the sum of it's keys.
type Object interface {
	Type() string
	Fields() map[string]string
	Get(field string) (value string, err error)
	Set(field, val string) error
}

// The Descriptor interface is used to marshal/unmarshal objects via key/values
type Descriptor interface {
	Construct() Object
	Type() string     // "EP"
	Fields() []string // "ADDRESS", "TLS"
}

// An ObjectStore is a storage where objects described according  can be
// stored. It uses xvals and descriptors to extract the keys/values that are used to
// build the objects.
type ObjectStore struct {
	descriptors map[string]Descriptor
	objects     map[string]Object
}

// NewObjectStore creates a new store for objects.
func NewObjectStore() *ObjectStore {
	s := &ObjectStore{
		descriptors: make(map[string]Descriptor),
		objects:     make(map[string]Object),
	}
	return s
}

// tu is strings.ToUpper
func tu(val string) string {
	return strings.ToUpper(val)
}

// AddDescriptor lets the store use a new descriptor for objects.
func (c *ObjectStore) AddDescriptor(descriptor Descriptor) {
	c.descriptors[tu(descriptor.Type())] = descriptor
}

// Reload the store from the set of key/values
func (c *ObjectStore) Reload(kv map[string]string) {
	for k, v := range kv {
		typ, name, field := c.extractTypeNameField(tu(k))
		if typ == "" {
			// Not a field of a recognizable object
			continue
		}
		var (
			obj Object
			ok  bool
		)
		if obj, ok = c.objects[Key(typ, name)]; !ok {
			obj = c.descriptors[typ].Construct()
		}
		c.objects[Key(typ, name)] = obj
		obj.Set(field, v)
	}
}

// Objects returns the objects known to the store.
func (c *ObjectStore) Objects() map[string]Object {
	return c.objects
}

// Get an object based on type and name
func (c *ObjectStore) Get(typ, name string) (Object, error) {
	key := Key(tu(typ), tu(name))
	obj, ok := c.objects[key]
	if !ok {
		return nil, fmt.Errorf("object with key %s not found", key)
	}
	return obj, nil
}

func (c *ObjectStore) New(typ, name string) (Object, error) {
	d, ok := c.descriptors[tu(typ)]
	if !ok {
		return nil, fmt.Errorf("don't know how to create an object from type %s", typ)
	}
	obj := d.Construct()
	c.objects[Key(tu(typ), tu(name))] = obj
	return obj, nil
}

// GenericObject provides a default implementation for objects.
// It can be used for all objects that has no need for a formal
// go struct as base for the type.
type GenericObject struct {
	typ    string
	fields map[string]string
}

func (o *GenericObject) Type() string              { return o.typ }
func (o *GenericObject) Fields() map[string]string { return o.fields }
func (o *GenericObject) Get(field string) (value string, err error) {
	v, ok := o.fields[tu(field)]
	if !ok {
		return "", fmt.Errorf("no field %s in object", field)
	}
	return v, nil
}
func (o *GenericObject) Set(key, val string) error {
	o.fields[tu(key)] = val
	return nil
}

type genericDescriptor struct {
	typ    string
	fields []string
}

// GenericObjectDescriptor generates a ObjectDescriptor for a generic type object.
func GenericObjectDescriptor(
	Type string,
	Fields []string) Descriptor {
	d := &genericDescriptor{
		typ: Type, fields: Fields,
	}
	return d
}

func (d *genericDescriptor) Construct() Object {
	return &GenericObject{
		typ:    d.Type(),
		fields: make(map[string]string),
	}
}

func (d *genericDescriptor) Type() string {
	return tu(d.typ)
}
func (d *genericDescriptor) Fields() []string {
	return d.fields
}

// Utility functions for parsing keys as type/name/field

// Key returns the key based on type+name
func Key(typ, name string) string {
	return typ + "+" + name
}
func FromKey(key string) (typ, name string) {
	s := strings.Split(key, "+")
	if len(s) == 1 {
		return s[0], ""
	}
	if len(s) == 2 {
		return s[0], s[1]
	}
	return "", ""
}

func extractNameField(key string, typ string, fieldNames []string) (name, field string) {
	a := strings.TrimPrefix(key, typ+"_")
	if key == a {
		return "", ""
	}
	for _, field = range fieldNames {
		fu := tu(field)
		name = strings.TrimSuffix(a, "_"+fu)
		if name != a {
			return name, fu
		}
	}
	return "", ""
}
func (c *ObjectStore) extractTypeNameField(key string) (typ, name, field string) {
	for _, d := range c.descriptors {
		typ = d.Type()
		name, field = extractNameField(key, typ, d.Fields())
		if name != "" {
			return
		}
	}
	return "", "", ""
}
