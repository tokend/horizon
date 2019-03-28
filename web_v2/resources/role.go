package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/regources/rgenerated"
)

// NewAccountRole - creates role core account Role
func NewAccountRole(role core2.AccountRole) rgenerated.AccountRole {
	return rgenerated.AccountRole{
		Key: NewAccountRoleKey(role.ID),
		Attributes: rgenerated.AccountRoleAttributes{
			Details: role.Details,
		},
		Relationships: rgenerated.AccountRoleRelationships{
			Rules: rgenerated.RelationCollection{},
		},
	}
}

//NewAccountRoleKey - returns new instance of key for account role
func NewAccountRoleKey(id uint64) rgenerated.Key {
	return rgenerated.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: rgenerated.ACCOUNT_ROLES,
	}
}

// NewSignerRole - maps signer role
func NewSignerRole(role core2.SignerRole) rgenerated.SignerRole {
	return rgenerated.SignerRole{
		Key: NewSignerRoleKey(role.ID),
		Attributes: rgenerated.SignerRoleAttributes{
			Details: role.Details,
		},
		Relationships: rgenerated.SignerRoleRelationships{
			Owner: NewAccountKey(role.OwnerID).AsRelation(),
			Rules: &rgenerated.RelationCollection{},
		},
	}
}

//NewSignerRoleKey - creates new key for signer role
func NewSignerRoleKey(id uint64) rgenerated.Key {
	return rgenerated.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: rgenerated.SIGNER_ROLES,
	}
}
