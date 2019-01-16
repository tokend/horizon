package resources

import (
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewAccount - creates new account using core.Account. Returns nil if record is nil
func NewAccount(record *core.Account) *regources.Account {
	if record == nil {
		return nil
	}

	return &regources.Account{
		ID: record.Address,
	}
}

//NewAccountState - returns new instance of AccountState
func NewAccountState(record *core.Account) *regources.AccountState {
	return &regources.AccountState{
		ID: record.Address,
		BlockReasons: &regources.Mask{
			Mask:  record.BlockReasons,
			Flags: FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll),
		},
		IsBlocked: record.BlockReasons != 0,
	}
}
