package search

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const exampleSearchResult = `[{"name": "example-search.de.", "object_type": "zone", "zone_id": "example-search.de."}, {"content": "127.0.0.1", "disabled": false, "name": "example-search.de.", "object_type": "record", "ttl": 60, "type": "A", "zone": "example-search.de.", "zone_id": "example-search.de."}, {"content": "ns1.example.com.", "disabled": false, "name": "example-search.de.", "object_type": "record", "ttl": 3600, "type": "NS", "zone": "example-search.de.", "zone_id": "example-search.de."}, {"content": "ns2.example.com.", "disabled": false, "name": "example-search.de.", "object_type": "record", "ttl": 3600, "type": "NS", "zone": "example-search.de.", "zone_id": "example-search.de."}, {"content": "a.misconfigured.powerdns.server. hostmaster.example-search.de. 2019031901 10800 3600 604800 3600", "disabled": false, "name": "example-search.de.", "object_type": "record", "ttl": 3600, "type": "SOA", "zone": "example-search.de.", "zone_id": "example-search.de."}]`

func TestResultListCanBeJSONDecoded(t *testing.T) {
	out := make(ResultList, 0)
	err := json.Unmarshal([]byte(exampleSearchResult), &out)

	assert.Nil(t, err)
	assert.Len(t, out, 5)
}

func TestResultListFilterByObjectTypeFiltersCorrectly(t *testing.T) {
	out := make(ResultList, 0)
	err := json.Unmarshal([]byte(exampleSearchResult), &out)

	require.Nil(t, err)

	filtered := out.FilterByObjectType(ObjectTypeZone)

	require.NotNil(t, filtered)
	assert.Len(t, filtered, 1)
	assert.Equal(t, ObjectTypeZone, filtered[0].ObjectType)
}

func TestResultListFilterByRecordTypeFiltersCorrectly(t *testing.T) {
	out := make(ResultList, 0)
	err := json.Unmarshal([]byte(exampleSearchResult), &out)

	require.Nil(t, err)

	filtered := out.FilterByRecordType("NS")

	require.NotNil(t, filtered)
	require.Len(t, filtered, 2)
	assert.Equal(t, ObjectTypeRecord, filtered[0].ObjectType)
	assert.Equal(t, ObjectTypeRecord, filtered[1].ObjectType)
}