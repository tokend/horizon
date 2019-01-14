package resource

import (
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/resource/base"
)

// AccountResponse - JSON:API response for Account resource
type AccountResponse struct {
	Data     Account            `json:"data"`
	Included includedCollection `json:"included,omitempty"`
}

func NewAccountResponse(record *core.Account) *AccountResponse {
	return &AccountResponse{
		Data: NewAccount(record),
	}
}

//IncludeBalances - includes balances into response
func (a *AccountResponse) IncludeBalances(balances []core.Balance) {
	a.Data.Relationships.Balances = &RelationshipCollection{}
	for i := range balances {
		balanceResource := NewBalance(&balances[i])
		a.Data.Relationships.Balances.Add(balanceResource.Key)
		a.Included.Add(&balanceResource)
		if balances[i].Asset != nil {
			assetResource := NewAsset(balances[i].Asset)
			a.Included.Add(&assetResource)
		}
	}
}

// Account - resource object representing AccountEntry
type Account struct {
	Key
	Attributes    *AccountAttributes    `json:"attributes,omitempty"`
	Relationships *AccountRelationships `json:"relationships,omitempty"`
}

//PopulateFromCore - populates account using core.Account
func NewAccount(record *core.Account) Account {
	return Account{
		Key: Key{
			ID:   record.Address,
			Type: typeAccounts,
		},
		Attributes: &AccountAttributes{
			Role: AccountRole{
				// TODO: must use account role
				ID:   int64(record.AccountType),
				Name: xdr.AccountType(record.AccountType).ShortString(),
			},
			BlockReasons: Mask{
				Mask:  record.BlockReasons,
				Flags: base.FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll),
			},
			IsBlocked: record.BlockReasons != 0,
		},
		Relationships: &AccountRelationships{},
	}
}

//AccountRelationships -represents reference from Account to other resource objects
type AccountRelationships struct {
	Referrer  *Key                    `json:"referrer,omitempty"`
	Balances  *RelationshipCollection `json:"balances,omitempty"`
	Referrals *RelationshipCollection `json:"referrals,omitempty"`
	Signers   *RelationshipCollection `json:"signers,omitempty"`
}

// AccountAttributes - represents information about Account
type AccountAttributes struct {
	Role         AccountRole `json:"role"`
	BlockReasons Mask        `json:"block_reasons"`
	IsBlocked    bool        `json:"is_blocked"`
}

// AccountRole - represents account role which defines actions allowed to be performed by this account
type AccountRole struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
