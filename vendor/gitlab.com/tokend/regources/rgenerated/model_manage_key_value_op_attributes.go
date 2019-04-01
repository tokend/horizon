package rgenerated

type ManageKeyValueOpAttributes struct {
	// * 1: \"put\", * 2: \"remove\"
	Action xdr.ManageKvAction `json:"action"`
	// Key of key-value entry to manage
	Key   string              `json:"key"`
	Value *KeyValueEntryValue `json:"value,omitempty"`
}
