package resource

import (
	"time"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/httpx"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/render/hal"
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

	LedgersState         ledger.SystemState `json:"ledgers_state"`
	NetworkPassphrase    string             `json:"network_passphrase"`
	CommissionAccountID  string             `json:"commission_account_id"`
	OperationalAccountID string             `json:"operational_account_id"`
	StorageFeeAccountID  string             `json:"storage_fee_account_id"`
	MasterAccountID      string             `json:"master_account_id"`
	MasterExchangeName   string             `json:"master_exchange_name"`
	TxExpirationPeriod   int64              `json:"tx_expiration_period"`
	CurrentTime          int64              `json:"current_time"`
	Precision            int64              `json:"precision"`
	HorizonRevision      string             `json:"horizon_revision"`
}

//go:generate git rev-parse HEAD
// Populate fills in the details
func (res *Root) PopulateLedgerState(
	ctx context.Context,
	ledgerState ledger.SystemState,
) {
	res.LedgersState = ledgerState
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
