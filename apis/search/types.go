package search

import "fmt"

type ObjectType int

const (
	_             = iota
	ObjectTypeAll ObjectType = iota
	ObjectTypeZone
	ObjectTypeRecord
	ObjectTypeComment
)

func (t ObjectType) String() string {
	switch t {
	case ObjectTypeAll:
		return "all"
	case ObjectTypeZone:
		return "zone"
	case ObjectTypeRecord:
		return "record"
	case ObjectTypeComment:
		return "comment"
	}

	return ""
}

func (t *ObjectType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"all"`:
		*t = ObjectTypeAll
	case `"zone"`:
		*t = ObjectTypeZone
	case `"record"`:
		*t = ObjectTypeRecord
	case `"comment"`:
		*t = ObjectTypeComment
	default:
		return fmt.Errorf(`unknown search type: %s'`, string(b))
	}

	return nil
}

type SearchResult struct {
	Content    string     `json:"content"`
	Disabled   bool       `json:"disabled"`
	Name       string     `json:"name"`
	ObjectType ObjectType `json:"object_type"`
	ZoneID     string     `json:"zone_id"`
	Zone       string     `json:"zone"`
	Type       string     `json:"type"`
	TTL        int        `json:"ttl"`
}
