// Router for Upwork API
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
// Copyright:: Copyright 2021(c) Upwork.com
// License::   See LICENSE.txt and TOS - https://developers.upwork.com/api-tos.html
package graphql

import (
    "net/http"
    "github.com/upwork/golang-upwork-oauth2/api"
)

const (
    EntryPoint = "graphql"
)

type a struct {
    client *api.ApiClient
}

// Constructor
func New(c *api.ApiClient) *a {
    c.SetEntryPoint(EntryPoint)

    return &a{c}
}

// Execute GraphQL request
func (r a) Execute(jsonData map[string]string) (*http.Response, interface{}) {
    return r.client.Post("", jsonData)
}
