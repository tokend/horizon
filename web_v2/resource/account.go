package resource

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/horizon/web_v2/attributes"
)

type Account struct {
	Base

	Id         string
	Type       string
	Attributes attributes.Account
}

func (a *Account) FindOwner() error {
	record, err := a.CoreQ().Accounts().ByAddress(a.Id)
	if err != nil {
		return errors.Wrap(err, "Failed to get account by address")
	}
	a.Owner = record.AccountID
	return nil
}

func (a *Account) IsAllowed() (bool, error) {
	return a.isSignedByOwner() || a.isSignedByAdmin(), nil
}

func (a *Account) PopulateAttributes() error {
	record, err := a.CoreQ().Accounts().ByAddress(a.Id)

	if err != nil {
		return errors.New("Failed to get account by address")
	}

	a.Id = record.AccountID
	a.Type = TypeAccounts

	attrs := attributes.Account{}

	attrs.AccountType.Type = xdr.AccountType(record.AccountType).String()
	attrs.AccountType.TypeI = record.AccountType
	// TODO: move `FlagFromXdrBlockReasons` to regources?
	attrs.BlockReasons.Types = base.FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll)
	attrs.BlockReasons.TypeI = record.BlockReasons
	attrs.IsBlocked = record.BlockReasons > 0
	attrs.Policies.TypeI = record.Policies
	// TODO: move `FlagFromXdrAccountPolicy` to regources?
	attrs.Policies.Types = base.FlagFromXdrAccountPolicy(record.Policies, xdr.AccountPoliciesAll)
	attrs.Thresholds.HighThreshold = record.Thresholds[1]
	attrs.Thresholds.HighThreshold = record.Thresholds[2]
	attrs.Thresholds.HighThreshold = record.Thresholds[3]

	return nil
}
