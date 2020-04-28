package db2

import (
	"database/sql/driver"
	"encoding/json"
	"gitlab.com/distributed_lab/kit/pgdb"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

// DEPRECATED: use json.RawMessage
type Details map[string]interface{}

func (r Details) Value() (driver.Value, error) {
	result, err := pgdb.JSONValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	return result, nil
}

func (r *Details) Scan(src interface{}) error {
	err := pgdb.JSONScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details")
	}

	return nil
}

func (r *Details) ToRawMessage() regources.Details {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return regources.Details(bytes)
}
