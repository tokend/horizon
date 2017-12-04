package resource

import (
	"fmt"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/httpx"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource/base"
	"golang.org/x/net/context"
)

// Populate fills out the resource's fields
func (this *Account) Populate(
	ctx context.Context,
	ca core.Account,
	cs []core.Signer,
	cb []core.Balance,
	cl *core.Limits,
) (err error) {
	this.ID = ca.AccountID
	this.AccountID = ca.AccountID
	this.BlockReasonsI = ca.BlockReasons
	this.BlockReasons = base.FlagFromXdrBlockReasons(ca.BlockReasons, xdr.BlockReasonsAll)
	this.IsBlocked = ca.BlockReasons > 0
	this.AccountTypeI = ca.AccountType
	this.AccountType = xdr.AccountType(ca.AccountType).String()

	this.Referrer = ca.Referrer
	this.ShareForReferrer = amount.String(int64(ca.ShareForReferrer))
	this.Thresholds.Populate(ca)

	this.Balances = make([]Balance, len(cb))
	for index, balance := range cb {
		err := this.Balances[index].Populate(balance)
		if err != nil {
			return err
		}
	}

	// populate signers
	this.Signers.Populate(cs)
	if cl != nil {
		this.Limits.Populate(*cl)
	}

	this.Policies.Populate(ca.Policies)

	if ca.Statistics != nil {
		this.Statistics.Populate(*ca.Statistics)
	}

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	self := fmt.Sprintf("/accounts/%s", ca.AccountID)
	this.Links.Self = lb.Link(self)
	this.Links.Transactions = lb.PagedLink(self, "transactions")
	this.Links.Operations = lb.PagedLink(self, "operations")
	this.Links.Payments = lb.PagedLink(self, "payments")

	return
}

func (a Account) PagingToken() string {
	return a.ID
}
