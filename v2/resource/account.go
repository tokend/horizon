package resource

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/horizon/v2/model"
)

type Account struct {
	Base

	Id    string
	Model model.Account
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

func (a *Account) PopulateModel() error {
	record, err := a.CoreQ().Accounts().ByAddress(a.Id)

	if err != nil {
		return errors.New("Failed to get account by address")
	}

	accountModel := model.Account{}

	accountModel.Id = record.AccountID
	accountModel.ResourceType = TypeAccounts

	accountModel.Attributes = &model.AccountAttributes{}

	accountModel.Attributes.AccountType.Type = xdr.AccountType(record.AccountType).String()
	accountModel.Attributes.AccountType.TypeI = record.AccountType
	// TODO: move `FlagFromXdrBlockReasons` to regources?
	accountModel.Attributes.BlockReasons.Types = base.FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll)
	accountModel.Attributes.BlockReasons.TypeI = record.BlockReasons
	accountModel.Attributes.IsBlocked = record.BlockReasons > 0
	accountModel.Attributes.Policies.TypeI = record.Policies
	// TODO: move `FlagFromXdrAccountPolicy` to regources?
	accountModel.Attributes.Policies.Types = base.FlagFromXdrAccountPolicy(record.Policies, xdr.AccountPoliciesAll)
	accountModel.Attributes.Thresholds.HighThreshold = record.Thresholds[1]
	accountModel.Attributes.Thresholds.HighThreshold = record.Thresholds[2]
	accountModel.Attributes.Thresholds.HighThreshold = record.Thresholds[3]

	return nil
}

func (a *Account) MarshalModel() ([]byte, error) {
	return a.Model.MarshalSelf()
}
