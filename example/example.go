// Example shows how to work with Upwork API
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
package main

import (
    "fmt"
    "bufio"
    "os"
    "context" // uncomment if you need to setup a custom http client
    _ "net/http" // uncomment if you need to setup a custom http client

    _ "golang.org/x/oauth2" // uncomment if you need to work with oauth2.Token or other object, e.g. to store or re-cache token pair

    "github.com/upwork/golang-upwork-oauth2/api"
    "github.com/upwork/golang-upwork-oauth2/api/routers/auth"
    _ "github.com/upwork/golang-upwork-oauth2/api/routers/messages" // uncomment to test messages example
    _ "github.com/upwork/golang-upwork-oauth2/api/routers/graphql" // uncomment to test graphql example
)

const cfgFile = "config.json" // update the path to your config file, or provide properties directly in your code

func main() {
    // init context
    ctx := context.Background()

/* it is possible to set up properties from code
    settings := map[string]string{
        "client_id": "clientid",
        "client_secret": "clientsecret",
    }
    config := api.NewConfig(settings)

    //or read them from a specific configuration file
    config := api.ReadConfig(cfgFile)
    config.Print()
*/

/* it is possible to setup a custom http client if needed
    httpClient := &http.Client{Timeout: 2}
    config := api.ReadConfig(cfgFile)
    ctx = config.SetCustomHttpClient(ctx, httpClient)
    client := api.Setup(config)
*/
    ctx = context.TODO() // define NoContext if you do not use a custom client, otherwise use earlier defined context
    client := api.Setup(api.ReadConfig(cfgFile))
    // You can configure the package send the requests as application/json, by default PostForm is used.
    // This will be automatically set to true for GraphQL request
    // client.SetPostAsJson(true)

    // GraphQL requests require X-Upwork-API-TenantId header, which can be setup using the following method
    // client.SetOrgUidHeader(ctx, "1234567890") // Organization UID

/*
    // WARNING: oauth2 library will refresh the access token for you
    // Setup notify function for refresh-token-workflow
    // type Token, see https://godoc.org/golang.org/x/oauth2#Token
    f := func(t *oauth2.Token) error {
        // re-cache refreshed token
        _, err := fmt.Printf("The token has been refreshed, here is a new one: %#v\n", t)
        return err
    }
    client.SetRefreshTokenNotifyFunc(f)
*/
    // we need an access/refresh token pair in case we haven't received it yet
    if !client.HasAccessToken(ctx) {
	// required to authorize the application. Once you have an access/refresh token pair associated with
	// the user, no need to redirect to the authorization screen.
        aurl := client.GetAuthorizationUrl("random-state")

        // read code
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("Visit the authorization url and provide oauth_verifier for further authorization")
        fmt.Println(aurl)
        authzCode, _ := reader.ReadString('\n')

	// WARNING: be sure to validate FormValue("state") before getting access token

        // get access token
        token := client.GetToken(ctx, authzCode)
	fmt.Println(token) // type Token, see https://godoc.org/golang.org/x/oauth2#Token
    }

    // http.Response and specified type will be returned, you can use any
    // use client.SetApiResponseType to specify the response type: use api.ByteResponse
    // or api.ErrorResponse, see usage example below
    // by default api.ByteResponse is used, i.e. []byte is returned as second value
    _, jsonDataFromHttp1 := auth.New(&client).GetUserInfo()

    // here you can Unmarshal received json string, or do any other action(s) if you used ByteResponse
    fmt.Println(string(jsonDataFromHttp1.([]byte))) // []byte

    // if you used ErrorResponse, like
    // client.SetApiResponseType(api.ErrorResponse)
    // httpResponse, err := auth.New(&client).GetUserInfo()
    // if err == nil {
    //     ... do smth with http.Response
    // }

    // run a post request using parameters as an example
    // params := make(map[string]string)
    // params["story"] = `{"message": "test message", "userId": "~017xxxxx"}`
    // _, jsonDataFromHttp2 := messages.New(&client).SendMessageToRoom("company_id", "room_id", params)
    // fmt.Println(string(jsonDataFromHttp2.([]byte)))

    // getting reports example
    // params := make(map[string]string)
    // params["tq"] = "select memo where worked_on >= '05-08-2015'"
    // params["tqx"] = "out:json"
    // _, jsonDataFromHttp3 := timereports.New(&client).GetByFreelancerFull(params)
    // fmt.Println(string(jsonDataFromHttp3.([]byte)))

    // sending GraphQL request
    // jsonData := map[string]string{
    //     "query": `
    //       {
    //         user {
    //           id
    //           nid
    //         }
    //         organization {
    //           id
    //         }
    //      }
    //    `,
    //  }
    // _, jsonDataFromHttp4 := graphql.New(&client).Execute(jsonData)
    // fmt.Println(string(jsonDataFromHttp4.([]byte)))
}
