package core

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

type KeyValue struct {
	Key		string                `db:"key"`
	Body    []byte                `db:"value"`
}

func (a *KeyValue) Scan() (*xdr.KeyValueEntryValue, error) {
	var result xdr.KeyValueEntryValue
	err := xdr.SafeUnmarshal(a.Body,result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal keyValue")
	}

	return &result, nil
}

func (a *KeyValue) Value()([]byte,error) {
	return a.Body,nil
}