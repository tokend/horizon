package keyvalue

import (
	"gitlab.com/swarmfund/horizon/db2/core"
)


type KeyValue struct {
	Key 	string 						`json:"key"`
	Type 	int32						`json:"type"`
	Value 	int32 						`json:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key 	= keyValue.Key
	k.Type 	= int32(keyValue.Value.Type)
	k.Value = int32(*keyValue.Value.DefaultMask)
	return nil
}