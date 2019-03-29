package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
	"net/http"
)

// GetTransactions - processes request to get the list of transactions (with ledger changes)
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getTransactionsHandler{
		LedgerQ:        *history2.NewLedgerQ(historyRepo),
		LedgerChangesQ: history2.NewLedgerChangesQ(historyRepo),
		TransactionsQ:  history2.NewTransactionsQ(historyRepo),
		Log:            ctx.Log(r),
	}

	request, err := requests.NewGetTransactions(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	result, err := handler.GetTransactions(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get transactions list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getTransactionsHandler struct {
	TransactionsQ  history2.TransactionsQ
	LedgerChangesQ history2.LedgerChangesQ
	LedgerQ        history2.LedgerQ
	Log            *logan.Entry
}

func (h *getTransactionsHandler) getLatestLedger() (*history2.Ledger, error) {
	sequence, err := h.LedgerQ.GetLatestLedgerSeq()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load latest ledger sequence")
	}

	ledger, err := h.LedgerQ.GetBySequence(sequence)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger by sequence")
	}

	return ledger, nil
}

// GetTransactions returns the list of transactions with related resources
func (h *getTransactionsHandler) GetTransactions(request *requests.GetTransactions) (*regources.TransactionsResponse, error) {
	q := h.TransactionsQ.Page(*request.PageParams)

	if request.ShouldFilter(requests.FilterTypeTransactionListLedgerChangeTypes) {
		q = q.FilterByEffectTypes(request.Filters.ChangeTypes...)
	}

	if request.ShouldFilter(requests.FilterTypeTransactionListLedgerEntryTypes) {
		q = q.FilterByLedgerEntryTypes(request.Filters.EntryTypes...)
	}

	historyTransactions, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load transactions")
	}

	result := regources.TransactionsResponse{
		Data: make([]regources.Transaction, 0, len(historyTransactions)),
	}

	for _, historyTransaction := range historyTransactions {
		historyChanges, err := h.LedgerChangesQ.FilterByTransactionID(historyTransaction.ID).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load ledger changes")
		}

		transaction := resources.NewTransaction(historyTransaction)
		transaction.Relationships.Source = resources.NewAccountKey(historyTransaction.Account).AsRelation()
		transaction.Relationships.LedgerEntryChanges = &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(historyChanges)),
		}

		operations := make(map[int64]regources.Key)

		for _, historyChange := range historyChanges {
			change, err := resources.NewLedgerEntryChange(historyChange)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse ledger entry change")
			}
			transaction.Relationships.LedgerEntryChanges.Data = append(
				transaction.Relationships.LedgerEntryChanges.Data, change.Key,
			)
			if request.ShouldInclude(requests.IncludeTypeTransactionListLedgerEntryChanges) {
				result.Included.Add(change)
			}
			operations[historyChange.OperationID] = resources.NewOperationKey(historyChange.OperationID)
		}

		transaction.Relationships.Operations = &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(operations)),
		}
		for _, operation := range operations {
			transaction.Relationships.Operations.Data = append(
				transaction.Relationships.Operations.Data, operation,
			)
		}

		result.Data = append(result.Data, transaction)
	}

	if len(result.Data) > 0 {
		result.Links = request.GetCursorLinks(*request.PageParams, result.Data[len(result.Data)-1].ID)
	} else {
		result.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	latestLedger, err := h.getLatestLedger()
	// TODO: possible race condition may occur if new ledger will be closed between querying transactions and ledgers
	// need to find a solution and fix somehow
	if err != nil {
		return nil, errors.Wrap(err, "failed to load latest ledger")
	}

	result.Meta = regources.TransactionResponseMeta{
		LatestLedgerCloseTime: latestLedger.ClosedAt,
		LatestLedgerSequence: latestLedger.Sequence,
	}

	return &result, nil
}
