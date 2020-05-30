package cryptokeys

import (
	"context"
	"github.com/mittwald/go-powerdns/pdnshttp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient_CreateCryptokey(t *testing.T) {
	gock.New("http://dns.example").
		Post("/api/v1/servers/localhost/zones/pdns-test.de/cryptokeys").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
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
  "privatekey": "Private-key-format: v1.2\nAlgorithm: 13 (ECDSAP256SHA256)\nPrivateKey: 4Xt/Qdsasn/TBC3O/PVCIEO4c2NozvRpX50qVdEL/Ag=\n",
  "published": true,
  "type": "Cryptokey"
}`)

	hc := &http.Client{Transport: gock.DefaultTransport}
	c := pdnshttp.NewClient("http://dns.example", hc, &pdnshttp.APIKeyAuthenticator{APIKey: "secret"}, ioutil.Discard)
	cc := New(c)

	key, err := cc.CreateCryptokey(context.Background(), "localhost", "pdns-test.de")

	assert.Nil(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, "ECDSAP256SHA256", key.Algorithm)
	assert.Equal(t, 256, key.Bits)
}
