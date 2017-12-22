package history

import (
	"database/sql/driver"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SaleDetails struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Logo             string `json:"logo"`
}

func (r SaleDetails) Value() (driver.Value, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, errors.New("failed to marshal sale details")
	}

	return data, nil
}

func (r *SaleDetails) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for sale details")
	}

	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal sale details")
	}

	return nil
}
