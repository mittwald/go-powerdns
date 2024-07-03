package search

import (
	"context"
	"fmt"
	"github.com/mittwald/go-powerdns/pdnshttp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io"
	"net/http"
	"testing"
)

func TestSearchExecutesCorrectRequest(t *testing.T) {
	cases := []struct {
		query              string
		max                int
		objectType         ObjectType
		expectedObjectType string
	}{
		{"example.com", 10, ObjectTypeAll, "all"},
		{"example.com", 10, ObjectTypeZone, "zone"},
		{"example.com", 10, ObjectTypeRecord, "record"},
		{"example.com", 10, ObjectTypeComment, "comment"},
		{"example.com", 15, ObjectTypeComment, "comment"},
	}

	for i := range cases {
		t.Run(fmt.Sprintf("test with %+v", cases[i]), func(t *testing.T) {
			gock.New("http://dns.example").
				Get("/api/v1/servers/localhost/search-data").
				MatchParam("object_type", cases[i].expectedObjectType).
				MatchParam("max", fmt.Sprintf("%d", cases[i].max)).
				MatchParam("q", cases[i].query).
				Reply(http.StatusOK).
				SetHeader("Content-Type", "application/json").
				BodyString(exampleSearchResult)

			hc := &http.Client{Transport: gock.DefaultTransport}
			c := pdnshttp.NewClient("http://dns.example", hc, &pdnshttp.APIKeyAuthenticator{APIKey: "secret"}, io.Discard)
			sc := New(c)

			results, err := sc.Search(
				context.Background(),
				"localhost",
				cases[i].query,
				cases[i].max,
				cases[i].objectType,
			)

			assert.Nil(t, err)
			assert.IsType(t, ResultList{}, results)
			assert.Len(t, results, 5)

			assert.True(t, gock.IsDone())
		})
	}

}
