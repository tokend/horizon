package keyvalue

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/go/xdr"
)


type KeyValue struct {
	Key 	string 						`json:"key"`
	Value 	*xdr.KeyValueEntryValue 	`json:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue)  {
	k.Key = keyValue.Key
	k.Value,_ = keyValue.Scan()
}