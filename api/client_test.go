package api

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestClientConstants(t *testing.T) {
    assert.Equal(t, "https://www.upwork.com/", BaseHost)
    assert.Equal(t, "api", DefaultEpoint)
    assert.Equal(t, "https://www.upwork.com/ab/account-security/oauth2/authorize", AuthorizationEP)
    assert.Equal(t, "https://www.upwork.com/api/v3/oauth2/token", AccessTokenEP)
    assert.Equal(t, "json", DataFormat)
    assert.Equal(t, "http_method", OverloadParam)
}

func TestSetup(t *testing.T) {
    client := Setup(ReadConfig("../example/config.json"))
    if assert.NotNil(t, client) {
        assert.NotNil(t, client.oconf)
        assert.NotNil(t, client.token)
    }
}

func TestSetEntryPoint(t *testing.T) {
    client := Setup(ReadConfig("../example/config.json"))
    client.SetEntryPoint("gds")

    assert.Equal(t, "gds", client.ep)
}
