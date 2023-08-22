package requests

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
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
	// FilterTypeTransactionListBeforeTimestamp - defines if we need to filter the list before specified ledger close time
	FilterTypeTransactionListBeforeTimestamp = "before"
	// FilterTypeTransactionListAfterTimestamp - defines if we need to filter the list after specified ledger close time
	FilterTypeTransactionListAfterTimestamp = "after"
)

// GetTransactions - represents params to be specified for GetTransactions handler
type GetTransactions struct {
	*base
	Filters struct {
		EntryTypes      []int  `filter:"ledger_entry_changes.entry_types"`
		ChangeTypes     []int  `filter:"ledger_entry_changes.change_types"`
		BeforeTimestamp *int64 `filter:"before"`
		AfterTimestamp  *int64 `filter:"after"`
	}
	Includes struct {
		LedgerEntryChanges bool `include:"ledger_entry_changes"`
	}
	PageParams pgdb.CursorPageParams
	PageNumber *uint64 `page:"number"`
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
			FilterTypeTransactionListBeforeTimestamp:   {},
			FilterTypeTransactionListAfterTimestamp:    {},
		},
	})
	if err != nil {
		return nil, err
	}

	request := GetTransactions{
		base: b,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = SetDefaultCursorPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	// use part of cursor params struct to prevent decode same token twice
	if request.PageNumber != nil {
		params := pgdb.OffsetPageParams{
			Limit:      request.PageParams.Limit,
			Order:      request.PageParams.Order,
			PageNumber: *request.PageNumber,
		}

		err = request.SetDefaultOffsetPageParams(&params)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set default page params")
		}

		request.PageParams.Limit = params.Limit
		request.PageParams.Order = params.Order
		request.PageNumber = &params.PageNumber
	}

	return &request, nil
}
