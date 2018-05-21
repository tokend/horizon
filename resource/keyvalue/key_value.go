package keyvalue

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
)


type KeyValue struct {
	Key 	string 						`json:"key"`
	Value 	xdr.KeyValueEntryValue 			`json:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key = keyValue.Key
	k.Value.Type = keyValue.Value.Type
	switch k.Value.Type {
	case xdr.KeyValueEntryTypeUint32:
		k.Value.Ui32Value = keyValue.Value.Ui32Value
	default:
		return errors.New("Unexpected type of KeyValueEntryValue")
	}

	return nil
}