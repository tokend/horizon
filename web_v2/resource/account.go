package resource

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/regources"
)

type Account struct {
	Data     AccountData   `json:"data"`
	Included []interface{} `json:"included,omitempty"`
}

func NewAccount(record *core.Account) *Account {
	return &Account{
		Data: AccountData{
			Id:   record.AccountID,
			Type: TypeAccounts,
			Attributes: &AccountAttributes{
				AccountType: AccountType{
					Type:  xdr.AccountType(record.AccountType).String(),
					TypeI: record.AccountType,
				},
				BlockReasons: AccountBlockReasons{
					Types: base.FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll),
					TypeI: record.BlockReasons,
				},
				IsBlocked: record.BlockReasons > 0,
				Policies: AccountPolicies{
					TypeI: record.Policies,
					Types: base.FlagFromXdrAccountPolicy(record.Policies, xdr.AccountPoliciesAll),
				},
				Thresholds: AccountThresholds{
					LowThreshold:  record.Thresholds[1],
					MedThreshold:  record.Thresholds[2],
					HighThreshold: record.Thresholds[3],
				},
			},
		},
	}
}

func (a *Account) IncludeBalances(balances []BalanceData) {
	for _, balance := range balances {
		a.Included = append(a.Included, balance)
	}
}

type AccountData struct {
	Id            string                `json:"id"`
	Type          string                `json:"type"`
	Attributes    *AccountAttributes    `json:"attributes,omitempty"`
	Relationships *AccountRelationships `json:"relationships,omitempty"`
}

type AccountRelationships struct {
	Referrer  *Account           `json:"referrer,omitempty"`
	Balances  *BalanceCollection `json:"balances,omitempty"`
	Referrals *AccountCollection `json:"referrals,omitempty"`
}

func (data *AccountData) RelateBalances(balances BalanceCollection) {
	if data.Relationships == nil {
		data.Relationships = &AccountRelationships{}
	}

	data.Relationships.Balances = balances.AsRelation()
}

type AccountAttributes struct {
	AccountType  AccountType         `json:"account_type"`
	BlockReasons AccountBlockReasons `json:"block_reasons"`
	IsBlocked    bool                `json:"is_blocked"`
	Policies     AccountPolicies     `json:"policies"`
	Thresholds   AccountThresholds   `json:"thresholds"`
}

type AccountBlockReasons struct {
	Types []regources.Flag `json:"types"`
	TypeI int32            `json:"type_i"`
}

type AccountType struct {
	Type  string `json:"type"`
	TypeI int32  `json:"type_i"`
}

type AccountThresholds struct {
	LowThreshold  byte `json:"low_threshold"`
	MedThreshold  byte `json:"med_threshold"`
	HighThreshold byte `json:"high_threshold"`
}

type AccountPolicies struct {
	TypeI int32            `json:"type_i"`
	Types []regources.Flag `json:"types"`
}
