package xvals

import (
	"testing"
)

func TestGetAllEndpoints(t *testing.T) {

	ep1 := &Endpoint{
		Address:      "Address",
		TLS:          "TLS",
		ServerCACert: "ServerCACert",
		ServerCert:   "ServerCert",
		ServerKey:    "ServerKey",
		ClientCACert: "ClientCACert",
		ClientCert:   "ClientCert",
		ClientKey:    "ClientKey",
	}
	epVars := map[string]string{
		"ep_ep1_Address":       "Address",
		"ep_ep1_TLS":           "TLS",
		"ep_ep1_Server_CACert": "ServerCACert",
		"ep_ep1_Server_Cert":   "ServerCert",
		"ep_ep1_Server_Key":    "ServerKey",
		"ep_ep1_Client_CACert": "ClientCACert",
		"ep_ep1_Client_Cert":   "ClientCert",
		"ep_ep1_Client_Key":    "ClientKey",
	}
	os := NewObjectStore()
	os.AddDescriptor(EndpointDescr)
	os.Reload(epVars)
	ep1r, e := os.Get("ep", "ep1")
	if e != nil {
		t.FailNow()
	}
	for v := range ep1.Fields() {
		f1, e := ep1.Get(v)
		if e != nil {
			t.FailNow()
		}
		f2, e := ep1r.Get(v)
		if e != nil {
			t.FailNow()
		}
		if f1 != f2 {
			t.FailNow()
		}
	}
}
