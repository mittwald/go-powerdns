package cryptokeys

import (
	"context"
	"github.com/mittwald/go-powerdns/pdnshttp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io"
	"net/http"
	"testing"
)

func TestClient_DeleteCryptokey(t *testing.T) {
	gock.New("http://dns.example").
		Delete("/api/v1/servers/localhost/zones/pdns-test.de/cryptokeys/102").
		Reply(http.StatusNoContent)

	hc := &http.Client{Transport: gock.DefaultTransport}
	c := pdnshttp.NewClient("http://dns.example", hc, &pdnshttp.APIKeyAuthenticator{APIKey: "secret"}, io.Discard)
	cc := New(c)

	err := cc.DeleteCryptokey(context.Background(), "localhost", "pdns-test.de", 102)

	assert.Nil(t, err)
}
