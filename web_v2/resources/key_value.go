package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

func NewKeyValue(record core2.KeyValue) rgenerated.KeyValueEntry {
	return rgenerated.KeyValueEntry{
		Key: rgenerated.Key{
			Type: rgenerated.KEY_VALUE_ENTRIES,
			ID:   record.Key,
		},
		Attributes: rgenerated.KeyValueEntryAttributes{
			Value: rgenerated.KeyValueEntryValue{
				Type: record.Value.Type,
				U32:  (*uint32)(record.Value.Ui32Value),
				U64:  (*uint64)(record.Value.Ui64Value),
				Str:  (*string)(record.Value.StringValue),
			},
		},
	}
}
