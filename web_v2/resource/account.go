package resource

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/horizon/web_v2/attributes"
)

type Account struct {
	Base `json:"-"`

	Id         string             `json:"id"`
	Type       ResourceType       `json:"type"`
	Attributes attributes.Account `json:"attributes"`

	record *core.Account `json:"-"`
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

func (a *Account) PopulateAttributes() error {
	record, err := a.CoreQ().Accounts().ByAddress(a.Id)

	if err != nil {
		return errors.New("Failed to get account by address")
	}

	a.Id = record.AccountID
	a.Type = TypeAccounts

	a.Attributes.AccountType.Type = xdr.AccountType(record.AccountType).String()
	a.Attributes.AccountType.TypeI = record.AccountType
	// TODO: move `FlagFromXdrBlockReasons` to regources?
	a.Attributes.BlockReasons.Types = base.FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll)
	a.Attributes.BlockReasons.TypeI = record.BlockReasons
	a.Attributes.IsBlocked = record.BlockReasons > 0
	a.Attributes.Policies.TypeI = record.Policies
	// TODO: move `FlagFromXdrAccountPolicy` to regources?
	a.Attributes.Policies.Types = base.FlagFromXdrAccountPolicy(record.Policies, xdr.AccountPoliciesAll)
	a.Attributes.Thresholds.HighThreshold = record.Thresholds[1]
	a.Attributes.Thresholds.HighThreshold = record.Thresholds[2]
	a.Attributes.Thresholds.HighThreshold = record.Thresholds[3]

	return nil
}

func (a *Account) Response() (interface{}, error) {
	return a, nil
}

type AccountCollection struct {
	Base `json:"-"`

	resources []Account
}

func (c *AccountCollection) FindOwner() error {
	return nil
}

func (c *AccountCollection) Fetch() error {
	return nil
}

func (c *AccountCollection) IsAllowed() (bool, error) {
	return c.isSignedByAdmin(), nil
}

func (c *AccountCollection) PopulateAttributes() error {
	for _, r := range c.resources {
		err := r.PopulateAttributes()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *AccountCollection) Response() (interface{}, error) {
	return c.resources, nil
}
