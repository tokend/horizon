package resource

import "gitlab.com/swarmfund/horizon/db2/core"

type KeyValue struct {
	Key 	string 					`json:"key"`
	Value 	map[string]interface{} 	`json:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue)  {
	k.Key = keyValue.Key
	k.Value,_ = keyValue.GetDetails()
}