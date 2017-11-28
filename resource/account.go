package resource

import (
	"fmt"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/httpx"
	"gitlab.com/tokend/horizon/render/hal"
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
	this.BlockReasons = ca.BlockReasons
	this.IsBlocked = ca.BlockReasons > 0
	this.AccountTypeI = ca.AccountType
	this.AccountType = xdr.AccountType(ca.AccountType).String()

	this.Referrer = ca.Referrer
	this.ShareForReferrer = amount.String(int64(ca.ShareForReferrer))
	this.CreatedAt = ca.CreatedAt
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
