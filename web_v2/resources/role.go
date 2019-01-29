package resources

import "gitlab.com/tokend/regources/v2"

// NewRole - creates role from account address
func NewRole(accountAddress string) regources.Role {
	return regources.Role{
		Key: regources.Key{
			ID:   accountAddress,
			Type: regources.TypeRoles,
		},
		Attributes: regources.RoleAsstr{
			Details: map[string]interface{}{
				"name": "Name of the Mocked Role",
			},
		},
	}
}

// NewRoleKey - creates role key from account address
func NewRoleKey(accountAddress string) regources.Key {
	return regources.Key{
		ID:   accountAddress,
		Type: regources.TypeRoles,
	}
}
