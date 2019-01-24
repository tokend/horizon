package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageExternalSystemPool(id int64, details history2.ManageExternalSystemPoolDetails) *regources.ManageExternalSystemPool {
	result := &regources.ManageExternalSystemPool{
		Key: regources.NewKeyInt64(id, regources.TypeManageExternalSystemAccountIDPoolEntry),
	}

	switch details.Action {
	case xdr.ManageExternalSystemAccountIdPoolEntryActionCreate:
		result.Attributes.Create = new(regources.CreateExternalSystemPool)
		*result.Attributes.Create = regources.CreateExternalSystemPool(*details.Create)
	case xdr.ManageExternalSystemAccountIdPoolEntryActionRemove:
		result.Attributes.Remove = new(regources.RemoveExternalSystemPool)
		*result.Attributes.Remove = regources.RemoveExternalSystemPool(*details.Remove)
	default:
		panic(errors.From(errors.New("unexpected action for manage ex sys id pool"), logan.F{
			"action": details.Action,
		}))
	}

	return result
}
