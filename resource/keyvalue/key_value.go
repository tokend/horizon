package keyvalue

import (
	"gitlab.com/swarmfund/horizon/db2/core"
)


type KeyValue struct {
	Key 	string 						`db:"key"`
	Value 	*core.KeyValueEntry 		`db:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key = keyValue.Key
	return k.Value.Scan(keyValue.Body)
}