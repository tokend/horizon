package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

func NewKeyValue(record core2.KeyValue) regources.KeyValueEntry {
	return regources.KeyValueEntry{
		Key: regources.Key{
			Type: regources.TypeKeyValueEntries,
			ID:   record.Key,
		},
		Attributes: regources.KeyValueEntryAttrs{
			Value: regources.KeyValueEntryValue{
				Type: record.Value.Type,
				U32:  (*uint32)(record.Value.Ui32Value),
				U64:  (*uint64)(record.Value.Ui64Value),
				Str:  (*string)(record.Value.StringValue),
			},
		},
	}
}
