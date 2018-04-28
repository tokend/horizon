package keyvalue

import (
	"gitlab.com/swarmfund/horizon/db2/core"
)


type KeyValue struct {
	Key 	string 						`db:"key"`
	Type 	int32
	Value 	int32 						`db:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key = keyValue.Key
	var kvBody core.KeyValueEntry
	if err :=kvBody.Scan(keyValue.Body); err!=nil{
		return err
	}
	k.Value = int32(*kvBody.DefaultMask)
	k.Type = int32(kvBody.Type)
	return nil
}