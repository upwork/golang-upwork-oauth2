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
package metadata

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

// Get categories (V2)
func (r a) GetCategoriesV2() (*http.Response, interface{}) {
    return r.client.Get("/profiles/v2/metadata/categories", nil)
}

// Get skills
func (r a) GetSkills() (*http.Response, interface{}) {
    return r.client.Get("/profiles/v1/metadata/skills", nil)
}

// Get skills (V2)
func (r a) GetSkillsV2() (*http.Response, interface{}) {
    return r.client.Get("/profiles/v2/metadata/skills", nil)
}

// Get specialties
func (r a) GetSpecialties() (*http.Response, interface{}) {
    return r.client.Get("/profiles/v1/metadata/specialties", nil)
}

// Get regions
func (r a) GetRegions() (*http.Response, interface{}) {
    return r.client.Get("/profiles/v1/metadata/regions", nil)
}

// Get tests
func (r a) GetTests() (*http.Response, interface{}) {
    return r.client.Get("/profiles/v1/metadata/tests", nil)
}

// Get reasons
func (r a) GetReasons(params map[string]string) (*http.Response, interface{}) {
    return r.client.Get("/profiles/v1/metadata/reasons", params)
}
