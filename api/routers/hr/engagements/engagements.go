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
// Copyright:: Copyright 2015(c) Upwork.com
// License::   See LICENSE.txt and TOS - https://developers.upwork.com/api-tos.html
package engagements

import (
	"net/http"

	"github.com/upwork/golang-upwork-oauth2/api"
)

const (
	EntryPoint = "api"
)

type a struct {
	client *api.ApiClient
}

// Constructor
func New(c *api.ApiClient) *a {
	c.SetEntryPoint(EntryPoint)

	return &a{c}
}

// Get list of engagements
func (r a) GetList(params map[string]string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Get specific engagement
func (r a) GetSpecific(reference string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}
