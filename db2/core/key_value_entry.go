package core

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"encoding/json"
)

type KeyValue struct {
	Key          string                `db:"key"`
	Value        []byte                `db:"value"`
}

func (a KeyValue) GetDetails() (db2.Details, error) {
	var result db2.Details
	err := json.Unmarshal(a.Value, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal asset details")
	}

	return result, nil
}