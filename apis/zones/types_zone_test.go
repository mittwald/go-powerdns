package zones

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNilNameserversAreSerializedToEmptyArray(t *testing.T) {
	z := Zone{
		ID:                 "1",
		Name:               "foo.example",
		Type:               ZoneTypeZone,
		Kind:               ZoneKindMaster,
		Nameservers:        nil,
	}

	j, err := json.Marshal(z)

	require.Nil(t, err)
	require.Equal(t, `{"id":"1","name":"foo.example","type":"Zone","kind":"Master","nameservers":[]}`, string(j))
}

func TestNameserversAreSerializedToArray(t *testing.T) {
	z := Zone{
		ID:                 "1",
		Name:               "foo.example",
		Type:               ZoneTypeZone,
		Kind:               ZoneKindMaster,
		Nameservers:        ZoneNameservers{"ns.foo.example"},
	}

	j, err := json.Marshal(z)

	require.Nil(t, err)
	require.Equal(t, `{"id":"1","name":"foo.example","type":"Zone","kind":"Master","nameservers":["ns.foo.example"]}`, string(j))
}