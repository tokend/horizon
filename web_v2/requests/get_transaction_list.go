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
	PageParams *pgdb.CursorPageParams
	Filters    struct {
		EntryTypes  []int `filter:"ledger_entry_changes.entry_types"`
		ChangeTypes []int `filter:"ledger_entry_changes.change_types"`
	}
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

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetTransactions{
		base:       b,
		PageParams: pageParams,
	}

	err = urlval.Decode(r.URL.Query(), &request.Filters)

	return &request, nil
}

//func (r *GetTransactions) getIntSlice(name string) ([]int, error) {
//	valuesStr := strings.Split(r.getString(name), ",")
//
//	if len(valuesStr) > 0 {
//		valuesInt := make([]int, 0, len(valuesStr))
//		for _, v := range valuesStr {
//			if v != "" {
//				valueInt, err := strconv.Atoi(v)
//				if err != nil {
//					return nil, validation.Errors{
//						v: err,
//					}
//				}
//
//				valuesInt = append(valuesInt, valueInt)
//			}
//		}
//
//		return valuesInt, nil
//	}
//
//	return []int{}, nil
//}
//
//func (r *GetTransactions) populateFilters() (err error) {
//	r.Filters.EntryTypes, err = r.getIntSlice(
//		fmt.Sprintf("filter[%s]", FilterTypeTransactionListLedgerEntryTypes),
//	)
//	if err != nil {
//		return err
//	}
//
//	r.Filters.ChangeTypes, err = r.getIntSlice(
//		fmt.Sprintf("filter[%s]", FilterTypeTransactionListLedgerChangeTypes),
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
