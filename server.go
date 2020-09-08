package puppetdb

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type body io.Reader

/*
Server Representation of a PuppetDB server instance.

Use NewServer to create a new instance.
*/
type Server struct {
	BaseURL           string
	CACertificateFile string
	HTTPTransport     http.RoundTripper
	HTTPTimeout       time.Duration
	Headers           map[string]string
	Body              body
}

// SetHTTPTimeout to set custom Timeout of http.Client
func (s *Server) SetHTTPTimeout(t time.Duration) {
	s.HTTPTimeout = t
}

// SetHeader the header
func (s *Server) SetHeader(key string, value string) {
	s.Headers[key] = value
}

// SetToken - Sets up Puppet Enterprise RBAC Token Header
func (s *Server) SetToken(token string) {
	s.SetHeader("X-Authentication", token)
}

// Authenticate - get Puppet Token from RBAC (using puppet-access) for Puppet Server
func (s *Server) Authenticate() {
	// See if the user supplied a token from the ENV or cli if not, try to fetch an existing one from disk or attempt to create it
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	tokenDir := homeDir + "/.puppetlabs/token"

	// If the token already exists in the user's home dir read it, otherwise create it
	if _, err := os.Stat(tokenDir); os.IsNotExist(err) {
		user, err := user.Current()
		if err != nil {
			log.Fatal("Unable to determine user")
		}
		fmt.Printf("Authenticate using puppet-access...\n")
		cmd := exec.Command("/opt/puppetlabs/bin/puppet-access", "login",
			"--username", user.Username,
			"--lifetime", "4h",
			"--service-url", fmt.Sprintf("https://%s:4433/rbac-api", s.puppetServer()))
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	token, err := ioutil.ReadFile(tokenDir)
	if err != nil {
		log.Fatalf("Unable to read %s and no token was specified or abled to be created", tokenDir)
	}
	s.SetToken(string(token))
	// Validate by checking PuppetDB version
	ver, _ := s.QueryVersion()
	log.Debug("PuppetDB Version: " + ver.Version)
	if len(ver.Version) < 1 {
		log.Info("Token expired or invalid, removing token, re-authenticating...")
		cmd2 := exec.Command("/opt/puppetlabs/bin/puppet-access", "delete-token-file")
		cmd2.Stdout = os.Stdout
		cmd2.Stdin = os.Stdin
		cmd2.Stderr = os.Stdout
		err = cmd2.Run()
		s.Authenticate()
	}
}

// SetCACertificate sets CA Cert
func (s *Server) SetCACertificate(cacert string) {
	s.CACertificateFile = cacert

}

/*
NewServer Create a new instance of a Server for usage later.

This is usually the main entry point of this SDK, where you would create
this initial object and use it to perform activities on the instance in
question.
*/
func NewServer(baseURL string) Server {
	return newServer(baseURL, nil)
}

/*
NewServerWithTransport Create a new instance of a Server for usage later.

Comparable to NewServer, but with an additional parameter to specify the http transport
(i.e. SSL options)
*/
func NewServerWithTransport(baseURL string, httpTransport http.RoundTripper) Server {
	return newServer(baseURL, httpTransport)
}

/*
NewSSLServer - Create new instance of a server with SSL

Secondary entry point of SDK, in case you are using SSL with PuppetDB
*/
func NewSSLServer(baseURL string, cacert string) Server {
	server := newServer(baseURL, nil)
	server.SetCACertificate(cacert)
	server.setTransport()
	return server
}

func (s *Server) setTransport() {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	// Read in the cert file
	certs, err := ioutil.ReadFile(s.CACertificateFile)
	if err != nil {
		log.Fatalf("Failed to read %q : %v", s.CACertificateFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		RootCAs: rootCAs,
	}
	s.HTTPTransport = &http.Transport{TLSClientConfig: config}
}

// puppetServer - Parse PuppetServer from URL for ulterior motives
func (s *Server) puppetServer() string {
	// req, err := http.NewRequest("GET", s.BaseURL, nil)
	// if err != nil {
	// 	log.Panic("PuppetServer hostname could not be determined")
	// }
	// return req.URL.Host
	return strings.Split(strings.Split(s.BaseURL, "/")[2], ":")[0] // janky, but it will work for now
}

func newServer(baseURL string, httpTransport http.RoundTripper) Server {
	return Server{
		BaseURL:       baseURL,
		HTTPTransport: httpTransport,
		HTTPTimeout:   time.Second * 30,
		Headers:       make(map[string]string),
	}
}
