package resource

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/valueflag"
)

type KeyValue struct {
	Key         string         `json:"key"`
	Type        valueflag.Flag `json:"type,omitempty"`
	Ui32Value   *uint32        `json:"ui32_value,omitempty"`
	StringValue *string        `json:"string_value,omitempty"`
}

func (k *KeyValue) Populate(keyValue *core.KeyValue) error {
	k.Key = keyValue.Key

	k.Type.Name = keyValue.Value.Type.ShortString()
	k.Type.Value = int32(keyValue.Value.Type)

	switch keyValue.Value.Type {
	case xdr.KeyValueEntryTypeUint32:
		k.Ui32Value = nil
		if keyValue.Value.Ui32Value != nil {
			uint32Value := uint32(*keyValue.Value.Ui32Value)
			k.Ui32Value = &uint32Value
		}
	case xdr.KeyValueEntryTypeString:
		k.StringValue = keyValue.Value.StringValue
	default:
		return errors.New("Unexpected key value type")
	}

	return nil
}
