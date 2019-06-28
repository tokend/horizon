/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ExternalSystemData struct {
	Data ExternalSystemDataEntry `json:"data"`
	// Possible values: * address * address_with_payload
	Type ExternalDataType `json:"type"`
}

//Value - implements db driver method for auto marshal
func (r ExternalSystemData) Value() (driver.Value, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal ExternalDataType data")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *ExternalSystemData) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for jsonb")
	}

	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ExternalSystemData data")
	}

	return nil
}
