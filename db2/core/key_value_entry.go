package core

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

type KeyValueEntry xdr.KeyValueEntryValue

func (k *KeyValueEntry) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type")
	}
	err := xdr.SafeUnmarshal(data,k);
	if err!=nil {
		return  errors.New("Faild to unmarshal key_value")
	}

	return nil
}

type KeyValue struct {
	Key		string                `db:"key"`
	Body    []byte                `db:"value"`
}
