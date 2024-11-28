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
package milestones

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

// Get active Milestone for the Contract
func (r a) GetActiveMilestone(contractId string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Get all submissions for the active Milestone
func (r a) GetSubmissions(milestoneId string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Create a new Milestone
func (r a) Create(params map[string]string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Edit an existing Milestone
func (r a) Edit(milestoneId string, params map[string]string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Activate an existing Milestone
func (r a) Activate(milestoneId string, params map[string]string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Approve an existing Milestone
func (r a) Approve(milestoneId string, params map[string]string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}

// Delete an existing Milestone
func (r a) Delete(milestoneId string) (*http.Response, interface{}) {
	panic("The legacy API was deprecated. Please, use GraphQL call - see example in this library.")
}
