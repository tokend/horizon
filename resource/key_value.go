package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
)


type KeyValue struct {
	Key 		string 					`json:"key"`
	Type        KeyValueType			`json:"type,omitempty"`
	Ui32Value   *uint32           		`json:"ui32_value,omitempty"`
	StringValue *string           		`json:"string_value,omitempty"`
}

type KeyValueType struct{
	Name 	string	`json:"name"`
	Value 	int		`json:"value"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key = keyValue.Key

	k.Type.Name = keyValue.Value.Type.ShortString()
	k.Type.Value = int(keyValue.Value.Type)

	k.Ui32Value = nil
	if keyValue.Value.Ui32Value != nil {
		uint32Value := uint32(*keyValue.Value.Ui32Value)
		k.Ui32Value = &uint32Value
	}
	k.StringValue = keyValue.Value.StringValue
	return nil
}