package api

import (
    "time"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestImports(t *testing.T) {
    if assert.Equal(t, 1, 1) != true {
        t.Error("Something is wrong.")
    }
}

func TestNewConfig(t *testing.T) {
    settings := map[string]string{
        "client_id": "consumerkey",
        "client_secret": "consumersecret",
        "access_token": "accesstoken",
        "refresh_token": "refreshtoken",
	"redirect_uri": "http://a.redirect.uri",
	"expires_at": "2018-01-01T01:00:00.000Z",
	"expires_in": "100",
        "debug": "on",
    }
    config := NewConfig(settings)

    if assert.NotNil(t, config) {
        assert.Equal(t, "consumerkey", config.ClientId)
        assert.Equal(t, "consumersecret", config.ClientSecret)
        assert.Equal(t, "accesstoken", config.AccessToken)
        assert.Equal(t, "refreshtoken", config.RefreshToken)
        assert.Equal(t, "http://a.redirect.uri", config.RedirectUri)
        ttime, _ := time.Parse(TIMEFORMAT, "2018-01-01T01:00:00.000Z")
        assert.Equal(t, ttime, config.ExpiresAt)
        assert.Equal(t, "100", config.ExpiresIn)
        assert.True(t, true, config.Debug)
    }
}

func TestReadConfig(t *testing.T) {
    config := ReadConfig("../example/config.json")

    if assert.NotNil(t, config) {
        assert.Equal(t, "YOUR_CONSUMER_KEY", config.ClientId)
        assert.Equal(t, "YOUR_CONSUMER_SECRET", config.ClientSecret)
        assert.Equal(t, "access-token-if-known-otherwise-remove", config.AccessToken)
        assert.Equal(t, "access-token-if-known-otherwise-remove", config.RefreshToken)
        assert.Equal(t, "https://a.callback.url", config.RedirectUri)
        ttime, _ := time.Parse(TIMEFORMAT, "2018-01-01T01:00:00.000Z")
        assert.Equal(t, ttime, config.ExpiresAt)
    }
}
