package zones

import "fmt"

type ZoneSOAEditAPI int

const (
	_                                    = iota
	ZoneSOAEditAPIDefault ZoneSOAEditAPI = iota
	ZoneSOAEditAPIIncrease
	ZoneSOAEditAPIEpoch
	ZoneSOAEditAPISoaEdit
	ZoneSOAEditAPISoaEditIncrease
	ZoneSOAEditAPINone
)

func (v ZoneSOAEditAPI) MarshalJSON() ([]byte, error) {
	switch v {
	case ZoneSOAEditAPIDefault:
		return []byte(`"DEFAULT"`), nil
	case ZoneSOAEditAPIIncrease:
		return []byte(`"INCREASE"`), nil
	case ZoneSOAEditAPIEpoch:
		return []byte(`"EPOCH"`), nil
	case ZoneSOAEditAPISoaEdit:
		return []byte(`"SOA-EDIT"`), nil
	case ZoneSOAEditAPISoaEditIncrease:
		return []byte(`"SOA-EDIT-INCREASE"`), nil
	case ZoneSOAEditAPINone:
		return []byte(`"NONE"`), nil
	default:
		return nil, fmt.Errorf("unsupported SOA-EDIT-API value: %d", v)
	}
}

func (v *ZoneSOAEditAPI) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"DEFAULT"`:
		*v = ZoneSOAEditAPIDefault
	case `"INCREASE"`:
		*v = ZoneSOAEditAPIIncrease
	case `"EPOCH"`:
		*v = ZoneSOAEditAPIEpoch
	case `"SOA-EDIT"`:
		*v = ZoneSOAEditAPISoaEdit
	case `"SOA-EDIT-INCREASE"`:
		*v = ZoneSOAEditAPISoaEditIncrease
	case `"NONE"`:
		*v = ZoneSOAEditAPINone
	default:
		return fmt.Errorf("unsupported SOA-EDIT-API value: %s", string(input))
	}
	return nil
}
