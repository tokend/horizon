package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/regources/v2"
)

// NewAccountRole - creates role core account Role
func NewAccountRole(role core2.AccountRole) regources.AccountRole {
	return regources.AccountRole{
		Key: NewAccountRoleKey(role.ID),
		Attributes: regources.AccountRoleAttrs{
			Details: role.Details,
		},
		Relationships: regources.RoleRelation{
			Rules: &regources.RelationCollection{},
		},
	}
}

//NewAccountRoleKey - returns new instance of key for account role
func NewAccountRoleKey(id uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: regources.TypeAccountRoles,
	}
}

// NewSignerRole - maps signer role
func NewSignerRole(role core2.SignerRole) regources.SignerRole {
	return regources.SignerRole{
		Key: NewSignerRoleKey(role.ID),
		Attributes: regources.SignerRoleAttrs{
			Details: role.Details,
		},
		Relationships: regources.SignerRoleRelation{
			Owner: NewAccountKey(role.OwnerID).AsRelation(),
			Rules: &regources.RelationCollection{},
		},
	}
}

//NewSignerRoleKey - creates new key for signer role
func NewSignerRoleKey(id uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: regources.TypeSignerRoles,
	}
}
