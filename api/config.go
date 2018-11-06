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
    "log"
    "fmt"
    "time"
    "context"
    "encoding/json"
    "io/ioutil"
    "net/http"

    "golang.org/x/oauth2"
)

const (
    TIMEFORMAT = "2006-01-02T15:04:05.000Z" // NOTE: time.RFC3339 does not work for unclear reason?
)

// Config
type Config struct {
    ClientId string
    ClientSecret string
    AccessToken string
    RefreshToken string
    RedirectUri string
    ExpiresIn string
    ExpiresAt time.Time
    State string
    Debug bool
}

// List of required configuration keys
var requiredKeys = [3]string{"client_id", "client_secret", "redirect_uri"}

// Create a new config
func NewConfig(data map[string]string) (settings *Config) {
    cfg := &Config{
        ClientId: data["client_id"],
        ClientSecret: data["client_secret"],
	RedirectUri: data["redirect_uri"],
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

// Configure a context with the custom http client
func (cfg *Config) SetCustomHttpClient(ctx context.Context, httpClient *http.Client) context.Context {
    return context.WithValue(ctx, oauth2.HTTPClient, httpClient)
}

// Test print of found/assigned key
func (cfg *Config) Print() {
    fmt.Println("assigned client id (key):", cfg.ClientId)
}
