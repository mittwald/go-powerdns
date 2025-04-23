package zones

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneSOAEditSerializesCorrectly(t *testing.T) {
	data := []struct {
		v ZoneSOAEdit
		e string
	}{
		{ZoneSOAEditIncrementWeeks, `"INCREMENT-WEEKS"`},
		{ZoneSOAEditInceptionEpoch, `"INCEPTION-EPOCH"`},
		{ZoneSOAEditInceptionIncrement, `"INCEPTION-INCREMENT"`},
		{ZoneSOAEditEpoch, `"EPOCH"`},
		{ZoneSOAEditNone, `"NONE"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			j, err := json.Marshal(data[i].v)

			assert.Nil(t, err)
			assert.Equal(t, data[i].e, string(j))
		})
	}
}

func TestZoneSOAEditSerializationReturnErrorOnUnknownValue(t *testing.T) {
	var v ZoneSOAEdit = 123

	_, err := json.Marshal(v)
	assert.NotNil(t, err)
}

func TestZoneSOAEditUnserializesCorrectly(t *testing.T) {
	data := []struct {
		v ZoneSOAEdit
		e string
	}{
		{ZoneSOAEditIncrementWeeks, `"INCREMENT-WEEKS"`},
		{ZoneSOAEditInceptionEpoch, `"INCEPTION-EPOCH"`},
		{ZoneSOAEditInceptionIncrement, `"INCEPTION-INCREMENT"`},
		{ZoneSOAEditEpoch, `"EPOCH"`},
		{ZoneSOAEditNone, `"NONE"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			var out ZoneSOAEdit

			err := json.Unmarshal([]byte(data[i].e), &out)

			assert.Nil(t, err)
			assert.Equal(t, data[i].v, out)
		})
	}
}

func TestZoneSOAEditUnserializationReturnErrorOnUnknownValue(t *testing.T) {
	e := []byte(`"FOO"`)

	var out ZoneSOAEdit

	err := json.Unmarshal(e, &out)
	assert.NotNil(t, err)
}
