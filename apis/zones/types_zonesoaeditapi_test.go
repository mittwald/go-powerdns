package zones

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneSOAEditAPISerializesCorrectly(t *testing.T) {
	data := []struct {
		v ZoneSOAEditAPI
		e string
	}{
		{ZoneSOAEditAPIDefault, `"DEFAULT"`},
		{ZoneSOAEditAPIIncrease, `"INCREASE"`},
		{ZoneSOAEditAPIEpoch, `"EPOCH"`},
		{ZoneSOAEditAPISoaEdit, `"SOA-EDIT"`},
		{ZoneSOAEditAPISoaEditIncrease, `"SOA-EDIT-INCREASE"`},
		{ZoneSOAEditAPINone, `"NONE"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			j, err := json.Marshal(data[i].v)

			assert.Nil(t, err)
			assert.Equal(t, data[i].e, string(j))
		})
	}
}

func TestZoneSOAEditAPISerializationReturnErrorOnUnknownValue(t *testing.T) {
	var v ZoneSOAEditAPI = 123

	_, err := json.Marshal(v)
	assert.NotNil(t, err)
}

func TestZoneSOAEditAPIUnserializesCorrectly(t *testing.T) {
	data := []struct {
		v ZoneSOAEditAPI
		e string
	}{
		{ZoneSOAEditAPIDefault, `"DEFAULT"`},
		{ZoneSOAEditAPIIncrease, `"INCREASE"`},
		{ZoneSOAEditAPIEpoch, `"EPOCH"`},
		{ZoneSOAEditAPISoaEdit, `"SOA-EDIT"`},
		{ZoneSOAEditAPISoaEditIncrease, `"SOA-EDIT-INCREASE"`},
		{ZoneSOAEditAPINone, `"NONE"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			var out ZoneSOAEditAPI

			err := json.Unmarshal([]byte(data[i].e), &out)

			assert.Nil(t, err)
			assert.Equal(t, data[i].v, out)
		})
	}
}

func TestZoneSOAEditAPIUnserializationReturnErrorOnUnknownValue(t *testing.T) {
	e := []byte(`"FOO"`)

	var out ZoneSOAEditAPI

	err := json.Unmarshal(e, &out)
	assert.NotNil(t, err)
}
