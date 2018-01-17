package db2

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Details map[string]interface{}

func (r Details) Value() (driver.Value, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, errors.New("failed to marshal sale details")
	}

	return data, nil
}

func (r *Details) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for details")
	}

	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal details")
	}

	return nil
}
