package keyvalue

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)


type KeyValue struct {
	Key 	string 						`json:"key"`
	Value 	xdr.KeyValueEntryValue 			`json:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key = keyValue.Key
	k.Value = xdr.KeyValueEntryValue(keyValue.Value)
	return nil
}