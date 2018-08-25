package resource

import (
	"time"

	"gitlab.com/swarmfund/horizon/httpx"
	"gitlab.com/swarmfund/horizon/ledger"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/tokend/go/amount"
	"golang.org/x/net/context"
)

// Root is the initial map of links into the api.
type Root struct {
	Links struct {
		Account             hal.Link `json:"account"`
		AccountTransactions hal.Link `json:"account_transactions"`
		Metrics             hal.Link `json:"metrics"`
		Self                hal.Link `json:"self"`
		Transaction         hal.Link `json:"transaction"`
		Transactions        hal.Link `json:"transactions"`
	} `json:"_links"`

	HorizonSequence      int32  `json:"history_latest_ledger"`
	HistoryElderSequence int32  `json:"history_elder_ledger"`
	CoreSequence         int32  `json:"core_latest_ledger"`
	CoreElderSequence    int32  `json:"core_elder_ledger"`
	NetworkPassphrase    string `json:"network_passphrase"`
	CommissionAccountID  string `json:"commission_account_id"`
	OperationalAccountID string `json:"operational_account_id"`
	StorageFeeAccountID  string `json:"storage_fee_account_id"`
	MasterAccountID      string `json:"master_account_id"`
	MasterExchangeName   string `json:"master_exchange_name"`
	TxExpirationPeriod   int64  `json:"tx_expiration_period"`
	CurrentTime          int64  `json:"current_time"`
	Precision            int64  `json:"precision"`
}

// Populate fills in the details
func (res *Root) PopulateLedgerState(
	ctx context.Context,
	ledgerState ledger.State,
) {
	res.HorizonSequence = ledgerState.HistoryLatest
	res.HistoryElderSequence = ledgerState.HistoryElder
	res.CoreSequence = ledgerState.CoreLatest
	res.CoreElderSequence = ledgerState.CoreElder
	res.CurrentTime = time.Now().Unix()

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	res.Links.Account = lb.Link("/accounts/{account_id}")
	res.Links.AccountTransactions = lb.PagedLink("/accounts/{account_id}/transactions")
	res.Links.Metrics = lb.Link("/metrics")
	res.Links.Self = lb.Link("/")
	res.Links.Transaction = lb.Link("/transactions/{hash}")
	res.Links.Transactions = lb.PagedLink("/transactions")
	res.Precision = amount.One
}
