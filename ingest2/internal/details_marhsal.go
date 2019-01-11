package internal

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

func MarshalCustomDetails(details xdr.Longstring) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(details), &result)
	if err != nil {
		return map[string]interface{}{
			"_meta": map[string]interface{}{
				"raw_data": string(details),
				"error":    errors.Wrap(err, "expected json object").Error(),
			},
		}
	}

	return result
}
