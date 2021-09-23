package xvals

import (
	"fmt"
	"io"
	"strings"
)

/*
	* The client need to know the following
		- Address
		- Port
		- Protocol
		- Path
		- If tls is used
		- CA Certificate used to sign servers private key

	* The server needs the following
		- Address
		- Port
		- Protocol
		- Path
		- If tls is used
		- Private Key
		- Certificate (Address + public key)

	* Outside, there has to be a CA that
		- Generate private keys
		- Generate a certificate from a private key + address
		- Provides the CA certificate for clients.

	The protocol is selected by the user by selecting protocol
	specific operations.
	 - GrpcClientConn("arne")
	 - HttpClientConn("arne")
	 - GrpcServerListener("arne")
	 - HttpServerListener("arne")

	The path is handled by the user

	For a client user we store the following as values
	- TLS=[none|server|mTls]
	- Address=host:port
	- ServerCACert=[file|string|"external"]
	- ClientCertificate=[file|string]
	- ClientPrivateKey=[file|string]

	For a server user we store the following as values
	- TLS=[true|false]
	- Address=host:port
	- ServerCertificate=[file|string]
	- ServerPrivateKey=[file|string]
	- ClientCACert=[file|string]

	A Server can decide to listen to all interfaces, but a client has
	to connect to a specific address. This is addressed in the protocol function
	- GrpcServerListener("arne",WithAllInterfaces())

	Using this formula, we can reuse the value for both client and the server.
	The client only uses the ServerCACert part and the server only uses
	ServerCertificate/ServerPrivateKey Same goes for mTLS But we can leave out the
	parts that is not needed, so the client can operate without the
	Certificate/PrivateKey beeing available.

	External CACert means that the CA is somehow already available to the executing
	program.

*/

// TpEndpoint is the type name for an endpoint
const TpEndpoint = "EP"

// Endpoint is either side of a client-server connection
type Endpoint struct {
	Address      string `yaml:"address"`
	TLS          string `yaml:"tls,omitempty"` // Empty is no TLS
	ServerCACert string `yaml:"server_ca_cert,omitempty"`
	ServerCert   string `yaml:"server_cert,omitempty"`
	ServerKey    string `yaml:"server_key,omitempty"`
	ClientCACert string `yaml:"client_ca_cert,omitempty"`
	ClientCert   string `yaml:"client_cert,omitempty"`
	ClientKey    string `yaml:"client_key,omitempty"`
	Path         string `yaml:"path,omitempty"`
}

func (e *Endpoint) Write(name string, w io.Writer) {
	un := strings.ToUpper(name)

	fmt.Fprintf(w, "EP_%s_ADDRESS=%s", un, e.Address)
	fmt.Fprintf(w, "EP_%s_TLS=%s", un, e.TLS)
	fmt.Fprintf(w, "EP_%s_SERVER_CACERT=%s", un, e.ServerCACert)
	fmt.Fprintf(w, "EP_%s_SERVER_CERT=%s", un, e.ServerCert)
	fmt.Fprintf(w, "EP_%s_SERVER_KEY=%s", un, e.ServerKey)
	fmt.Fprintf(w, "EP_%s_CLIENT_CACERT=%s", un, e.ClientCACert)
	fmt.Fprintf(w, "EP_%s_CLIENT_CERT=%s", un, e.ClientCert)
	fmt.Fprintf(w, "EP_%s_CLIENT_KEY=%s", un, e.ClientKey)
	fmt.Fprintf(w, "EP_%s_PATH=%s", un, e.Path)
}

// UseTLS returns false if TLS should not be considered
func (e *Endpoint) UseTLS() bool {
	return e.TLS == "server" || e.TLS == "mTLS"
}

// Get field of the endpoing
func (e *Endpoint) Get(field string) (val string, err error) {

	switch tu(field) {
	case "ADDRESS":
		return e.Address, nil
	case "TLS":
		return e.TLS, nil
	case "SERVER_CACERT":
		return e.ServerCACert, nil
	case "SERVER_CERT":
		return e.ServerCert, nil
	case "SERVER_KEY":
		return e.ServerKey, nil
	case "CLIENT_CACERT":
		return e.ClientCACert, nil
	case "CLIENT_CERT":
		return e.ClientCert, nil
	case "CLIENT_KEY":
		return e.ClientKey, nil
	case "PATH":
		return e.Path, nil
	default:
		return "", fmt.Errorf("field %s is not valid for an endpoint", field)
	}
}

// Set field of the endpoint
func (e *Endpoint) Set(field, val string) error {
	switch tu(field) {
	case "ADDRESS":
		e.Address = val
	case "TLS":
		e.TLS = val
	case "SERVER_CACERT":
		e.ServerCACert = val
	case "SERVER_CERT":
		e.ServerCert = val
	case "SERVER_KEY":
		e.ServerKey = val
	case "CLIENT_CACERT":
		e.ClientCACert = val
	case "CLIENT_CERT":
		e.ClientCert = val
	case "CLIENT_KEY":
		e.ClientKey = val
	case "PATH":
		e.Path = val
	}
	return fmt.Errorf("field %s is not valid for an endpoint", field)
}

// Type returns the type name of the endpoint
func (e *Endpoint) Type() string { return TpEndpoint }

// Fields returns the full set of fields with values
func (e *Endpoint) Fields() map[string]string {
	r := map[string]string{
		"ADDRESS":       e.Address,
		"TLS":           e.TLS,
		"SERVER_CACERT": e.ServerCACert,
		"SERVER_CERT":   e.ServerCert,
		"SERVER_KEY":    e.ServerKey,
		"CLIENT_CACERT": e.ClientCACert,
		"CLIENT_CERT":   e.ClientCert,
		"CLIENT_KEY":    e.ClientKey,
		"PATH":          e.Path,
	}
	return r
}

// EndpointDescr contains the description of the endpoint
var EndpointDescr = &epDescriptor{}

// EpDescriptor is a type descriptor for the endpoint
type epDescriptor struct{}

func (d *epDescriptor) Type() string { return TpEndpoint }
func (d *epDescriptor) Fields() []string {
	return []string{"ADDRESS", "TLS", "SERVER_CACERT", "SERVER_CERT", "SERVER_KEY", "CLIENT_CACERT", "CLIENT_CERT", "CLIENT_KEY", "PATH"}
}
func (d *epDescriptor) Construct() Object {
	return &Endpoint{}
}
