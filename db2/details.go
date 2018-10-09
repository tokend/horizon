package db2

import (
	"database/sql/driver"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Details map[string]interface{}

func (r Details) Value() (driver.Value, error) {
	result, err := DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	return result, nil
}

func (r *Details) Scan(src interface{}) error {
	err := DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details")
	}

	return nil
}

type DetailsWithPayload struct {
	Details   Details `json:"details"`
	Author    string  `json:"author"`
	CreatedAt int64   `json:"created_at"`
}

func (r DetailsWithPayload) Value() (driver.Value, error) {
	result, err := DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details array")
	}

	return result, nil
}

func (r *DetailsWithPayload) Scan(src interface{}) error {
	err := DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details array")
	}

	return nil
}
