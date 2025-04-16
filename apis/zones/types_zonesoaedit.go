package zones

import "fmt"

type ZoneSOAEdit int

const (
	ZoneSOAEditUnset          ZoneSOAEdit = iota
	ZoneSOAEditIncrementWeeks             = iota
	ZoneSOAEditInceptionEpoch
	ZoneSOAEditInceptionIncrement
	ZoneSOAEditEpoch
	ZoneSOAEditNone
)

func (v ZoneSOAEdit) MarshalJSON() ([]byte, error) {
	switch v {
	case ZoneSOAEditIncrementWeeks:
		return []byte(`"INCREMENT-WEEKS"`), nil
	case ZoneSOAEditInceptionEpoch:
		return []byte(`"INCEPTION-EPOCH"`), nil
	case ZoneSOAEditInceptionIncrement:
		return []byte(`"INCEPTION-INCREMENT"`), nil
	case ZoneSOAEditEpoch:
		return []byte(`"EPOCH"`), nil
	case ZoneSOAEditNone:
		return []byte(`"NONE"`), nil
	default:
		return nil, fmt.Errorf("unsupported SOA-EDIT value: %d", v)
	}
}

func (v *ZoneSOAEdit) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"INCREMENT-WEEKS"`:
		*v = ZoneSOAEditIncrementWeeks
	case `"INCEPTION-EPOCH"`:
		*v = ZoneSOAEditInceptionEpoch
	case `"INCEPTION-INCREMENT"`:
		*v = ZoneSOAEditInceptionIncrement
	case `"EPOCH"`:
		*v = ZoneSOAEditEpoch
	case `"NONE"`:
		*v = ZoneSOAEditNone
	case `""`:
		*v = ZoneSOAEditUnset
	default:
		return fmt.Errorf("unsupported SOA-EDIT value: %s", string(input))
	}
	return nil
}
