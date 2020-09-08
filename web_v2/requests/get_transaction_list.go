package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeTransactionListLedgerEntryChanges - defines if the ledger entry changes should be included in the response
	IncludeTypeTransactionListLedgerEntryChanges = "ledger_entry_changes"

	// FilterTypeLedgerEntryType - defines if we need to filter the list by ledger entries transactions affected
	FilterTypeTransactionListLedgerEntryTypes = "ledger_entry_changes.entry_types"
	// FilterTypeLedgerEntryType - defines if we need to filter the list by ledger changes transactions produced
	FilterTypeTransactionListLedgerChangeTypes = "ledger_entry_changes.change_types"
)

// GetTransactions - represents params to be specified for GetTransactions handler
type GetTransactions struct {
	*base
	Filters struct {
		EntryTypes  []int `filter:"ledger_entry_changes.entry_types"`
		ChangeTypes []int `filter:"ledger_entry_changes.change_types"`
	}
	Includes struct {
		LedgerEntryChanges bool `include:"ledger_entry_changes"`
	}
	PageParams pgdb.CursorPageParams
}

// NewGetTransactions returns the new instance of GetTransactions request
func NewGetTransactions(r *http.Request) (*GetTransactions, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeTransactionListLedgerEntryChanges: {},
		},
		supportedFilters: map[string]struct{}{
			FilterTypeTransactionListLedgerEntryTypes:  {},
			FilterTypeTransactionListLedgerChangeTypes: {},
		},
	})
	if err != nil {
		return nil, err
	}

	request := GetTransactions{
		base: b,
	}

	err = urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = SetDefaultCursorPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
