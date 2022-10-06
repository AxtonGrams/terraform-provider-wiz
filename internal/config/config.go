package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Settings holds all the information necessary to configure the provider
type Settings struct {
	WizURL                 string
	WizAuthURL             string
	WizAuthGrantType       string
	WizAuthClientID        string
	WizAuthClientSecret    string
	WizAuthAudience        string
	Proxy                  bool
	ProxyServer            string
	CAChain                string
	HTTPClientRetryMax     int
	HTTPClientRetryWaitMin int
	HTTPClientRetryWaitMax int
}

// ProviderConf holds structures that are useful to the provider at runtime
type ProviderConf struct {
	Settings   *Settings
	TokenType  string
	Token      string
	HTTPClient *http.Client
	UserAgent  string
}

// AuthorizationResponse contains the reponse from the authorization api
type AuthorizationResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GetHTTPClient creates a http client
func GetHTTPClient(ctx context.Context, settings *Settings) *http.Client {
	tflog.Info(ctx, "GetHTTPClient called...")

	// load trusted certificate authorities in a certpool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(settings.CAChain))

	// configure the transport with trusted certificate authorities
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	transport := &http.Transport{
		TLSClientConfig:   tlsConfig,
		MaxConnsPerHost:   10,
		DisableKeepAlives: false,
	}
	if settings.Proxy {
		proxyURL, _ := url.Parse(settings.ProxyServer)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// configure the client
	client := retryablehttp.NewClient()

	// override with the trusted certificate authorities
	client.HTTPClient.Transport = transport
	client.RetryWaitMin = time.Duration(settings.HTTPClientRetryWaitMin) * 1000000000
	client.RetryWaitMax = time.Duration(settings.HTTPClientRetryWaitMax) * 1000000000
	client.RetryMax = settings.HTTPClientRetryMax

	return client.StandardClient()
}

// NewProviderConf creates a new structure containing all configuration data
func NewProviderConf(ctx context.Context, settings *Settings, userAgent string) (*ProviderConf, diag.Diagnostics) {
	tflog.Info(ctx, "NewProviderConf called...")

	tokenType, token, diags := GetSessionToken(ctx, settings)

	pcfg := &ProviderConf{
		Settings:   settings,
		Token:      token,
		TokenType:  tokenType,
		HTTPClient: GetHTTPClient(ctx, settings),
		UserAgent:  userAgent,
	}
	return pcfg, diags
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) (*Settings, error) {
	cfg := &Settings{
		WizURL:                 d.Get("wiz_url").(string),
		WizAuthURL:             d.Get("wiz_auth_url").(string),
		WizAuthGrantType:       d.Get("wiz_auth_grant_type").(string),
		WizAuthClientID:        d.Get("wiz_auth_client_id").(string),
		WizAuthClientSecret:    d.Get("wiz_auth_client_secret").(string),
		WizAuthAudience:        d.Get("wiz_auth_audience").(string),
		Proxy:                  d.Get("proxy").(bool),
		ProxyServer:            d.Get("proxy_server").(string),
		CAChain:                d.Get("ca_chain").(string),
		HTTPClientRetryMax:     d.Get("http_client_retry_max").(int),
		HTTPClientRetryWaitMin: d.Get("http_client_retry_wait_min").(int),
		HTTPClientRetryWaitMax: d.Get("http_client_retry_wait_max").(int),
	}

	return cfg, nil
}

// GetSessionToken retrieves a new session token
func GetSessionToken(ctx context.Context, settings *Settings) (string, string, diag.Diagnostics) {
	tflog.Info(ctx, "GetSessionToken called...")

	var diags diag.Diagnostics

	// get an http client
	httpclient := GetHTTPClient(ctx, settings)

	// setup the request
	data := url.Values{}
	data.Set("grant_type", settings.WizAuthGrantType)
	data.Set("client_id", settings.WizAuthClientID)
	data.Set("client_secret", settings.WizAuthClientSecret)
	data.Set("audience", settings.WizAuthAudience)
	request, err := http.NewRequestWithContext(ctx, "POST", settings.WizAuthURL, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		return "", "", append(diags, diag.FromErr(err)...)
	}

	// log the request
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return "", "", append(diags, diag.FromErr(err)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("authentication request: %s", reqDump))

	// call the api
	resp, err := httpclient.Do(request)
	if err != nil {
		return "", "", append(diags, diag.FromErr(err)...)
	}

	// log the response
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return "", "", append(diags, diag.FromErr(err)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("auth response: %s", respDump))

	defer resp.Body.Close()

	// parse the response
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", append(diags, diag.FromErr(err)...)
	}

	// validate successful response
	responseBody := &AuthorizationResponse{}
	if resp.StatusCode != http.StatusOK {
		return "", "", append(diags, diag.FromErr(err)...)
	}
	err = json.Unmarshal(rbody, &responseBody)
	if err != nil {
		return "", "", append(diags, diag.FromErr(err)...)
	}

	// return
	return responseBody.TokenType, responseBody.AccessToken, diags
}
