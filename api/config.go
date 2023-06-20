// Package implements access to Upwork API
//
// Licensed under the Upwork's API Terms of Use;
// you may not use this file except in compliance with the Terms.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author::    Maksym Novozhylov (mnovozhilov@upwork.com)
// Copyright:: Copyright 2018(c) Upwork.com
// License::   See LICENSE.txt and TOS - https://developers.upwork.com/api-tos.html
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	TIMEFORMAT = "2006-01-02T15:04:05.000Z" // NOTE: time.RFC3339 does not work for unclear reason?
)

type HeadersTransport struct {
	rt              http.RoundTripper
	uaHeader        string
	xTenantIdHeader string
}

// Config
type Config struct {
	ClientId            string
	ClientSecret        string
	AccessToken         string
	RefreshToken        string
	RedirectUri         string
	ExpiresIn           string
	ExpiresAt           time.Time
	State               string
	GrantType           string
	Debug               bool
	HasCustomHttpClient bool
	TenantIdHeader      string // X-Upwork-API-TenantId required for GraphQL requests
}

// List of required configuration keys
var requiredKeys = [2]string{"client_id", "client_secret"}

// Create a new config
func NewConfig(data map[string]string) (settings *Config) {
	cfg := &Config{
		ClientId:     data["client_id"],
		ClientSecret: data["client_secret"],
		RedirectUri:  data["redirect_uri"],
	}

	// save access token if defined
	if val, ok := data["access_token"]; ok {
		cfg.AccessToken = val
	}

	// save refresh token if defined
	if val, ok := data["refresh_token"]; ok {
		cfg.RefreshToken = val
	}

	// save expires_in if defined
	if val, ok := data["expires_in"]; ok {
		cfg.ExpiresIn = val
	}

	// save expiresat if defined
	if val, ok := data["expires_at"]; ok {
		cfg.ExpiresAt, _ = time.Parse(TIMEFORMAT, val)
	}

	// save state if defined
	if val, ok := data["state"]; ok {
		cfg.State = val
	}

	// save grant_type if defined
	if val, ok := data["grant_type"]; ok {
		cfg.GrantType = val
	}

	// save debug flag if defined
	if debug, ok := data["debug"]; ok && debug == "on" {
		cfg.Debug = true
	}

	return cfg
}

// Read a specific configuration (json) file
func ReadConfig(fn string) (settings *Config) {
	// read from config file if exists
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal("config file: ", err)
	}

	// parse json config
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		log.Fatal("config file: ", err)
	}

	_, ok := data["redirect_uri"]
	if data["grant_type"] != "client_credentials" && !ok {
		log.Fatal("config file: redirect_uri is missing in " + fn)
	}

	// test required properties
	for _, v := range requiredKeys {
		_, ok := data[v]
		if !ok {
			log.Fatal("config file: " + v + " is missing in " + fn)
		}
	}

	// convert
	config := make(map[string]string)
	for k, v := range data {
		config[k] = v.(string)
	}

	return NewConfig(config)
}

// RoundTrip for the RoundTripper interface
func (t *HeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.uaHeader)

	if t.xTenantIdHeader != "" {
		req.Header.Set("X-Upwork-API-TenantId", t.xTenantIdHeader)
	}

	return t.rt.RoundTrip(req)
}

// Configure X-Upwork-API-TenantId header for OwnHttpClient
func (cfg *Config) SetOrgUidHeader(tenantId string) {
	if cfg.HasCustomHttpClient == true {
		panic("SetOrgUidHeader can not be used with the custom client. Add X-Upwork-API-TenantId header manually.")
	}
	cfg.TenantIdHeader = tenantId
}

// Configure a context with the custom http client
func (cfg *Config) SetCustomHttpClient(ctx context.Context, httpClient *http.Client) context.Context {
	cfg.HasCustomHttpClient = true
	return context.WithValue(ctx, oauth2.HTTPClient, httpClient)
}

// Configure a context with the own http client over RoundTripper
func (cfg *Config) SetOwnHttpClient(ctx context.Context) context.Context {
	cfg.HasCustomHttpClient = false

	// Prepare wrapper to fix User-Agent header
	var transport http.RoundTripper = &HeadersTransport{http.DefaultTransport, UPWORK_LIBRARY_USER_AGENT, cfg.TenantIdHeader}

	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: transport})
}

// Test print of found/assigned key
func (cfg *Config) Print() {
	fmt.Println("assigned client id (key):", cfg.ClientId)
}
