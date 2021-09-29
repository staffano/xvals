package xvals

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GetEndpoint retrieves and endpoint from the external context
func GetEndpoint(name string) (*Endpoint, error) {
	obj, err := GetObject(TpEndpoint, name)
	if err != nil {
		return nil, err
	}
	ep, ok := obj.(*Endpoint)
	if !ok {
		return nil, fmt.Errorf("could not interpret %s as an Endpoint", name)
	}
	return ep, nil
}

// GetClientTLSConfig returns a *tls.Config suitable for a server
// We need the following values from the endpoint
// - TLS=[none|server|mtls]
// - ServerCACert=[file|string|"external"]
// - ClientCertificate=[file|string]
// - ClientPrivateKey=[file|string]
func (e *Endpoint) GetClientTLSConfig() (*tls.Config, error) {

	var err error
	tlsConfig := new(tls.Config)

	if e.TLS == "mtls" || e.TLS == "server" {

		serverCaCert := fileOrContent(e.ServerCACert)
		if string(serverCaCert) == "" {
			return nil, fmt.Errorf("server TLS requires a valid server CA certificate")
		}
		// Always include the system certificates
		tlsConfig.RootCAs, err = x509.SystemCertPool()
		if err != nil {
			// if we couldn't load the system certs create an empty certpool and
			// continue
			log.Printf("failed to load system certificates, will continues with an empty cert pool")
			tlsConfig.RootCAs = x509.NewCertPool()
		}
		if string(serverCaCert) != "external" {
			if !tlsConfig.RootCAs.AppendCertsFromPEM(serverCaCert) {
				return nil, fmt.Errorf("failed to add Server CA certificate to cert pool")
			}
		}
	}

	// If TLS=mTls then we need the client certificates too
	if e.TLS == "mtls" {
		clientCert := fileOrContent(e.ClientCert)
		clientKey := fileOrContent(e.ClientKey)
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client key pair")
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}

// GetServerTLSConfig returns a *tls.Config suitable for a server
// We need the following values from the endpoint
// - TLS=[none|server|mtls]
// - ServerCertificate=[file|string]
// - ServerPrivateKey=[file|string]
// - ClientCACert=[file|string]
func (e *Endpoint) GetServerTLSConfig() (*tls.Config, error) {

	tlsConfig := new(tls.Config)

	if e.TLS == "server" {
		// Load server key pair
		serverCert := fileOrContent(e.ServerCert)
		serverKey := fileOrContent(e.ServerKey)
		cert, err := tls.X509KeyPair(serverCert, serverKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load server key pair")
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// If TLS=mTls then we need the client certificates too
	if e.TLS == "mtls" {
		clientCaCert := fileOrContent(e.ClientCACert)
		if string(clientCaCert) == "" {
			return nil, fmt.Errorf("client CA cert required for mTLS")
		}
		tlsConfig.ClientCAs = x509.NewCertPool()
		if !tlsConfig.ClientCAs.AppendCertsFromPEM(clientCaCert) {
			return nil, fmt.Errorf("failed to add client CA certificate to cert pool")
		}
	}
	return tlsConfig, nil
}

// GrpcClientConn creates a client grpc connection to the endpoint
func (e *Endpoint) GrpcClientConn(ctx context.Context, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	if e.TLS == "server" || e.TLS == "mTLS" {
		tlsConfig, err := e.GetClientTLSConfig()
		if err != nil {
			return nil, err
		}
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		options = append(options, grpc.WithInsecure())
	}
	return grpc.Dial(e.Address, options...)
}

// GrpcServerListener creates a grpc server and listener on the endpoint
func (e *Endpoint) GrpcServerListener(allInterfaces bool, options ...grpc.ServerOption) (*grpc.Server, net.Listener, error) {
	var listenAddress string
	if allInterfaces {
		_, port, err := net.SplitHostPort(e.Address)
		if err != nil {
			return nil, nil, fmt.Errorf("malformed address in endpoint %s  %w", e.Address, err)
		}
		listenAddress = fmt.Sprintf(":%s", port)
	}

	if e.TLS == "server" || e.TLS == "mTLS" {
		tlsConfig, err := e.GetServerTLSConfig()
		if err != nil {
			return nil, nil, err
		}
		options = append(options, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	// Create the listener
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start listener")
	}

	// Create a new gRPC server
	server := grpc.NewServer(options...)
	return server, listener, nil
}

// HTTPClientConnInfo creates a client and an url, which the user can use to
// setup the client connection to the server.
func (e *Endpoint) HTTPClientConnInfo(name string) (client *http.Client, u url.URL, err error) {
	if e.TLS == "server" || e.TLS == "mTLS" {
		tlsConfig, err := e.GetClientTLSConfig()
		if err != nil {
			return client, u, err
		}
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	} else {
		// Use default http.RoundTripper transport
		client = &http.Client{}
	}

	if e.TLS != "none" {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Host = e.Address
	u.Path = e.Path
	return
}

/* HTTPServerListener creates a http.Server and a net.Listener to be used by a http server.
...
	ep,_ := xvals.GetEndpoint("my_endpoint")
	server,listener,_ := ep.HTTPServerListener(true)
	server.Handler = http.NotFoundHandler()
	server.Serve(listener)
...
*/
func (e *Endpoint) HTTPServerListener(allInterfaces bool) (*http.Server, net.Listener, error) {
	var (
		listenAddress string
		err           error
		listener      net.Listener
	)
	if allInterfaces {
		_, port, err := net.SplitHostPort(e.Address)
		if err != nil {
			return nil, nil, fmt.Errorf("malformed address in endpoint %s  %w", e.Address, err)
		}
		listenAddress = fmt.Sprintf(":%s", port)
	}
	tlsConfig := new(tls.Config)
	if e.TLS == "server" || e.TLS == "mTLS" {
		tlsConfig, err = e.GetServerTLSConfig()
		if err != nil {
			return nil, nil, err
		}
		listener, err = tls.Listen("tcp", listenAddress, tlsConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to start tls listener")
		}
	} else {
		listener, err = net.Listen("tcp", listenAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to start non-tls listener")
		}
	}

	// Create the listener

	// Create a new http server
	server := &http.Server{}
	return server, listener, nil
}

// ====================== //
//    Utility functions   //
// ====================== //

// fileOrContent checks if val exists as a file and then returns its
// content as a string. Otherwise val is returned
func fileOrContent(val string) []byte {
	if stat, err := os.Stat(val); err == nil {
		if !stat.IsDir() {
			data, err := os.ReadFile(val)
			if err == nil {
				return data
			}
		}
	}
	return []byte(val)
}
