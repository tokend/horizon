package regources

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//Details - custom type for external details
type Details json.RawMessage

//UnmarshalJSON - casts data to Details
func (d *Details) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("regources.Details: UnmarshalJSON on nil pointer")
	}
	*d = append((*d)[0:0], data...)
	return nil
}

//MarshalJSON - casts Details to []byte
func (d Details) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	return d, nil
}

func (d Details) String() string {
	return string(d)
}
