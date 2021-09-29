package xvals

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
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
	r := WithEnvironment()
	GetGood(t, rk, rv)

	os.Unsetenv(rk)

	GetGood(t, rk, rv)
	r.Reload()
	GetBad(t, rk)
}

func TestWithConfigFile(t *testing.T) {
	GetBad(t, "ep_staffan_address")
	p := WithConfigFile("testdata/testctx1.yaml")
	GetGood(t, "ep_staffan_address", "localhost:12345")
	p.Reload()
	GetGood(t, "ep_staffan_address", "localhost:12345")
}

func TestWithProfile(t *testing.T) {
	GetBad(t, "key1")
	p := WithProfile("testdata/testprofiles.yaml")
	GetGood(t, "key1", "val1")
	p.Reload()
	GetGood(t, "key2", "val2")

	// get a value from a profile other than the current
	GetBad(t, "key3")
}

func TestCreateProfileFile(t *testing.T) {
	profiles := &ProfileFile{
		CurrentProfile: "profile1",
		Profiles: map[string]map[string]string{
			"profile1": {"key1": "val1", "key2": "val2"},
			"profile2": {"key3": "val3"}},
	}
	data, err := yaml.Marshal(profiles)
	if err != nil {
		t.FailNow()
	}
	os.WriteFile("testdata/testprofiles.yaml", data, os.ModePerm)
}
