package regources

import "gitlab.com/tokend/go/xdr"

//ManageKeyValue - stores details of create account operation
type ManageKeyValue struct {
	Key
	Attributes ManageKeyValueAttrs `json:"attributes"`
}

//ManageKeyValueAttrs - details of ManageKeyValueOp
type ManageKeyValueAttrs struct {
	Key    string             `json:"key"`
	Action xdr.ManageKvAction `json:"action"`
	Value  *KeyValue          `json:"value,omitempty"`
}
