package zones

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

func TestCreateZoneCreatesZone(t *testing.T) {
	gock.New("http://dns.example").
		Post("/api/v1/servers/localhost/zones").
		Reply(http.StatusCreated).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
			"account": "", 
			"api_rectify": false,
			"dnssec": false, 
			"id": "some-generated-id",
			"kind": "Native",
			"last_check": 0,
			"masters": [],
			"name": "test.example.",
			"notified_serial": 0,
			"nsec3narrow": false,
			"nsec3param": "",
			"rrsets": [{
				"comments": [],
				"name": "www.test.example.",
				"records": [{
					"content": "127.0.0.1",
					"disabled": false
				}],
				"ttl": 60,
				"type": "A"
			}, {
				"comments": [],
				"name": "test.example.",
				"records": [{
					"content": "a.misconfigured.powerdns.server. hostmaster.example3.de. 2019031801 10800 3600 604800 3600",
					"disabled": false
				}],
				"ttl": 3600,
				"type": "SOA"
			}, {
				"comments": [],
				"name": "test.example.",
				"records": [{
					"content": "ns1.example.com.",
					"disabled": false
				}, {
					"content": "ns2.example.com.", 
					"disabled": false
				}],
				"ttl": 3600,
				"type": "NS"
			}],
			"serial": 2019031801,
			"soa_edit": "",
			"soa_edit_api": "DEFAULT",
			"url": "/api/v1/servers/localhost/zones/example3.de."
		}`)

	hc := &http.Client{Transport: gock.DefaultTransport}
	c := pdnshttp.NewClient("http://dns.example", hc, &pdnshttp.APIKeyAuthenticator{APIKey: "secret"}, ioutil.Discard)
	sc := New(c)

	zone, err := sc.CreateZone(
		context.Background(),
		"localhost",
		Zone{
			Name: "test.example.",
			Kind: ZoneKindNative,
			Nameservers: []string{
				"ns1.example.com",
				"ns2.example.com",
			},
			ResourceRecordSets: []ResourceRecordSet{
				{
					Name: "www.test.example",
					Type: "A",
					TTL:  3600,
					Records: []Record{
						{Content: "127.0.0.1"},
					},
				},
			},
		},
	)

	assert.Nil(t, err)
	require.NotNil(t, zone)
	assert.Equal(t, "some-generated-id", zone.ID)
}
