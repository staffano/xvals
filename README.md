# xvals

Introduction 
------------
xvals is a library for accessing external values. The library can be configured to read
values from environment variables, config file and command line paramters.

It also contains an Object model, so that values grouped using a specific scheme
can be treated as an object. 

This initial use case for this was to unify the
notion of Endpoints when dealing with multiple microservices. It facilitates the
handling of these endpoints in regards of both dev ops and implementation.

Installation and usage
----------------------

To install it, run:

    go get github.com/staffano/xvals

License
-------

The xvals package is licensed under the Apache License 2.0. Please see the [LICENSE](LICENSE.txt) file for details.


Example
-------

```Go
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

```
