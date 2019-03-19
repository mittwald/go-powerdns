package pdns

import (
	"os"
)

// This example uses with WithAPIKeyAuthentication function to add API-key based authentication
// to the PowerDNS client.
func ExampleNew_withAPIKeyAuthentication() {
	client, err := New(
		WithBaseURL("http://your-dns-server.example:8081"),
		WithAPIKeyAuthentication("super-secret"),
	)

	if err != nil {
		panic(err)
	}

	client.Status()
}

// This example uses the WithDebuggingOutput function; this will cause all HTTP requests
// and responses to be logged to the io.Writer that is supplied to WithDebuggingOutput.
func ExampleNew_withDebugging() {
	client, err := New(
		WithBaseURL("http://your-dns-server.example:8081"),
		WithAPIKeyAuthentication("super-secret"),
		WithDebuggingOutput(os.Stdout),
	)

	if err != nil {
		panic(err)
	}

	client.Status()
}
