package db2

import (
	"database/sql/driver"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func DriverValue(data interface{}) (driver.Value, error) {
	data, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("failed to marshal details")
	}

	return data, nil
}

func DriveScan(src, dest interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for jsonb")
	}

	err := json.Unmarshal(data, dest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal jsonb")
	}

	return nil
}
