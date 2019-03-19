package zones

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeTypeSerializesCorrectly(t *testing.T) {
	data := []struct {
		v RecordSetChangeType
		e string
	}{
		{ChangeTypeDelete, `"DELETE"`},
		{ChangeTypeReplace, `"REPLACE"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			j, err := json.Marshal(data[i].v)

			assert.Nil(t, err)
			assert.Equal(t, data[i].e, string(j))
		})
	}
}

func TestChangeTypeSerializationReturnErrorOnUnknownValue(t *testing.T) {
	var v RecordSetChangeType = 123

	_, err := json.Marshal(v)
	assert.NotNil(t, err)
}

func TestChangeTypeUnserializesCorrectly(t *testing.T) {
	data := []struct {
		v RecordSetChangeType
		e string
	}{
		{ChangeTypeDelete, `"DELETE"`},
		{ChangeTypeReplace, `"REPLACE"`},
	}

	for i := range data {
		t.Run(fmt.Sprintf("serializes to %s", data[i].e), func(t *testing.T) {
			var out RecordSetChangeType

			err := json.Unmarshal([]byte(data[i].e), &out)

			assert.Nil(t, err)
			assert.Equal(t, data[i].v, out)
		})
	}
}

func TestChangeTypeUnserializationReturnErrorOnUnknownValue(t *testing.T) {
	e := []byte(`"FOO"`)

	var out RecordSetChangeType

	err := json.Unmarshal(e, &out)
	assert.NotNil(t, err)
}
