package pdnshttp

import (
	"net/http"
)

type ClientAuthenticator interface {
	OnRequest(*http.Request) error
	OnConnect(*http.Client) error
}
