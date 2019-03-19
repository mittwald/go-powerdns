# PowerDNS client library for Go

[![GoDoc](https://godoc.org/github.com/mittwald/go-powerdns?status.svg)](https://godoc.org/github.com/mittwald/go-powerdns)

This package contains a Go library for accessing the [PowerDNS][powerdns] API.

## Installation

Install using `go get`:

```console
> go get github.com/mittwald/go-powerdns
```

## Usage

First, instantiate a client using `pdns.New`:

```go
client, err := pdns.New(
    pdns.WithBaseURL("http://localhost:8081"),
    pdns.WithAPIKeyAuthentication("supersecret"),
)
```

The client then offers more specialiced sub-clients, for example for managing server and zones.
Have a look at this library's [documentation][godoc] for more information.

## Complete example

```go
package main

import "context"
import "github.com/mittwald/go-powerdns"
import "github.com/mittwald/go-powerdns/apis/zones"

func main() {
    client, err := pdns.New(
        pdns.WithBaseURL("http://localhost:8081"),
        pdns.WithAPIKeyAuthentication("supersecret"),
    )
	
    if err != nil {
    	panic(err)
    }
    
    client.Zones().CreateZone(context.Background(), "localhost", zones.Zone{
        Name: "example.com.",
    })
}
```

[powerdns]: https://github.com/PowerDNS/pdns
[godoc]: https://godoc.org/github.com/mittwald/go-powerdns