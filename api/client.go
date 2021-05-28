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
    "strings"
    "fmt"
    "bytes"
    "context"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"

    "golang.org/x/oauth2"
)

// Define end points
const (
    // oauth2 and api flow
    BaseHost = "https://www.upwork.com/"
    DefaultEpoint = "api"
    GqlEndpoint = "https://api.upwork.com/graphql"
    AuthorizationEP = BaseHost + "ab/account-security/oauth2/authorize"
    AccessTokenEP = BaseHost + DefaultEpoint + "/v3/oauth2/token"
    DataFormat = "json"
    OverloadParam = "http_method"
    UPWORK_LIBRARY_USER_AGENT = "Github Upwork API Golang Library"

    // response types
    ByteResponse = "[]byte"
    ErrorResponse = "error"
)

// Api client
type ApiClient struct {
    // oauth2
    oconf *oauth2.Config
    token *oauth2.Token
    oclient *http.Client
    config *Config

    // refresh token notify function
    rnfunc TokenNotifyFunc

    // client
    ep string
    respType string
    sendPostAsJson bool
    hasCustomHttpClient bool
}

// TokenNotifyFunc is a function that accepts an oauth2 Token upon refresh, and
// returns an error if it should not be used.
type TokenNotifyFunc func(*oauth2.Token) error

// NotifyingTokenSource is an oauth2.TokenSource that calls a function when a
// new token is obtained.
type NotifyingTokenSource struct {
    f   TokenNotifyFunc
    src oauth2.TokenSource
}

// Setup client using specific config
func Setup(config *Config) (client ApiClient) {
    var c ApiClient

    c.config = config

    c.oconf = &oauth2.Config{
        ClientID:     config.ClientId,
        ClientSecret: config.ClientSecret,
	RedirectURL: config.RedirectUri,
	Endpoint: oauth2.Endpoint{
	    TokenURL: AccessTokenEP,
	    AuthURL:  AuthorizationEP,
        },
    }

    c.hasCustomHttpClient = config.HasCustomHttpClient

    // Force setup of client_id as a parameter
    oauth2.RegisterBrokenAuthHeaderProvider(BaseHost)

    c.token = new(oauth2.Token)
    c.token.TokenType = "Bearer"
    c.token.AccessToken = config.AccessToken
    c.token.RefreshToken = config.RefreshToken
    c.token.Expiry = config.ExpiresAt

    c.SetApiResponseType(ByteResponse)
    c.SetPostAsJson(false) // send by default using PostForm

    return c
}

// NewNotifyingTokenSource creates a NotifyingTokenSource from an underlying src
// and calls f when a new token is obtained.
func NewNotifyingTokenSource(src oauth2.TokenSource, f TokenNotifyFunc) *NotifyingTokenSource {
    return &NotifyingTokenSource{f: f, src: src}
}

// Token fetches a new token from the underlying source.
func (s *NotifyingTokenSource) Token() (*oauth2.Token, error) {
    t, err := s.src.Token()
    if err != nil {
        return nil, err
    }
    if s.f == nil {
        return t, nil
    }
    return t, s.f(t)
}

// Set refresh token notify function
func (c *ApiClient) SetRefreshTokenNotifyFunc(f TokenNotifyFunc) () {
    c.rnfunc = f
}

// Set request type for non-Get requests
func (c *ApiClient) SetPostAsJson(t bool) {
    c.sendPostAsJson = t;
}

// Set API response type
func (c *ApiClient) SetApiResponseType(t string) {
    c.respType = t;
}

// Set entry point, e.g requested from a router
func (c *ApiClient) SetEntryPoint(ep string) {
    c.ep = ep
}

// Receive an authorization URL
func (c *ApiClient) GetAuthorizationUrl(stateString string) (authzUrl string) {
    url := c.oconf.AuthCodeURL(stateString)
    if url == "" {
        log.Fatal("Can not get authorization URL using OAuth2 library")
    }

    return url
}

// Get access token using a specific authorization code
func (c *ApiClient) GetToken(ctx context.Context, authzCode string) (*oauth2.Token) {
    ctx = c.config.SetOwnHttpClient(ctx)

    accessToken, err := c.oconf.Exchange(ctx, strings.Trim(authzCode, "\n"))
    if err != nil {
        log.Fatal(err)
    }

    c.token = accessToken
    c.setupOauth2Client(ctx)

    return accessToken
}

// Check if client contains already a access/refresh token pair
func (c *ApiClient) HasAccessToken(ctx context.Context) (bool) {
    has := (c.token != nil && (c.token.AccessToken != "" && c.token.RefreshToken != ""))
    if has {
        c.setupOauth2Client(ctx)
    }
    return has
}

// GET method for client
func (c *ApiClient) Get(uri string, params map[string]string) (r *http.Response, re interface{}) {
    // parameters must be encoded according to RFC 3986
    // hmmm, it seems to be the easiest trick?
    qstr := ""
    if params != nil {
        for k, v := range params {
            qstr += fmt.Sprintf("%s=%s&", k, v)
        }
        qstr = qstr[0:len(qstr)-1]
    }
    u := &url.URL{Path: qstr}

    // https://github.com/mrjones/oauth/issues/34
    encQuery := strings.Replace(u.String(), ";", "%3B", -1)
    encQuery = strings.Replace(encQuery, "./", "?", 1) // see URL.String method to understand when "./" is returned

    // non-empty string may miss "?"
    if encQuery !="" && encQuery[:1] != "?" {
        encQuery = "?" + encQuery
    }

    response, err := c.oclient.Get(formatUri(uri, c.ep) + encQuery)

    return c.getTypedResponse(response, err)
}

// POST method for client
func (c *ApiClient) Post(uri string, params map[string]string) (r *http.Response, re interface{}) {
    return c.sendPostRequest(uri, params)
}

// PUT method for client
func (c *ApiClient) Put(uri string, params map[string]string) (r *http.Response, re interface{}) {
    return c.sendPostRequest(uri, addOverloadParam(params, "put"))
}

// DELETE method for client
func (c *ApiClient) Delete(uri string, params map[string]string) (r *http.Response, re interface{}) {
    return c.sendPostRequest(uri, addOverloadParam(params, "delete"))
}

// setup/save authorized oauth2 client, based on received or provided access/refresh token pair
func (c *ApiClient) setupOauth2Client(ctx context.Context) {
    if (c.hasCustomHttpClient == false) {
	ctx = c.config.SetOwnHttpClient(ctx)
    }

    if c.rnfunc != nil {
        // setup notifier for token-refresh workflow - https://github.com/golang/oauth2/issues/84
        realSource := c.oconf.TokenSource(ctx, c.token)
        notifyingSrc := NewNotifyingTokenSource(realSource, c.rnfunc)
        notifyingWithInitialSrc := oauth2.ReuseTokenSource(c.token, notifyingSrc)
        // setup authorized oauth2 client
        c.oclient = oauth2.NewClient(ctx, notifyingWithInitialSrc)
    } else {
        c.oclient = c.oconf.Client(ctx, c.token)
    }
}

// setup X-Upwork-API-TenantId header
func (c *ApiClient) SetOrgUidHeader(ctx context.Context, tenantId string) {
    c.config.SetOrgUidHeader(tenantId)
    c.setupOauth2Client(ctx)
}

// run post/put/delete requests
func (c *ApiClient) sendPostRequest(uri string, params map[string]string) (*http.Response, interface{}) {
    var (
        response *http.Response
        err error
    )

    if c.ep == "graphql" {
        jsonStr, _ := json.Marshal(params) // params contain json data in this case
        response, err = c.oclient.Post(GqlEndpoint, "application/json", bytes.NewBuffer(jsonStr))
    } else if c.sendPostAsJson == true {
	// old style for backward compatibility with the old library
        var jsonStr = []byte("{}")
        if params != nil {
            str := ""
            for k, v := range params {
                str += fmt.Sprintf("\"%s\": \"%s\",", k, v)
            }
            jsonStr = []byte(fmt.Sprintf("{%s}", str[0:len(str)-1]))
        }

        response, err = c.oclient.Post(formatUri(uri, c.ep), "application/json", bytes.NewBuffer(jsonStr))
    } else {
	// prefered
        urlValues := url.Values{}
        if params != nil {
	    for k, v := range params {
                urlValues.Add(k, v)
            }
        }

        response, err = c.oclient.PostForm(formatUri(uri, c.ep), urlValues)
    }

    return c.getTypedResponse(response, err)
}

// return proper response type
func (c *ApiClient) getTypedResponse(resp *http.Response, re error) (*http.Response, interface{}) {
    if c.respType == ByteResponse {
	r, b := formatResponse(resp, re)
	return r, b.([]byte)
    } else {
        return resp, re
    }
}

// Check and format (preparate a byte body) http response routine
func formatResponse(resp *http.Response, err error) (*http.Response, interface{}) {
    if err != nil {
        log.Fatal("Can not execute the request, " + err.Error())
    }

    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        // do not exit, it can be a normal response
        // it's up to client/requester's side decide what to do
    }
    // read json http response
    jsonDataFromHttp, _ := ioutil.ReadAll(resp.Body)

    return resp, jsonDataFromHttp
}

// Create a path to a specific resource
func formatUri(uri string, ep string) (string) {
    format := ""
    if ep == DefaultEpoint {
        format += "." + DataFormat
    }
    return BaseHost + ep + uri + format
}

// add overload parameter to the map of parameters
func addOverloadParam(params map[string]string, op string) map[string]string {
    if params == nil {
        params = make(map[string]string)
    }
    params[OverloadParam] = op
    return params
}
