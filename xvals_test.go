package xvals

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}
func randString(n int) string {
	c := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	a := make([]byte, n)
	for i := range a {
		a[i] = byte(c[rand.Intn(len(c))])
	}
	return string(a)
}

var Fixture = map[string]string{
	"name":     "John",
	"phone":    "12345",
	"home":     "trUe",
	"not-home": "false",
}

func GetGood(t *testing.T, key, exp string) {
	v, e := Value(key)
	if e != nil {
		t.Logf("key: [%s] got:[%v] expected: [nil]", key, e)
		t.FailNow()
	}
	if v != exp {
		t.Logf("key: [%s] got:[%s] expected: [%s]", key, v, exp)
		t.FailNow()
	}
}
func GetGoodInt(t *testing.T, key string, exp int) {
	v, e := IntValue(key)
	if e != nil {
		t.Logf("key: [%s] got:[%v] expected: [nil]", key, e)
		t.FailNow()
	}
	if v != exp {
		t.Logf("key: [%s] got:[%d] expected: [%d]", key, v, exp)
		t.FailNow()
	}
}
func GetBadInt(t *testing.T, key string) {
	_, e := IntValue(key)
	if e == nil {
		t.Logf("key: [%s] expected error", key)
		t.FailNow()
	}
}
func GetGoodBool(t *testing.T, key string, exp bool) {
	v, e := BoolValue(key)
	if e != nil {
		t.Logf("key: [%s] got:[%v] expected: [nil]", key, e)
		t.FailNow()
	}
	if v != exp {
		t.Logf("key: [%s] got:[%v] expected: [%v]", key, v, exp)
		t.FailNow()
	}
}
func GetBad(t *testing.T, key string) {
	_, e := Value(key)
	if e == nil {
		t.Logf("key: [%s] expected error", key)
		t.FailNow()
	}
}
func GetBadBool(t *testing.T, key string) {
	_, e := BoolValue(key)
	if e == nil {
		t.Logf("key: [%s] expected error", key)
		t.FailNow()
	}
}
func TestWithMap(t *testing.T) {

	WithMap(Fixture)
	GetGood(t, "name", "John")
	GetGood(t, "phone", "12345")
	GetGoodInt(t, "phone", 12345)
	GetBadBool(t, "home")
	GetGoodBool(t, "not-home", false)
	GetBad(t, "lastname")
}

func TestWithEnvVars(t *testing.T) {
	rk := randString(10)
	rv := randString(200)
	defer os.Unsetenv(rk)

	os.Setenv(rk, rv)
	GetBad(t, rk)
	r, e := WithEnvironment()
	if e != nil {
		t.FailNow()
	}
	GetGood(t, rk, rv)

	os.Unsetenv(rk)
	if e != nil {
		t.FailNow()
	}
	GetGood(t, rk, rv)
	e = r.Reload()
	if e != nil {
		t.FailNow()
	}
	GetBad(t, rk)
}

func TestWithConfigFile(t *testing.T) {
	GetBad(t, "ep_staffan_address")
	p, e := WithConfigFile("testdata/testctx1.yaml")
	if e != nil {
		t.FailNow()
	}
	GetGood(t, "ep_staffan_address", "localhost:12345")
	if e := p.Reload(); e != nil {
		t.FailNow()
	}
	GetGood(t, "ep_staffan_address", "localhost:12345")
}
