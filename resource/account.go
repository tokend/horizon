package resource

import (
	"fmt"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/httpx"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource/base"
	"golang.org/x/net/context"
)

// Account is the summary of an account
type Account struct {
	Links struct {
		Self         hal.Link `json:"self"`
		Transactions hal.Link `json:"transactions"`
		Operations   hal.Link `json:"operations"`
		Payments     hal.Link `json:"payments"`
	} `json:"_links"`

	HistoryAccount
	IsBlocked     bool              `json:"is_blocked"`
	BlockReasonsI int32             `json:"block_reasons_i"`
	BlockReasons  []base.Flag       `json:"block_reasons"`
	AccountTypeI  int32             `json:"account_type_i"`
	AccountType   string            `json:"account_type"`
	Thresholds    AccountThresholds `json:"thresholds"`
	Balances      []Balance         `json:"balances"`
	Signers
	Limits                 `json:"limits"`
	Statistics             `json:"statistics"`
	Policies               AccountPolicies           `json:"policies"`
	ExternalSystemAccounts []ExternalSystemAccountID `json:"external_system_accounts"`
}

// Populate fills out the resource's fields
func (a *Account) Populate(ctx context.Context, ca core.Account) {
	a.ID = ca.AccountID
	a.AccountID = ca.AccountID
	a.BlockReasonsI = ca.BlockReasons
	a.BlockReasons = base.FlagFromXdrBlockReasons(ca.BlockReasons, xdr.BlockReasonsAll)
	a.IsBlocked = ca.BlockReasons > 0
	a.AccountTypeI = ca.AccountType
	a.AccountType = xdr.AccountType(ca.AccountType).String()
	a.Thresholds.Populate(ca.Thresholds)
	a.Policies.Populate(ca.Policies)
	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	self := fmt.Sprintf("/accounts/%s", ca.AccountID)
	a.Links.Self = lb.Link(self)
	a.Links.Transactions = lb.PagedLink(self, "transactions")
	a.Links.Operations = lb.PagedLink(self, "operations")
	a.Links.Payments = lb.PagedLink(self, "payments")
	a.Statistics.Populate(*ca.Statistics)
}

func (a *Account) SetBalances(balances []core.Balance) {
	a.Balances = make([]Balance, len(balances))
	for i := range balances {
		a.Balances[i].Populate(balances[i])
	}
}

func (a Account) PagingToken() string {
	return a.ID
}
