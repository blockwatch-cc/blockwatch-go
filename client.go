// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package blockwatch

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// Defines the default server URL.
	dataURL = "https://data.blockwatch.cc/v1/"

	// Defines the release version of this SDK.
	sdkVersion = "1.0.0"

	// Defines the expected media type.
	mediaType = "application/json"
)

var (
	// Defines the User Agent string sent to HTTP servers. May be overwritten
	// by clients.
	UserAgent = "Blockwatch-Data-SDK/" + sdkVersion

	// ERootCAFailed is the error returned when a failure to add a provided
	// root CA to the list of known CAs occurs.
	ERootCAFailed = errors.New("failed to add Root CAs to certificate pool")
)

// ConnConfig describes the connection configuration parameters for the client.
type ConnConfig struct {
	// HTTP tuning parameters
	DialTimeout           time.Duration `json:"dial_timeout"`
	KeepAlive             time.Duration `json:"keepalive"`
	IdleConnTimeout       time.Duration `json:"idle_timeout"`
	ResponseHeaderTimeout time.Duration `json:"response_timeout"`
	ExpectContinueTimeout time.Duration `json:"continue_timeout"`
	MaxIdleConns          int           `json:"idle_conns"`

	// Proxy specifies to connect through a SOCKS 5 proxy server.  It may
	// be an empty string if a proxy is not required.
	Proxy string `json:"proxy"`

	// ProxyUser is an optional username to use for the proxy server if it
	// requires authentication.  It has no effect if the Proxy parameter
	// is not set.
	ProxyUser string `json:"proxy_user"`

	// ProxyPass is an optional password to use for the proxy server if it
	// requires authentication.  It has no effect if the Proxy parameter
	// is not set.
	ProxyPass string `json:"proxy_pass"`

	// TLS configuration options
	ServerName         string   `json:"server_name"`
	AllowInsecureCerts bool     `json:"disable_tls"`
	TLSMinVersion      int      `json:"tls_min_version"`
	TLSMaxVersion      int      `json:"tls_max_version"`
	RootCaCerts        []string `json:"tls_ca"`
	RootCaCertsFile    string   `json:"tls_ca_file"`
	ClientCert         []string `json:"tls_cert"`
	ClientCertFile     string   `json:"tls_cert_file"`
	ClientKey          []string `json:"tls_key"`
	ClientKeyFile      string   `json:"tls_key_file"`
}

// sane defaults
func DefaultConnConfig() *ConnConfig {
	return &ConnConfig{
		DialTimeout:           5 * time.Second,
		KeepAlive:             180 * time.Second,
		IdleConnTimeout:       180 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		MaxIdleConns:          2,
	}
}

type transportCloser interface {
	CloseIdleConnections()
}

type Client struct {
	config     ConnConfig
	baseURL    *url.URL
	apikey     string
	userAgent  string
	httpClient *http.Client
}

// NewClient creates a new API client based on the provided connection configuration.
func NewClient(apikey string, config *ConnConfig) (*Client, error) {
	if apikey == "" {
		return nil, fmt.Errorf("missing Blockwatch Data API key")
	}
	if config == nil {
		config = DefaultConnConfig()
	}

	u, err := url.Parse(dataURL)
	if err != nil {
		return nil, err
	}

	httpClient, err := newHTTPClient(config)
	if err != nil {
		return nil, err
	}

	c := &Client{
		config:     *config,
		baseURL:    u,
		httpClient: httpClient,
		apikey:     apikey,
		userAgent:  UserAgent,
	}
	return c, nil
}

func (c *Client) Get(ctx context.Context, urlpath string, headers http.Header, result interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, urlpath, headers, nil, result)
	if err != nil {
		return err
	}
	return c.do(req, result)
}

func (c *Client) newRequest(ctx context.Context, method, urlStr string, headers http.Header, data, result interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	var body io.Reader
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("X-API-Key", c.apikey)

	if body != nil {
		req.Header.Add("Content-Type", mediaType)
	}

	if result != nil && headers.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	for n, v := range headers {
		for _, vv := range v {
			req.Header.Add(n, vv)
		}
	}

	// handle proxy auth
	if c.config.ProxyUser != "" && c.config.ProxyPass != "" {
		headers.Set("Proxy-Authorization", BasicAuth(c.config.ProxyUser, c.config.ProxyPass))
	}

	return req, nil
}

// Do retrieves values from the API and marshals them into the provided interface.
func (c *Client) do(req *http.Request, v interface{}) (err error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	statusClass := resp.StatusCode / 100
	if statusClass == 2 {
		if v == nil {
			return nil
		}
		return c.handleResponse(req.Context(), resp, v)
	}

	return handleError(resp)
}

func (c *Client) handleResponse(ctx context.Context, resp *http.Response, v interface{}) error {
	// process as stream when response interface is an io.Writer
	if stream, ok := v.(io.Writer); ok {
		_, err := io.Copy(stream, resp.Body)
		// close consumer if possible
		if closer, ok := v.(io.WriteCloser); ok {
			closer.Close()
		}
		return err
	}

	// process other responses as JSON
	return json.NewDecoder(resp.Body).Decode(v)
}

func handleError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	httpErr := httpError{
		status: resp.StatusCode,
		body:   bytes.Trim(body, "\n"),
	}

	// prepare special rate limit error
	if resp.StatusCode == 429 {
		resetTime, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)
		return newErrRateLimited(&httpErr, resetTime)
	}

	// unpack response errors
	var errs Errors
	json.Unmarshal(body, &errs)
	return apiError{
		httpError: &httpErr,
		errors:    errs,
	}
}

// newHTTPClient returns a new http client that is configured according to the
// proxy, TLS and timeout settings in the associated connection configuration.
func newHTTPClient(cc *ConnConfig) (*http.Client, error) {
	// Set proxy function if there is a proxy configured.
	proxyFunc := http.ProxyFromEnvironment
	if cc.Proxy != "" {
		proxyURL, err := url.Parse(cc.Proxy)
		if err != nil {
			return nil, err
		}
		proxyFunc = http.ProxyURL(proxyURL)
	}

	// Configure TLS if needed.
	tlsConfig, err := makeTLSConfig(cc)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   cc.DialTimeout,
				KeepAlive: cc.KeepAlive,
			}).Dial,
			Proxy:                 proxyFunc,
			TLSClientConfig:       tlsConfig,
			IdleConnTimeout:       cc.IdleConnTimeout,
			ResponseHeaderTimeout: cc.ResponseHeaderTimeout,
			ExpectContinueTimeout: cc.ExpectContinueTimeout,
			MaxIdleConns:          cc.MaxIdleConns,
			MaxIdleConnsPerHost:   cc.MaxIdleConns,
		},
	}
	return client, nil
}

func tlsVersion(v int) uint16 {
	switch v {
	// case 0:
	//  // insecure
	// 	return tls.VersionSSL30
	case 1:
		return tls.VersionTLS10
	case 2:
		return tls.VersionTLS11
	case 3:
		return tls.VersionTLS12
	default:
		// default to strongest TLS
		return tls.VersionTLS12
	}
}

func makeTLSConfig(cc *ConnConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cc.AllowInsecureCerts,
		ServerName:         cc.ServerName,
		MinVersion:         tlsVersion(cc.TLSMinVersion),
		MaxVersion:         tlsVersion(cc.TLSMaxVersion),
	}
	if len(cc.RootCaCerts) > 0 {
		// load from config
		rootCAs := x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM([]byte(strings.Join(cc.RootCaCerts, "\n"))) {
			return nil, ERootCAFailed
		}
		tlsConfig.RootCAs = rootCAs
	} else if len(cc.RootCaCertsFile) > 0 {
		// load from file
		caCert, err := ioutil.ReadFile(cc.RootCaCertsFile)
		if err != nil {
			return nil, fmt.Errorf("Could not load TLS CA [%s]: %v", cc.RootCaCertsFile, err)
		}
		rootCAs := x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(caCert) {
			return nil, ERootCAFailed
		}
		tlsConfig.RootCAs = rootCAs
	}
	if len(cc.ClientCert) > 0 && len(cc.ClientKey) > 0 {
		// load from config
		cert, err := tls.X509KeyPair(
			[]byte(strings.Join(cc.ClientCert, "\n")),
			[]byte(strings.Join(cc.ClientKey, "\n")),
		)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	} else if len(cc.ClientCertFile) > 0 && len(cc.ClientKeyFile) > 0 {
		// load from file
		cert, err := tls.LoadX509KeyPair(cc.ClientCertFile, cc.ClientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("Could not load TLS client cert or key [%s]: %v", cc.ClientCertFile, err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	}
	return tlsConfig, nil
}

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
	auth := strings.Join([]string{username, password}, ":")
	return strings.Join([]string{"Basic", base64.StdEncoding.EncodeToString([]byte(auth))}, " ")
}
