package internal

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/rgenerated"
)

//MarshalCustomDetails - tries to marshal details to json.RawMessage
// if fails to marshal populates details of the error
func MarshalCustomDetails(details xdr.Longstring) rgenerated.Details {
	var result rgenerated.Details
	err := json.Unmarshal([]byte(details), &result)
	if err != nil {
		return marshalError(details, err)
	}

	return result
}

func marshalError(rawData xdr.Longstring, cause error) rgenerated.Details {
	data := map[string]interface{}{
		"_meta": map[string]interface{}{
			"raw_data": string(rawData),
			"error":    errors.Wrap(cause, "expected json object").Error(),
		},
	}

	result, err := json.Marshal(data)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal error for details"))
	}

	return result
}
