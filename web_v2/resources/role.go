package resources

import (
	"encoding/json"
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/regources/v2"
)

// NewAccountRole - creates role from account address
func NewAccountRole(role core2.AccountRole) regources.AccountRole {
	var details regources.Details
	_ = json.Unmarshal([]byte(role.Details), &details)

	return regources.AccountRole{
		Key: regources.Key{
			ID:   strconv.FormatUint(role.ID, 10),
			Type: regources.TypeRoles,
		},
		Attributes: regources.AccountRoleAttrs{
			Details: details,
		},
	}
}

// NewAccountRoleKey - creates role key from account address
func NewAccountRoleKey(roleID uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(roleID, 10),
		Type: regources.TypeRoles,
	}
}

// NewAccountRole - creates role from account address
func NewSignerRole(role core2.SignerRole) regources.SignerRole {
	var details regources.Details
	_ = json.Unmarshal([]byte(role.Details), &details)

	return regources.SignerRole{
		Key: regources.Key{
			ID:   strconv.FormatUint(role.ID, 10),
			Type: regources.TypeRoles,
		},
		Attributes: regources.SignerRoleAttrs{
			Details: details,
			OwnerID: role.OwnerID,
		},
	}
}

// NewAccountRoleKey - creates role key from account address
func NewSignerRoleKey(roleID uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(roleID, 10),
		Type: regources.TypeSignerRoles,
	}
}
