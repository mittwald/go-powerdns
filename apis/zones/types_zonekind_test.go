package zones

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)
func TestZoneKindSerializesCorrectly(t *testing.T) {
	data := []struct {
		v ZoneKind
		e string
	}{
		{ZoneKindNative, `"Native"`},
		{ZoneKindMaster, `"Master"`},
		{ZoneKindSlave, `"Slave"`},
		{ZoneKindProducer, `"Producer"`},
		{ZoneKindConsumer, `"Consumer"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			j, err := json.Marshal(data[i].v)

			assert.Nil(t, err)
			assert.Equal(t, data[i].e, string(j))
		})
	}
}

func TestZoneTypeSerializationReturnErrorOnUnknownValue(t *testing.T) {
	var v ZoneKind = 123

	_, err := json.Marshal(v)
	assert.NotNil(t, err)
}

func TestZoneKindUnserializesCorrectly(t *testing.T) {
	data := []struct {
		v ZoneKind
		e string
	}{
		{ZoneKindNative, `"Native"`},
		{ZoneKindMaster, `"Master"`},
		{ZoneKindSlave, `"Slave"`},
		{ZoneKindProducer, `"Producer"`},
		{ZoneKindConsumer, `"Consumer"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			var out ZoneKind

			err := json.Unmarshal([]byte(data[i].e), &out)

			assert.Nil(t, err)
			assert.Equal(t, data[i].v, out)
		})
	}
}

func TestZoneKindUnserializationReturnErrorOnUnknownValue(t *testing.T) {
	e := []byte(`"FOO"`)

	var out ZoneKind

	err := json.Unmarshal(e, &out)
	assert.NotNil(t, err)
}
