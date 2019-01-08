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
	Id       string
	record   *core.Account
	response struct {
		Id         string             `json:"id"`
		Type       ResourceType       `json:"type"`
		Attributes attributes.Account `json:"attributes"`
	}
}

func NewAccount (id string) (*Account, error) {
	return &Account{
		Id: id,
	}, nil
}

func (a *Account) IsAllowed() (bool, error) {
	// TODO: can be optimized a bit
	return a.isSignedBy(a.Id) || a.isSignedByMaster(), nil
}

func (a *Account) Fetch() error {
	if a.record != nil {
		return nil
	}

	record, err := a.CoreQ().Accounts().ByAddress(a.Id)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch account")
	}

	a.record = record

	return nil
}

func (a *Account) Populate() error {
	a.response.Id = a.record.AccountID
	a.response.Type = TypeAccounts

	a.response.Attributes.AccountType.Type = xdr.AccountType(a.record.AccountType).String()
	a.response.Attributes.AccountType.TypeI = a.record.AccountType
	// TODO: move `FlagFromXdrBlockReasons` to regources?
	a.response.Attributes.BlockReasons.Types = base.FlagFromXdrBlockReasons(a.record.BlockReasons, xdr.BlockReasonsAll)
	a.response.Attributes.BlockReasons.TypeI = a.record.BlockReasons
	a.response.Attributes.IsBlocked = a.record.BlockReasons > 0
	a.response.Attributes.Policies.TypeI = a.record.Policies
	// TODO: move `FlagFromXdrAccountPolicy` to regources?
	a.response.Attributes.Policies.Types = base.FlagFromXdrAccountPolicy(a.record.Policies, xdr.AccountPoliciesAll)
	a.response.Attributes.Thresholds.HighThreshold = a.record.Thresholds[1]
	a.response.Attributes.Thresholds.HighThreshold = a.record.Thresholds[2]
	a.response.Attributes.Thresholds.HighThreshold = a.record.Thresholds[3]

	return nil
}

func (a *Account) Response() (interface{}, error) {
	return a.response, nil
}
