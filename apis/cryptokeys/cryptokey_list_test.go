package servers

import (
	"context"
	"github.com/mittwald/go-powerdns/pdnshttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient_ListCryptokeys(t *testing.T) {
	gock.New("http://dns.example").
		Get("/api/v1/servers/localhost/zones/pdns-test.de/cryptokeys").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`[
  {
    "active": true,
    "algorithm": "ECDSAP256SHA256",
    "bits": 256,
    "dnskey": "257 3 13 sO2Oog47gVFc0iDl0Ubm/RUJ/bdOks/tJmfNS4KX7IPEj2lymwvHBlXqXEvnpsVa+c4CGidwdoGyo7TDMDUIQg==",
    "ds": [
      "50747 13 1 63cdac4d2115c3ea8a8f5d311af58957c2270e32",
      "50747 13 2 336d41f466a29e65118a5d46c02b3680043e8194096e61d07c77931fb49269a8",
      "50747 13 4 03821c4f34a8d63ef80015383d3a5f12ce99e0cb8d8f5a3010b41098fac54f4127d63ea5021f7396bac8c079b6235bf3"
    ],
    "flags": 257,
    "id": 102,
    "keytype": "csk",
    "published": true,
    "type": "Cryptokey"
  }
]`)

	hc := &http.Client{Transport: gock.DefaultTransport}
	c := pdnshttp.NewClient("http://dns.example", hc, &pdnshttp.APIKeyAuthenticator{APIKey: "secret"}, ioutil.Discard)
	cc := New(c)

	list, err := cc.ListCryptokeys(context.Background(), "localhost", "pdns-test.de")

	assert.Nil(t, err)
	require.NotNil(t, list)
	require.Len(t, list, 1)
	assert.Equal(t, 102, list[0].ID)
}
