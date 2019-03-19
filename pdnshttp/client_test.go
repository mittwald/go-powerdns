package pdnshttp

import (
	"context"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetExecutedCorrectly(t *testing.T) {
	gock.New("http://test.example").
		Get("/api/v1/servers").
		MatchHeader("X-API-Key", "secret").
		Reply(http.StatusOK).
		JSON(map[string]string{"foo": "bar"})

	hc := &http.Client{Transport: gock.DefaultTransport}
	c := NewClient("http://test.example", hc, &APIKeyAuthenticator{APIKey: "secret"}, ioutil.Discard)

	var out interface{}

	err := c.Get(context.Background(), "/api/v1/servers2", &out)
	require.Nil(t, err)
	require.True(t, gock.IsDone(), "still has pending mocks")
}
