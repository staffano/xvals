package main

import (
	"fmt"
	"os"

	"github.com/staffano/xvals"
	"gopkg.in/yaml.v3"
)

func main() {
	os.Setenv("EP_EP1_ADDRESS", "localhost:123")
	os.Setenv("EP_EP1_TLS", "mTLS")
	os.Setenv("EP_EP1_SERVER_CACERT", "ServerCACert")
	os.Setenv("EP_EP1_SERVER_CERT", "ServerCert")
	os.Setenv("EP_EP1_SERVER_KEY", "ServerKey")
	os.Setenv("EP_EP1_CLIENT_CACERT", "ClientCACert")
	os.Setenv("EP_EP1_CLIENT_CERT", "ClientCert")
	os.Setenv("EP_EP1_CLIENT_KEY", "ClientKey")
	os.Setenv("EP_EP1_PATH", "/this/is/a/path")

	xvals.WithEnvironment()
	xvals.WithObject(xvals.EndpointDescr)
	xvals.ReloadObjects()

	ep1, _ := xvals.GetObject(xvals.TpEndpoint, "ep1")

	d, _ := yaml.Marshal(ep1)
	fmt.Println(string(d))
}
