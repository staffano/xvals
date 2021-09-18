package xvals

import "testing"

var objMap = map[string]string{
	"to_object1_name":  "John",
	"to_object1_age":   "44",
	"to_object1_phone": "122",
	"to_object2_name":  "Lisa",
	"to_object2_age":   "52",
	"to_object2_phone": "123",
	"to_object3_name":  "Andrew",
	"to_object3_age":   "12",
	"to_object3_phone": "124",
}
var TestObjectDescriptor = GenericObjectDescriptor("TO", []string{"name", "age", "phone", "email"})

func GetGoodObj(t *testing.T, os *ObjectStore, typ, name, field string, exp string) {
	o, e := os.Get(typ, name)
	if e != nil {
		t.FailNow()
	}
	v, e := o.Get(field)
	if e != nil {
		t.FailNow()
	}
	if v != exp {
		t.FailNow()
	}
}
func TestObjects(t *testing.T) {

	os := NewObjectStore()
	os.AddDescriptor(TestObjectDescriptor)
	os.Reload(objMap)

	GetGoodObj(t, os, "to", "object1", "name", "John")

	_, e := os.Get("apa", "arne")
	if e == nil {
		t.FailNow()
	}

	no, e := os.New("TO", "Kalle")
	if e != nil {
		t.FailNow()
	}

	e = no.Set("phone", "1234")
	if e != nil {
		t.FailNow()
	}
	GetGoodObj(t, os, "to", "kalle", "phone", "1234")
}
