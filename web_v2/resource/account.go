package resource

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/horizon/web_v2/attributes"
)

type Account struct {
	Base
	id         string
	record     *core.Account
	attributes *attributes.Account
}

func NewAccount(id string) (*Account, error) {
	return &Account{id: id}, nil
}

func (a *Account) IsAllowed() (bool, error) {
	// TODO: can be optimized a bit
	return a.isSignedBy(a.id) || a.isSignedByMaster(), nil
}

func (a *Account) Fetch() error {
	if a.record != nil {
		return nil
	}

	record, err := a.CoreQ.Accounts().ByAddress(a.id)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch account")
	}

	a.record = record

	return nil
}

func (a *Account) Populate() error {
	a.attributes = &attributes.Account{}

	a.attributes.AccountType.Type = xdr.AccountType(a.record.AccountType).String()
	a.attributes.AccountType.TypeI = a.record.AccountType
	// TODO: move `FlagFromXdrBlockReasons` to regources?
	a.attributes.BlockReasons.Types = base.FlagFromXdrBlockReasons(a.record.BlockReasons, xdr.BlockReasonsAll)
	a.attributes.BlockReasons.TypeI = a.record.BlockReasons
	a.attributes.IsBlocked = a.record.BlockReasons > 0
	a.attributes.Policies.TypeI = a.record.Policies
	// TODO: move `FlagFromXdrAccountPolicy` to regources?
	a.attributes.Policies.Types = base.FlagFromXdrAccountPolicy(a.record.Policies, xdr.AccountPoliciesAll)
	a.attributes.Thresholds.HighThreshold = a.record.Thresholds[1]
	a.attributes.Thresholds.HighThreshold = a.record.Thresholds[2]
	a.attributes.Thresholds.HighThreshold = a.record.Thresholds[3]

	return nil
}

func (a *Account) Response() (Response, error) {
	response := Response{
		Id: a.id,
		Type: TypeAccounts,
		Attributes: a.attributes,
	}

	return response, nil
}
