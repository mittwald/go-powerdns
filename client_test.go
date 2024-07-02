package pdns

import (
	"context"
	"flag"
	"fmt"
	"github.com/mittwald/go-powerdns/apis/search"
	"github.com/mittwald/go-powerdns/apis/zones"
	"github.com/mittwald/go-powerdns/pdnshttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		fmt.Println("skipping integration tests")
		os.Exit(0)
	}

	runOrPanic("docker-compose", "rm", "-sfv")
	runOrPanic("docker-compose", "down", "-v")
	runOrPanic("docker-compose", "up", "-d")

	defer func() {
		runOrPanic("docker-compose", "down", "-v")
	}()

	c, err := New(
		WithBaseURL("http://localhost:8081"),
		WithAPIKeyAuthentication("secret"),
	)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c.WaitUntilUp(ctx)

	e := m.Run()

	if e != 0 {
		fmt.Println("")
		fmt.Println("TESTS FAILED")
		fmt.Println("Leaving containers running for further inspection")
		fmt.Println("")
	} else {
		runOrPanic("docker-compose", "down", "-v")
	}

	os.Exit(e)
}

func runOrPanic(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		panic(err)
	}
}

func TestCanConnect(t *testing.T) {
	c := buildClient(t)

	statusErr := c.Status()
	assert.Nil(t, statusErr)
}

func TestListServers(t *testing.T) {
	c := buildClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	servers, err := c.Servers().ListServers(ctx)

	assert.Nil(t, err, "ListServers returned error")
	assert.Lenf(t, servers, 1, "ListServers should return one server")
}

func TestGetServer(t *testing.T) {
	c := buildClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server, err := c.Servers().GetServer(ctx, "localhost")

	require.Nil(t, err, "GetServer returned error")
	require.NotNil(t, server)
	require.Equal(t, "authoritative", server.DaemonType)
}

func TestGetEmptyZones(t *testing.T) {
	c := buildClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	z, err := c.Zones().ListZones(ctx, "localhost")

	require.Nil(t, err, "ListZones returned error")

	assert.Len(t, z, 0)
}

func TestCreateZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{
				Name: "example.de.",
				Type: "A",
				TTL:  60,
				Records: []zones.Record{
					{Content: "127.0.0.1"},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.Nil(t, err, "CreateZone returned error")

	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "example.de.", created.Name)
}

func TestCreateZoneProducedReadableErrorMessages(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name:        "test-error-message.de.",
		Type:        zones.ZoneTypeZone,
		Kind:        zones.ZoneKindNative,
		Nameservers: []string{"ns1.example.com.", "ns2.example.com."},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Zones().CreateZone(ctx, "localhost", zone)
	require.Nil(t, err, "CreateZone returned error")

	_, err2 := c.Zones().CreateZone(ctx, "localhost", zone)
	require.Error(t, err2, "CreateZone should return error")
	require.Equal(t, "unexpected status code 409: http://localhost:8081/api/v1/servers/localhost/zones Conflict", err2.Error())
}

func TestDeleteZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example-delete.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{
				Name: "example-delete.de.",
				Type: "A",
				TTL:  60,
				Records: []zones.Record{
					{Content: "127.0.0.1"},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.Nil(t, err, "CreateZone returned error")

	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "example-delete.de.", created.Name)

	deleteErr := c.Zones().DeleteZone(ctx, "localhost", created.ID)
	require.Nil(t, deleteErr, "DeleteZone returned error")

	_, getErr := c.Zones().GetZone(ctx, "localhost", created.ID)
	assert.NotNil(t, getErr)
	assert.IsType(t, pdnshttp.ErrNotFound{}, getErr)
	assert.True(t, pdnshttp.IsNotFound(getErr))
}

func TestAddRecordToZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example2.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "foo.example2.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.Nil(t, err, "CreateZone returned error")

	err = c.Zones().AddRecordSetToZone(ctx, "localhost", created.ID, zones.ResourceRecordSet{
		Name:    "bar.example2.de.",
		Type:    "A",
		TTL:     60,
		Records: []zones.Record{{Content: "127.0.0.2"}},
	})

	require.Nil(t, err, "AddRecordSetToZone returned error")

	updated, err := c.Zones().GetZone(ctx, "localhost", created.ID)

	require.Nil(t, err)

	rs := updated.GetRecordSet("bar.example2.de.", "A")
	require.NotNil(t, rs)
}

func TestSelectZoneWithoutRRSets(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example5.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "foo.example5.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.NoError(t, err, "CreateZone returned error")

	zoneWithoutRRSets, err := c.Zones().GetZone(ctx, "localhost", created.ID, zones.WithoutResourceRecordSets())
	require.NoError(t, err, "GetZone returned error")
	require.Len(t, zoneWithoutRRSets.ResourceRecordSets, 0, "ResourceRecordSets should be empty")
}

func TestSelectFilteredRRSetsFromZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example4.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "foo.example4.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
			{Name: "bar.example4.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "10.0.0.1"}}},
			{Name: "bar.example4.de.", Type: "TXT", TTL: 60, Records: []zones.Record{{Content: `"Hello!"`}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.NoError(t, err, "CreateZone returned error")

	zoneWithRRSets, err := c.Zones().GetZone(ctx, "localhost", created.ID, zones.WithResourceRecordSetFilter("bar.example4.de.", "TXT"))

	require.NoError(t, err)
	require.Len(t, zoneWithRRSets.ResourceRecordSets, 1)
	require.Equal(t, "bar.example4.de.", zoneWithRRSets.ResourceRecordSets[0].Name)
	require.Equal(t, "TXT", zoneWithRRSets.ResourceRecordSets[0].Type)
	require.Equal(t, `"Hello!"`, zoneWithRRSets.ResourceRecordSets[0].Records[0].Content)
}

func TestRemoveRecordFromZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example3.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "foo.example3.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.Nil(t, err, "CreateZone returned error")

	err = c.Zones().AddRecordSetToZone(ctx, "localhost", created.ID, zones.ResourceRecordSet{
		Name:    "bar.example3.de.",
		Type:    "A",
		TTL:     60,
		Records: []zones.Record{{Content: "127.0.0.2"}},
	})

	require.Nil(t, err, "AddRecordSetToZone returned error")

	updated, err := c.Zones().GetZone(ctx, "localhost", created.ID)
	require.Nil(t, err)
	rs := updated.GetRecordSet("bar.example3.de.", "A")
	require.NotNil(t, rs)

	err = c.Zones().RemoveRecordSetFromZone(ctx, "localhost", created.ID, "bar.example3.de.", "A")
	require.Nil(t, err, "RemoveRecordSetFromZone returned error")

	updated, err = c.Zones().GetZone(ctx, "localhost", created.ID)
	require.Nil(t, err)
	rs = updated.GetRecordSet("bar.example3.de.", "A")
	require.Nil(t, rs)
}

func TestSearchZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example-search.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "example-search.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.Nil(t, err, "CreateZone returned error")

	results, sErr := c.Search().Search(ctx, "localhost", "example-search.de", 10, search.ObjectTypeZone)

	require.Nil(t, sErr)
	require.True(t, len(results) > 0, "number of search results should be > 0")

	assert.Equal(t, "example-search.de.", results[0].Name)
	assert.Equal(t, search.ObjectTypeZone, results[0].ObjectType)
}

func TestExportZone(t *testing.T) {
	c := buildClient(t)

	zone := zones.Zone{
		Name: "example-export.de.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "example-export.de.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	created, err := c.Zones().CreateZone(ctx, "localhost", zone)

	require.Nil(t, err, "CreateZone returned error")

	export, sErr := c.Zones().ExportZone(ctx, "localhost", created.ID)

	date := time.Now().Format("20060102") + "01"

	require.Nil(t, sErr)
	require.Equal(t, "example-export.de.\t60\tIN\tA\t127.0.0.1\nexample-export.de.\t3600\tIN\tNS\tns1.example.com.\nexample-export.de.\t3600\tIN\tNS\tns2.example.com.\nexample-export.de.\t3600\tIN\tSOA\ta.misconfigured.dns.server.invalid. hostmaster.example-export.de. "+date+" 10800 3600 604800 3600\n", string(export))
}

func buildClient(t *testing.T) Client {
	debug := io.Discard

	if testing.Verbose() {
		debug = os.Stderr
	}

	c, err := New(
		WithBaseURL("http://localhost:8081"),
		WithAPIKeyAuthentication("secret"),
		WithDebuggingOutput(debug),
	)

	assert.Nil(t, err)
	return c
}

// This example uses the "context.WithTimeout" function to wait until the PowerDNS API is reachable
// up until a given timeout is reached. After that, the "WaitUntilUp" method will return with an error.
func ExampleClient_waitUntilUp() {
	client, _ := New(
		WithBaseURL("http://localhost:8081"),
		WithAPIKeyAuthentication("secret"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.WaitUntilUp(ctx)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_listServers() {
	client, _ := New(
		WithBaseURL("http://localhost:8081"),
		WithAPIKeyAuthentication("secret"),
	)

	servers, err := client.Servers().ListServers(context.Background())
	if err != nil {
		panic(err)
	}
	for i := range servers {
		fmt.Printf("found server: %s\n", servers[i].ID)
	}
}

func ExampleClient_getServer() {
	client, _ := New(
		WithBaseURL("http://localhost:8081"),
		WithAPIKeyAuthentication("secret"),
	)

	server, err := client.Servers().GetServer(context.Background(), "localhost")
	if err != nil {
		if pdnshttp.IsNotFound(err) {
			// handle not found
		} else {
			panic(err)
		}
	}

	fmt.Printf("found server: %s\n", server.ID)
}

// This example uses the "Zones()" sub-client to create a new zone.
func ExampleClient_createZone() {
	client, _ := New(
		WithBaseURL("http://localhost:8081"),
		WithAPIKeyAuthentication("secret"),
	)

	input := zones.Zone{
		Name: "mydomain.example.",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: []zones.ResourceRecordSet{
			{Name: "foo.mydomain.example.", Type: "A", TTL: 60, Records: []zones.Record{{Content: "127.0.0.1"}}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	zone, err := client.Zones().CreateZone(ctx, "localhost", input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("zone ID: %s\n", zone.ID)
	// Output: zone ID: mydomain.example.
}
