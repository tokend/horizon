package regources

import "gitlab.com/tokend/go/xdr"

//KeyValue - represents xdr.KeyValueEntryValue
type KeyValue struct {
	Type xdr.KeyValueEntryType `json:"type"`
	U32  *uint32               `json:"u_32,omitempty"`
	Str  *string               `json:"str,omitempty"`
	U64  *uint64               `json:"u_64,omitempty"`
}
