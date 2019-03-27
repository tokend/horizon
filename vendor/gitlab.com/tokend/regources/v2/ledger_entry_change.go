package regources

import "gitlab.com/tokend/go/xdr"

type LedgerEntryChange struct {
	Key
	Attributes LedgerEntryChangeAttrs `json:"attributes"`
}

type LedgerEntryChangeAttrs struct {
	Payload    string                    `json:"payload"`
	ChangeType xdr.LedgerEntryChangeType `json:"change_type"`
	EntryType  xdr.LedgerEntryType       `json:"entry_type"`
}
