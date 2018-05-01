package core

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

type KeyValueEntry xdr.KeyValueEntryValue

func (k *KeyValueEntry) Scan(src interface{}) error {
	var data string
	switch rawData := src.(type) {
	case []byte:
		data = string(rawData)
	case string:
		data = rawData

	default:
		return errors.New("Unexpected type")
	}
	err := xdr.SafeUnmarshalBase64(data,k);
	if err!=nil {
		return  errors.New("Faild to unmarshal key_value")
	}

	return nil
}

type KeyValue struct {
	Key		 string         `db:"key"`
	Value    KeyValueEntry 	`db:"value"`
}
